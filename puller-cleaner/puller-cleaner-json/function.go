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

// Package puller_cleaner_json contains the main function.
// This component is suited for the use case in which there is an unexpected change in the messages sent.
// This repository implements the component as a GCP Cloud function.
package puller_cleaner_json

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"

	"github.com/syntio/puller-cleaner-json/cleaner"
	"github.com/syntio/puller-cleaner-json/configuration"
	"github.com/syntio/puller-cleaner-json/puller"
)

var projectID string

var Cfg configuration.Config

var subIdJSON string
var validTopic string
var invalidTopicCSV string
var deadLetterTopic string

var timeDurationSeconds time.Duration
var maxBatchSize int
var maxThroughput int

var contentType string

var schemaRegistryURL string = os.Getenv("SCHEMA_REGISTRY_URL")
var resourcePath string = os.Getenv("EVOLUTION_PATH")
var schemaRegistryEvolutionURL string = schemaRegistryURL + resourcePath

func init() {
	projectID = os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Printf("ERROR: Couldn't read PROJECT_ID environment variable.")
		return
	}

	Cfg, err := configuration.RetrieveConfig()
	if err != nil {
		log.Printf("ERROR: Couldn't read configuration parameters from a config file: %v", err)
		return
	}

	subIdJSON = Cfg.Subscriptions.SubIdJSON
	validTopic = Cfg.Topics.ValidTopic
	invalidTopicCSV = Cfg.Topics.InvalidTopicCSV
	deadLetterTopic = Cfg.Topics.DeadLetterTopic

	timeDurationSeconds = Cfg.PullerCleanerJSON.TimeDurationSeconds
	maxBatchSize = Cfg.PullerCleanerJSON.MaxBatchSize
	maxThroughput = Cfg.PullerCleanerJSON.MaxThroughput

	contentType = Cfg.ContentType
}

// PullerCleaner is the main function.
// It pulls messages from the topic of JSON invalid messages
// and tries to clean them using the Schema Registry system.
//
// Input paramaters are http.ResponseWriter and http.Request
// because of the HTTP trigger on Cloud Function.
func PullerCleaner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Pull messages from PubSub subscription to a slice
	var msgs []pubsub.Message
	length, err := puller.Pull(&msgs, projectID, subIdJSON, timeDurationSeconds, maxBatchSize, maxThroughput)
	if length == 0 {
		if err != nil {
			log.Print(err)
		}

		log.Printf("Pulled 0 messages.")
		return
	}
	log.Printf("Pulled %d messages!", length)

	// Clean messages
	cleaner.Clean(ctx, msgs, projectID, validTopic, invalidTopicCSV, deadLetterTopic,
		schemaRegistryURL, schemaRegistryEvolutionURL, contentType)

	fmt.Fprint(w, "Finished execution")
}
