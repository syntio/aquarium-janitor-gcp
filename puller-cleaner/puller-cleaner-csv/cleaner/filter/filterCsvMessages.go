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


// Package filter handles removing messages with non-CSV format (and/or invalid metadata). It also implements several helper function.
package filter

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"

	"github.com/syntio/puller-cleaner-csv/cleaner/sender"
)

// GetAttributes extracts messages's attributes schemaID and format.
//
// Input parameter is a message that attributes need to be extracted from.
//
// Output parameters are schemaID and format, and the indicators of whether 
// the attributes are well extracted.
func GetAttributes(msg pubsub.Message) (string, bool, string, bool) {
	att := msg.Attributes

	schemaID, foundSchemaID := att["schemaId"]
	format, foundFormat := att["format"]

	return schemaID, foundSchemaID, format, foundFormat
}

// SetAttributes sets schema ID and version ID attributes in a message.
//
// Input parameters are a PubSub message which attributes need to be set, and
// schema ID and version ID to set.
func SetAttributes(msg *pubsub.Message, schemaIDstring, versionIDstring string) {
	(*msg).Attributes["schemaId"] = schemaIDstring
	(*msg).Attributes["versionId"] = versionIDstring
}

// isDeadCSVLetter checks whether the message is an obvious CSV 
// dead letter. It is, if it does not contain the metadata "schemaId" 
// and "format", or if the format is different of the "csv" value.
//
// Input parameter is a message that needs to be checked if it is 
// a dead letter or a valid CSV message.
//
// Output paramaters are indicator of whether a message is dead letter, 
// indicator of whether schema ID is found, and messages's format.
func isDeadCSVLetter(msg pubsub.Message) (bool, bool, string) {
	_, foundSchemaID, format, foundFormat := GetAttributes(msg)

	if !foundSchemaID || !foundFormat || format != "csv" {
		return true, foundSchemaID, format
	}

	return false, foundSchemaID, format
}

// RemoveFromSlice is a helper function for removing message 
// (specified with its index) from a slice of messages.
//
// Input parameters are slice of messages containing the message, 
// message's index in the slice, and the current length of the slice.
func RemoveFromSlice(msgs *[]pubsub.Message, i, n *int) {
	(*msgs) = append((*msgs)[:(*i)], (*msgs)[(*i)+1:]...)

	(*i)-- // Since we just deleted a[i], we must redo that index
	(*n)--
}

// RemoveInvalidFormats removes from slice all the messages that are not of the wanted 
// format to their respective topic and forwards them to the corresponding topic.
//
// Input parameters are context for communication with PubSub, slice of messages from 
// which the message should be removed, and project's and topics' IDs where removed 
// messages should be forwarded to.
func RemoveInvalidFormats(ctx context.Context, msgs *[]pubsub.Message, projectID, 
	invalidTopicJSON, deadLetterTopic string) {

	n := len(*msgs)

	for i := 0; i < n; i++ {
		msg := (*msgs)[i]

		isDeadCSVmessage, isFoundSchemaID, format := isDeadCSVLetter(msg)
		if !isDeadCSVmessage {
			continue
		}

		var topic string

		if format == "json" && isFoundSchemaID {
			topic = invalidTopicJSON
		} else {
			topic = deadLetterTopic
		}

		if forwarded, err := sender.ForwardMessage(ctx, projectID, topic, &msg); !forwarded {
			log.Printf("Couldn't forward message to PubSub topic, but will delete it from slice: %v", err)
		}

		RemoveFromSlice(msgs, &i, &n)
	}
}
