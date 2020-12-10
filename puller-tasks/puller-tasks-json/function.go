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
//
// Package puller_tasks_json contains the main function.
// It is invoked as a Cloud Function in regular time invervals using the Cloud Scheduler service
// It creates Cloud Tasks that invoke the Puller & Cleaner Cloud Function, which does the schema recovery and evolution
package puller_tasks_json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/syntio/puller-tasks-json/configuration"
	"github.com/syntio/puller-tasks-json/task"
)

var projectID string
var region string

var Cfg configuration.Config

var pullerCleanerURL string
var pullerTaskQueue string

type Request struct {
	PullerNumber int `json:"puller-number"`
}

// Load the neccessary parameters from the configuration file and the environment
func init() {
	projectID = os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Printf("ERROR: Couldn't read PROJECT_ID environment variable.")
		return
	}

	region = os.Getenv("REGION")
	if region == "" {
		log.Printf("ERROR: Couldn't read REGION environment variable.")
		return
	}

	Cfg, err := configuration.RetrieveConfig()
	if err != nil {
		log.Printf("ERROR: Couldn't read configuration parameters from a config file: %v", err)
		return
	}

	pullerCleanerURL = Cfg.Functions.PullerCleanerJsonURL
	pullerTaskQueue = Cfg.PullerTaskQueue
}

// The main function that is invoked by the cloud scheduler to create Cloud Tasks
// Input paramaters are http.ResponseWriter and http.Request
// because of the HTTP trigger on Cloud Function.
func HTTPPullerTasks(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var request Request
	if err = json.Unmarshal(requestBody, &request); err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < request.PullerNumber; i++ {
		if err = task.CreateHTTPTargetTask(projectID, region, pullerTaskQueue, pullerCleanerURL); err != nil {
			log.Println(err)
			return
		}
		log.Println("Puller function number: ", i+1)
	}
}
