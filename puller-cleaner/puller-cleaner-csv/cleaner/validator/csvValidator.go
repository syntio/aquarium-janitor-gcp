// Copyright 2020 Syntio Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


// Package validator handles validating messages with specified CSV schema.
package validator

import (
	"context"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"

	"github.com/syntio/puller-cleaner-csv/cleaner/filter"
	"github.com/syntio/puller-cleaner-csv/cleaner/sender"
)

type csvValidationRequest struct {
	Data   string `json:"data"`
	Schema string `json:"schema"`
}

type csvValidationResponse struct {
	Validation bool   `json:"validation"`
	Info       string `json:"info"`
}

// ValidateMessageWithSchemaCSV validates the given CSV message with the given CSV schema.
// It is done by using a separate cloud function for validation.
//
// Input parameters are message to be validated, schema to validate message with, cloud function 
// that does the validation, and a content type for communication with that cloud function.
//
// Output parameter is a bool that indicates whether message is successfully validated.
func ValidateMessageWithSchemaCSV(message, schema []byte, csvValidatorURL, contentType string) bool {
	// Create request
	req := csvValidationRequest{
		Data:   string(message),
		Schema: string(schema),
	}

	jsonReq, _ := json.Marshal(req)

	// Post request
	response, err := http.Post(csvValidatorURL, contentType, bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Println("ERROR: Can't connect to the CSV validator")
		return false
	}

	defer response.Body.Close()

	// Parse response
	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Save response to struct 'csvValidationResponse'
	responseStruct := &csvValidationResponse{}
	json.Unmarshal(responseByte, responseStruct)

	// If response is okay and contains information, return response 
	if response.StatusCode == http.StatusOK {
		return responseStruct.Validation
	} else {
		log.Printf("ERROR: Can't connect to the CSV validator. Response code: %d. Message: %s", 
			response.StatusCode, responseStruct.Info)
		return false
	}
}

// cleanMessage cleans the message that has been validated with one schema - it attaches schema's ID
// and version to a message, forwards message to valid topic, and deletes message from slice.
//
// Input parameters are context for communication with PubSub, slice of messages to remove message
// from, message to be cleaned, message's index in the slice, slice's current length, schema's ID
// and version that need to be attached to a message, and projectID and validTopic where message
// needs to be sent.
func cleanMessage(ctx context.Context, msgs *[]pubsub.Message, msg *pubsub.Message, i, length *int, 
	schemaIDstring, versionIDstring, projectID, validTopic string) {

	filter.SetAttributes(msg, schemaIDstring, versionIDstring)

	if forwarded, err := sender.ForwardMessage(ctx, projectID, validTopic, msg); !forwarded {
		log.Printf("Couldn't forward message to PubSub topic, but will delete it from slice: %v", err)
	}

	filter.RemoveFromSlice(msgs, i, length)
}

// CheckAgainstNewSchema checks all the messages in the slice if they follow the provided schema. If not, it just leaves the 
// message in the slice. If yes, it forwards message to a valid topic with newly specified schema and version, and it deletes 
// message from slice.
//
// Input parameters are context for communication with PubSub, messages to be checked, current length of slice of pulled messages, 
// schema specification that messages should be checked against, ID of message from which schema specification was inferred from, 
// schema's ID and version, project ID and topics to forward messages to corresponding places, cloud function that does the 
// validation, and a content type for communication with that cloud function.
func CheckAgainstNewSchema(ctx context.Context, msgs *[]pubsub.Message, length *int, schemaSpecification []byte, 
	firstMsgID, schemaIDstring, versionIDstring, projectID, validTopic, deadLetterTopic, csvValidatorURL, contentType string) {

	for i := 0; i < (*length); i++ {
		msg := (*msgs)[i]

		if isValid := ValidateMessageWithSchemaCSV(msg.Data, schemaSpecification, csvValidatorURL, contentType); !isValid {
			// Couldn't validate the very same message that Schema Registry inferred the schema from
			if msg.ID == firstMsgID {
				sender.ForwardAndDelete(ctx, projectID, deadLetterTopic, msg, msgs, length)
				break
			}
			continue
		}

		cleanMessage(ctx, msgs, &msg, &i, length, schemaIDstring, versionIDstring, projectID, validTopic)
	}
}
