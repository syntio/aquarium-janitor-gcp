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


// Package registry handles sending requests and receiving responses from Schema Registry's REST server regarding Schema Evolution.
// This file implements handler for Schema Evolution.
package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"cloud.google.com/go/pubsub"

	"github.com/syntio/puller-cleaner-json/cleaner/filter"
)

type postResponse struct {
	Identification string `json:"identification"`
	Version        int32  `json:"version"`
	Message        string `json:"message"`
}

// InferSchema is the Schema Evolution handler function.
// It sends the message and its format to the Schema Registry. Schema Registry's REST server returns whether 
// the schema is inferred from the message, and if so, the ID and version of derived schema.
//
// Input parameters are a message that a schema needs to be inferred from, an URL of Schema Registry's Evolution 
// component for communication with Schema Registry about evolution, and a content type for communication with 
// Schema Registry's REST server.
//
// Output parameters are a struct that contains details of the inferred schema, a bool which indicates whether 
// schema is successfully inferred or not, and a possible error occurred while communication with Schema Registry.
func InferSchema(msg pubsub.Message, schemaRegistryEvolutionURL, contentType string) (*postResponse, bool, error) {
	schemaIDstring, _, format, _ := filter.GetAttributes(msg)

	// Create request to send to Schema Registry
	sendMessage := map[string]string{"Data": string(msg.Data), "format": format}
	jsonRequest, err := json.Marshal(&sendMessage)
	if err != nil {
		return nil, false, fmt.Errorf("ERROR: Couldn't marshal request to send to Schema Registry to infer schema: %v", err)
	}

	// Send request to Schema Registry
	surl := fmt.Sprintf(schemaRegistryEvolutionURL, schemaIDstring)
	response, err := http.Post(surl, contentType, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, false, fmt.Errorf("ERROR: While sending schema registration request: %v", err)
	}

	defer response.Body.Close()

	// Read response (response is incoming byte stream)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, false, fmt.Errorf("ERROR: While reading schema registration response: %v", err)
	}

	// Save Schema Registry's response (info about schema) to struct 'postResponseInfo' (it contains id, version, message)
	postResponseInfo := &postResponse{}
	json.Unmarshal(responseBody, postResponseInfo)

	// If response from Schema Registry is okay and contains information, return response 
	if response.StatusCode == http.StatusOK {
		if (&postResponseInfo.Identification != nil && &postResponseInfo.Version != nil) &&
			(postResponseInfo.Identification != "" && postResponseInfo.Version != 0) {
			return postResponseInfo, true, nil
		} else {
			return nil, false, fmt.Errorf("ERROR: Schema's ID or version received from Schema Registry is not ok: ID = %s, version = %d", 
				postResponseInfo.Identification, postResponseInfo.Version)
		}

	} else {
		return nil, false, fmt.Errorf("ERROR: Unhandled error while inferring schema, status code: %d", response.StatusCode)
	}
}

// convertToString is a helper function.
// It converts int32 variable into a string.
//
// Input parameter is int32 variable.
//
// Output is converted string.
func convertToString(number int32) string {
	base := 10
	
	return strconv.FormatInt(int64(number), base)
}

// ExtractIDAndVersion extracts ID and version from schema, as strings.
//
// Input parameter is a schema.
//
// Output parameters are ID and version of input schema, as strings.
func ExtractIDAndVersion(schemaInfo *postResponse) (string, string){
	schemaIDstring := schemaInfo.Identification
	versionIDstring := convertToString(schemaInfo.Version)

	return schemaIDstring, versionIDstring
}