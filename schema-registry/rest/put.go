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

package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/syntio/schema-registry/business_logic"
	"github.com/syntio/schema-registry/model/dto"
)

// PutSchema registers a new schema version in the Schema Registry. The new version is connected to other schemas
// by ID from the request URL.
func PutSchema(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeInfoResponse(w, "Connection Error, could not read data", http.StatusServiceUnavailable)
		return
	}

	putRequestInfo, err := PutRequestDeserializeJSON(requestBody)
	if err != nil {
		writeInfoResponse(w, "Bad request. Content-Type must be 'application/json'.", http.StatusBadRequest)
		return
	}

	response, err := service.UpdateSchema(r.Context(), id, putRequestInfo, false)
	if err != nil {
		writeInfoResponse(w, "Could not update schema", http.StatusInternalServerError)
		return
	}
	writeValidResponse(w, response, http.StatusOK)
}

// PutRequestDeserializeJSON deserializes the request body from a byte array to the
// resultung dto.SpecificationDTO, an error is returned were the unmarshalling unsuccessful.
func PutRequestDeserializeJSON(requestBody []byte) (*dto.SpectificationDTO, error) {
	putRequestInfo := dto.SpectificationDTO{}
	err := json.Unmarshal(requestBody, &putRequestInfo)
	return &putRequestInfo, err
}
