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


// Package validator handles validating messages with specified JSON schema.
package validator

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/xeipuuv/gojsonschema"

	"github.com/syntio/puller-cleaner-json/cleaner/filter"
	"github.com/syntio/puller-cleaner-json/cleaner/sender"
)

// ValidateMessageWithSchemaJSON validates the given JSON message with the 
// given JSON schema.
//
// Input parameters are the message to be validated and the schema by which 
// the message is to be validated.
//
// Output parameters are bool which indicates whether message can be validated 
// with a given schema or not, and a possible error occurred while validation.
func validateMessageWithSchemaJSON(message, schema []byte) (bool, error) {
	documentLoader := gojsonschema.NewBytesLoader(message)
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	
	messageFollowsSchema, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, err
	}

	return messageFollowsSchema.Valid(), nil
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

// CheckAgainstNewSchema checks all the messages in the slice if they follow the provided schema. If not, it 
// just leaves the message in the slice. If yes, it forwards message to a valid topic with newly specified 
// schema and version, and it deletes message from slice.
//
// Input parameters are context for communication with PubSub, messages to be checked, current length of slice 
// of pulled messages, schema specification that messages should be checked against, ID of message from which 
// schema specification was inferred from, schema's ID and version, project ID and valid topic to forward those 
// messages that follow the schema.
func CheckAgainstNewSchema(ctx context.Context, msgs *[]pubsub.Message, length *int, schemaSpecification []byte, 
	firstMsgID, schemaIDstring, versionIDstring, projectID, validTopic, deadLetterTopic string) {

	for i := 0; i < (*length); i++ {
		msg := (*msgs)[i]

		if isValid, _ := validateMessageWithSchemaJSON(msg.Data, schemaSpecification); !isValid {
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
