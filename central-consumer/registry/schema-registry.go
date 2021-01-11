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

// Package registry is used for communication with the schema registry component.
package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/syntio/central-consumer/configuration"
)

var Cfg configuration.Config = configuration.RetrieveConfig()
var schemaRegistryURL = os.Getenv("SCHEMA_REGISTRY_URL")

// Schema represents a message schema from Schema Registry.
type Schema struct {
	Id            string           `json:"id,omitempty"`
	SchemaType    string           `json:"schema-type"`
	Autogenerated bool             `json:"autogenerated"`
	Description   string           `json:"description"`
	CreationDate  time.Time        `json:"creation-date"`
	Name          string           `json:"name"`
	SchemaDetails []*SchemaDetails `json:"schemas"`
}

// SchemaDetails represents details of a Schema.
type SchemaDetails struct {
	Version       int32  `json:"version"`
	Specification string `json:"specification"`
	SchemaHash    string `json:"schema-hash"`
}

// Report represents a info massage when the Schema Registry couldn't retrieve the required message schema.
type Report struct {
	Message string `json:"message"`
}

// GetSchemaByIDAndVersion retrieves a message schema from the Schema Registry using the schema ID and version.
// HTTP is used as the method of communication with the Schema Registry service.
//
// Function returns the SchemaInfo structure if the required schema is found (boolean indicator). An error is returned
// if any errors occur during the function execution.
func GetSchemaByIDAndVersion(id, version string) (*Schema, bool, error) {
	var schemaInfo *Schema
	var found bool = false
	var err error

	getURL := fmt.Sprintf("%s/schema/%s/version/%s", schemaRegistryURL, id, version)

	response, err := http.Get(getURL)
	if err != nil {
		return schemaInfo, found, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return schemaInfo, found, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		schemaInfo, err = JSONToSchema(responseBody)
		if err != nil {
			return schemaInfo, found, err
		}
		found = true
	} else {
		report, err := JSONToReport(responseBody)
		if err != nil {
			return schemaInfo, found, err
		}
		log.Printf("Schema couldn't be retrieved for ID = %s and Version = %s. Message: %s.\n", id, version, report.Message)
	}
	return schemaInfo, found, err
}

// JSONToSchema converts response body to a Schema structure and returns it. Error is set to nil if the conversion was successful.
func JSONToSchema(responseBody []byte) (*Schema, error) {
	schemaInfo := Schema{}
	if err := json.Unmarshal(responseBody, &schemaInfo); err != nil {
		return nil, err
	}
	return &schemaInfo, nil
}

// JSONToReport converts json message to a Report structure and returns it. Error is set to nil if the conversion was successful.
func JSONToReport(jsonMessage []byte) (*Report, error) {
	report := Report{}
	if err := json.Unmarshal(jsonMessage, &report); err != nil {
		return nil, err
	}
	return &report, nil
}
