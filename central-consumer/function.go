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

// Package central_consumer is a starting point of a cloud function. Messages that trigger central_consumer are
// transmitted to the supporting topic. Topic is determined based on a message content.
package central_consumer

import (
	"context"
	"encoding/base64"
	"log"
	"strings"

	"github.com/syntio/central-consumer/pubsub"
	"github.com/syntio/central-consumer/registry"
	"github.com/syntio/central-consumer/validator"
)

var b64coder *base64.Encoding
// Location of the specific validator implementations.
var baseLocation string

// Init function.
func init() {
	b64coder = base64.StdEncoding
	baseLocation = "github.com/syntio/central-consumer/validator/impl."
}

// CentralConsumerHandler represents a consumer handler function working with Schema Registry (Janitor).
//
// The function defines workflow steps for push-message processing. Workflow consists of checking the input message
// metadata validity, retrieving the required message schema from the Schema Registry by the ID and Version metadata
// and validating the input message with the retrieved schema.
//
// An error is returned if any errors occur during the function execution.
func CentralConsumerHandler(ctx context.Context, message pubsub.Message) error {
	valid, id, version, format := retrieveMetadata(message.Attributes)

	if !valid {
		log.Printf("Message metadata invalid, required fields are NOT set. ID = %v, Version = %v, Format = %v.\n",
			id, version, format)
		handleTransmission(message, pubsub.DeadLetterTopic)
		return nil
	}

	// Retrieve the message schema from the Schema Registry
	schemaInfo, found, err := registry.GetSchemaByIDAndVersion(id, version)

	if err != nil {
		log.Printf("ERROR: during schema registry request. %v.\n", err)
		handleTransmission(message, pubsub.DeadLetterTopic)
		return nil
	}

	// Validate and transmit the message depending on the retrieved schema
	if !found {
		_, invalidTopic, _, _ := chooseTopic(format)
		handleTransmission(message, invalidTopic)
	} else {
		transmitValidMessage(format, message, schemaInfo)

	}

	return nil
}

// retrieveMetadata returns if the metadata is valid and corresponding metadata values.
// Invalid metadata is logged with the corresponding description.
func retrieveMetadata(metadata map[string]string) (bool, string, string, string) {
	id := metadata["schemaId"]
	version := metadata["versionId"]
	format := metadata["format"]

	valid := len(id) != 0 && len(version) != 0 && len(format) != 0

	return valid, id, version, format
}

// handleValidationAndTransmission validates the input message with the retrieved message schema from the Schema Registry.
// Depending on the validation result, the message is transmitted to the validated, invalidated or error topic.
// Also, depending on the message format, the function receives as a parameter a validator for the specific message format.
func handleValidationAndTransmission(message pubsub.Message, schemaInfo *registry.Schema, validatedTopic, invalidatedTopic,
	errorTopic string, validator func([]byte, []byte) (bool, error)) {
	content := message.Data
	schema, err := b64coder.DecodeString(schemaInfo.SchemaDetails[0].Specification)

	if err != nil {
		log.Printf("ERROR: during base64 decoding. %v.\n", err)
		handleTransmission(message, errorTopic)
		return
	}

	valid, err := validator(content, schema)

	switch {
	case err != nil:
		log.Printf("ERROR: during message validation. %v\n", err)
		handleTransmission(message, errorTopic)
	case valid:
		handleTransmission(message, validatedTopic)
	default:
		handleTransmission(message, invalidatedTopic)
	}
}

// handleTransmission transmits the input message to a specific topic.
func handleTransmission(message pubsub.Message, topicName string) {
	if err := pubsub.Transmit(message, topicName); err != nil {
		log.Printf("ERROR: during message transmission to pubsub. %v.\n", err)
	}
}

// transmitValidMessage transmits a valid message to a specific topic.
func transmitValidMessage(format string, message pubsub.Message, schemaInfo *registry.Schema) {
	validator := validator.ValidatorFactory(baseLocation + strings.Title(format) + "Validator")
	validTopic, invalidTopic, errorTopic, err := chooseTopic(format)

	if err {
		handleTransmission(message, errorTopic)
	} else {
		handleValidationAndTransmission(message, schemaInfo, validTopic, invalidTopic, errorTopic, validator.Validate)
	}
 }

// chooseTopic returns corresponding formats based on a format of a message. For invalid input format, err is set to true.
func chooseTopic(format string) (validTopic, invalidTopic, errorTopic string, err bool) {
	validTopic = pubsub.ValidTopic
	invalidTopic = pubsub.DeadLetterTopic
	errorTopic = pubsub.DeadLetterTopic
	err = false

	if format == "json" {
		invalidTopic = pubsub.InvalidTopicJSON
		errorTopic = pubsub.InvalidTopicJSON
	} else if format == "csv" {
		invalidTopic = pubsub.InvalidTopicCSV
		errorTopic = pubsub.InvalidTopicCSV
	} else if !(format == "avro" || format == "protobuf" || format == "xml") {
		err = true
	}

	return
}
