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

// Package rest contains the Schema registry REST Server configuration and start-up functions.
package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syntio/schema-registry/model/dto"
)

//
// SetupAndStartServer starts the REST server and configures it on the port localhost:8080.
//
// Configuration includes listening on handles:
//  - "/schema/{id}/version/{version}" for schema retrieval
//  - "/schema" for schema registration
//  - "/schema/{id} for schema versioning
//	- "/schema/{id}/evolution for schema evolution
//	- "/schema/resolver/backward-transite/{id} for schema list retrieval
//
func SetupAndStartServer() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/schema/{id}/version/{version}", GetSchemaByIdAndVersion).Methods("GET")
	router.HandleFunc("/schema/", PostSchema).Methods("POST")
	router.HandleFunc("/schema/{id}", PutSchema).Methods("PUT")
	router.HandleFunc("/schema/{id}/evolution", EvolutionSchema).Methods("POST")
	router.HandleFunc("/schema/resolver/{id}", BackwardResolver).Methods("GET")

	fmt.Println("Schema register REST server ready on port :8080")
	fmt.Println(http.ListenAndServe(":8080", router))
}

// wrieteValid response writes any informational or error response into the designated writer.
// Input arguments are a http writer, a serialized response and a HTTP status code.
func writeInfoResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResponse, err := InfoResponseSerializeJSON(message)
	if err != nil {
		log.Printf("Info response could't be serialized properly.\nError: %s", err)
		return
	}
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("Get HTTP response couldn't be send properly.\nError: %s", err)
		return
	}
}

// wrieteValid response writes any response containing a body into the designated writer.
// Input arguments are a http writer, a serialized response and a HTTP status code.
func writeValidResponse(w http.ResponseWriter, response []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err := w.Write(response)
	if err != nil {
		log.Printf("Put HTTP response couldn't be send properly.\nError: %s", err)
		return
	}
}

//
// InfoResponseSerializeJSON serializes the input message to a simple JSON object containing only one field  - message.
//
func InfoResponseSerializeJSON(message string) ([]byte, error) {
	jsonResponse, err := json.Marshal(dto.ReportDTO{Message: message})
	return jsonResponse, err
}
