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


// Package cleaner implements the cleaner part of puller/cleaner component. It cleans the messages by communicating with Schema Registry.
package cleaner

import (
	"context"
	"encoding/base64"
	"log"

	"cloud.google.com/go/pubsub"

	"github.com/syntio/puller-cleaner-csv/cleaner/filter"
	"github.com/syntio/puller-cleaner-csv/cleaner/registry"
	"github.com/syntio/puller-cleaner-csv/cleaner/sender"	
	"github.com/syntio/puller-cleaner-csv/cleaner/validator"
)

// b64coder is used to encode/decode schemas.
var b64coder *base64.Encoding

// init is executed before anything else in the file. 
// It initiates the base64 encoder.
func init() {
	b64coder = base64.StdEncoding
}

// schemaBase64Decode decodes given schema into a byte array.
//
// Input parameter is a schema as a string.
//
// Output parameters are a byte representation of a schema, a bool 
// which indicates whether schema is successfully decoded, and a
// possible error occurred while decoding.
func schemaBase64Decode(schemaBase64 string) ([]byte, bool, error) {
	decodedBytes, err := b64coder.DecodeString(schemaBase64)
	if err != nil {
		return nil, false, err
	}
	
	return decodedBytes, true, nil
}

// Clean cleans the messages (infers schema from them, or sends them to dead letter topic).
// It goes through messages and then tries to: 
//		1. infer schema from message, 
//		2. retrieve schema from SR, 
//		3. retrieve schema specification, 
//		4. validate the rest of the messages with the retrieved schema
//
// Input parameters are context for communication with PubSub, array of messages that need to be cleaned, 
// project's and topics' names where resolved messages need to be sent to, URL of Schema Registry for 
// communication with Schema Registry, URL of Schema Registry's Evolution component for communication with 
// Schema Registry about schema evolution, and a content type for communication with Schema Registry's REST server.
func Clean(ctx context.Context, msgs []pubsub.Message, projectID, validTopic, invalidTopicJSON, deadLetterTopic, 
	schemaRegistryURL, schemaRegistryEvolutionURL, contentType, csvValidatorURL string) {
	
	// First remove all messages with invalid format (faulty metadata or non-corresponding formats)
	filter.RemoveInvalidFormats(ctx, &msgs, projectID, invalidTopicJSON, deadLetterTopic)
	
	length := len(msgs)

	// 1. Go through all the messages and try to infer the schema
	for length != 0 {
		msgFirst := msgs[0]

		// If schema is not inferred
		if schemaInfo, isCleaned, err := registry.InferSchema(msgFirst, schemaRegistryEvolutionURL, contentType); !isCleaned {
			log.Printf("Couldn't infer the schema from a message, it is a dead letter: %v", err)
			sender.ForwardAndDelete(ctx, projectID, deadLetterTopic, msgFirst, &msgs, &length)

		// If schema is inferred
		} else {
			schemaIDstring, versionIDstring := registry.ExtractIDAndVersion(schemaInfo)

			// 2. Retrieve the schema
			schema, found, err := registry.GetSchema(schemaIDstring, versionIDstring, schemaRegistryURL)
			if !found {
				log.Printf("Couldn't retrieve the schema, it is a dead letter: %v", err)
				sender.ForwardAndDelete(ctx, projectID, deadLetterTopic, msgFirst, &msgs, &length)

				continue
			}

			// 3. Retrieve the schema specification (to validate other messages with it)
			schemaSpecBytes, success, err := schemaBase64Decode(schema.SchemaDetails[0].Specification)
			if !success {
				log.Printf("Couldn't parse the schema using base64 decoding, it is a dead letter: %v", err)
				sender.ForwardAndDelete(ctx, projectID, deadLetterTopic, msgFirst, &msgs, &length)

				continue
			}

			// 4. Let's start cleaning the rest of the messages with the retrieved schema specification!
			validator.CheckAgainstNewSchema(ctx, &msgs, &length, schemaSpecBytes, msgFirst.ID, schemaIDstring, 
				versionIDstring, projectID, validTopic, deadLetterTopic, csvValidatorURL, contentType)
			log.Printf("Just cleaned input down to %d messages", length)
		}
	}
}