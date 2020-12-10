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


// Package puller implements the puller part of puller/cleaner component. It pulls the messages from PubSub topic.
package puller

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// receiveMessages receives messages from PubSub topic via provided PubSub subscription.
// It uses asynchronous function *Subscription.Receive, but ensures synchronisation by using mutex.
// Pulling is limited by maximum number of pulled messages, bytes (throughput) pulled, and seconds 
// of pulling.
//
// Input parameters are context for connection with PubSub, slice where messages will be saved to, 
// PubSub subscription which pulls messages, and limits for receiving messages.
//
// Output is an error occured while receiving messages.
func receiveMessages(cctx context.Context, msgs *[]pubsub.Message, sub *pubsub.Subscription, 
	timeDurationSeconds time.Duration, maxBatchSize, maxThroughput int) error {

	// Define in context that you will receive messages for 'timeDurationSeconds' seconds
	cctx, ccancel := context.WithTimeout(cctx, timeDurationSeconds*time.Second)

	// Create a channel to handle messages to as they come in
	var mu sync.Mutex
	current_bytes_taken := 0
	
	// Receive messages until the passed in context is done
	err := sub.Receive(cctx, func(cctx context.Context, msg *pubsub.Message) {
		*msgs = append(*msgs, *msg)
		msg.Ack()
		mu.Lock()
		defer mu.Unlock()

		current_bytes_taken += len(msg.Data)

		if len(*msgs) >= maxBatchSize || current_bytes_taken >= maxThroughput {
			ccancel()
		}
	})

	return err
}

// Pull pulls messages from PubSub topic by using specified PubSub subscription and GCP project.
// Pulling is limited by specified maximum number of messages, bytes (throughput), and seconds of pulling.
//
// Input parameters are empty array of messages to fill in, project's and topic's subscription's IDs, 
// duration of pulling, maximum number of messages to pull, and a maximum throughput.
//
// Output is the length of the array of pulled messages, and the error occured while receiving messages.
func Pull(msgs *[]pubsub.Message, projectID, subIdCSV string, timeDurationSeconds time.Duration, 
	maxBatchSize, maxThroughput int) (int, error) {
	
	cctx := context.Background()

	client, err := pubsub.NewClient(cctx, projectID)
	if err != nil {
		return 0, fmt.Errorf("ERROR: Couldn't make connection to pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subIdCSV)

	err = receiveMessages(cctx, msgs, sub, timeDurationSeconds, maxBatchSize, maxThroughput)
	if err != nil && status.Code(err) != codes.Canceled {
		return 0, fmt.Errorf("ERROR: Couldn't receive messages from PubSub topic: %v", err)
	}

	return len(*msgs), nil
}
