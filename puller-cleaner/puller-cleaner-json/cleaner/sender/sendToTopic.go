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


// Package sender handles sending messages to a specified PubSub topic.
package sender

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

// ForwardMessage forwards the given message to the specified topic in the specified project.
//
// Input parameters are context for the connection with PubSub, project's and topic's IDs where message needs 
// to be sent, and a message.
//
// Output parameters are a bool which indicates whether or not the message was successfully forwarded to the
// topic, and a possible error occurred while connection and sending to PubSub.
func ForwardMessage(ctx context.Context, projectID, topicName string, message *pubsub.Message) (bool, error) {
	clientPubSub, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("ERROR: Couldn't establish connection to PubSub to forward the message.")
	}
	defer clientPubSub.Close()

	topic := clientPubSub.Topic(topicName)
	result := topic.Publish(ctx, message)
	
	_, err = result.Get(ctx)
	if err != nil {
		return false, fmt.Errorf("ERROR: Couldn't publish the message.")
	}

	return true, nil
}

// ForwardAndDelete forwards message to the specified topic and deletes it from the slice msgs.
//
// Input parameters are context for the connection with PubSub, project's and topic's IDs where 
// message needs to be sent, a message that needs to be sent, slice of messages that the message 
// will be deleted from, and slice's current length.
func ForwardAndDelete(ctx context.Context, projectID, deadLetterTopic string, 
	message pubsub.Message, msgs *[]pubsub.Message, length *int) {

	if forwarded, err := ForwardMessage(ctx, projectID, deadLetterTopic, &message); !forwarded {
		log.Printf("Couldn't forward message to PubSub topic, but will delete it from slice: %v", err)
	}

	deleteFirstMessage(msgs, length) // we know that message is the first element in msgs
}

// deleteFirstMessage is a helper function that removes first 
// message from a slice of messages.
//
// Input parameters are messages to remove message from, and
// its current length.
func deleteFirstMessage(msgs *[]pubsub.Message, length *int) {
	// Erase element
	(*msgs)[0] = (*msgs)[(*length) - 1] // Copy last element to index 0.
	(*msgs) = (*msgs)[:(*length)-1]     // Truncate slice.

	// Decrease number of messages in the slice
	(*length)--
}