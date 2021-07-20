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

// Package puller_cleaner_csv contains the main function.
// This component is suited for the use case in which there is an unexpected change in the messages sent.
// This repository implements the component as a GCP Cloud function.
package puller_cleaner_csv

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"

	"github.com/syntio/puller-cleaner-csv/cleaner"
	"github.com/syntio/puller-cleaner-csv/configuration"
	"github.com/syntio/puller-cleaner-csv/puller"
)

var projectID string

var Cfg configuration.Config

var subIdCSV string
var validTopic string
var invalidTopicJSON string
var deadLetterTopic string

var timeDurationSeconds time.Duration
var maxBatchSize int
var maxThroughput int

var contentType string

var csvValidatorURL string = os.Getenv("CSV_VALIDATOR_URL")
var schemaRegistryURL string = os.Getenv("SCHEMA_REGISTRY_URL")

const resourcePath string = "/schema/%s/evolution"

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

	subIdCSV = Cfg.Subscriptions.SubIdCSV
	validTopic = Cfg.Topics.ValidTopic
	invalidTopicJSON = Cfg.Topics.InvalidTopicJSON
	deadLetterTopic = Cfg.Topics.DeadLetterTopic

	timeDurationSeconds = Cfg.PullerCleanerCSV.TimeDurationSeconds
	maxBatchSize = Cfg.PullerCleanerCSV.MaxBatchSize
	maxThroughput = Cfg.PullerCleanerCSV.MaxThroughput

	contentType = Cfg.ContentType

}

// PullerCleaner is the main function.
// It pulls messages from the topic of CSV invalid messages
// and tries to clean them using the Schema Registry system.
//
// Input paramaters are http.ResponseWriter and http.Request
// because of the HTTP trigger on Cloud Function.
func PullerCleaner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Pull messages from PubSub subscription to a slice
	var msgs []pubsub.Message
	length, err := puller.Pull(&msgs, projectID, subIdCSV, timeDurationSeconds, maxBatchSize, maxThroughput)
	if length == 0 {
		if err != nil {
			log.Print(err)
		}

		log.Printf("Pulled 0 messages.")
		return
	}
	log.Printf("Pulled %d messages!", length)

	// Clean messages
	cleaner.Clean(ctx, msgs, projectID, validTopic, invalidTopicJSON, deadLetterTopic,
		schemaRegistryURL, schemaRegistryEvolutionURL, contentType, csvValidatorURL)

	fmt.Fprint(w, "Finished execution")
}
