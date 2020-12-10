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

// Package pubsub is used to transmit an input message to the provided topic. Message transmission is done using Transmit method.
package pubsub

import (
	lib "cloud.google.com/go/pubsub"
	"context"
	"os"

	"github.com/syntio/central-consumer/registry"
)

var projectID = os.Getenv("PROJECT_ID")
var ValidTopic = registry.Cfg.Topics.ValidTopic
var InvalidTopicJSON = registry.Cfg.Topics.InvalidTopicJSON
var InvalidTopicCSV = registry.Cfg.Topics.InvalidTopicCSV
var DeadLetterTopic = registry.Cfg.Topics.DeadLetterTopic

// Message represents a standard message which contains a payload and metadata. Structure is used for message
// re-transmission after its validation.
type Message struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

// Transmit transmits a message to the required topic.
//
// An error is returned if any errors occur during the function execution.
func Transmit(message Message, topicName string) error {
	ctx := context.Background()
	client, err := lib.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	topic := client.Topic(topicName)

	libMessage := lib.Message{
		Data:       message.Data,
		Attributes: message.Attributes,
	}

	result := topic.Publish(ctx, &libMessage)
	_, err = result.Get(ctx)
	if err != nil {
		return err
	}

	return nil
}