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
	"log"
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/syntio/schema-registry/business_logic"
	"github.com/syntio/schema-registry/model/dto"
)

// EvolutionSchema registers a new schema from the input message. Evolved schema is connected with other schemas
// by ID from the request URL. Evolved schema has an new, incremented version.
func EvolutionSchema(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeInfoResponse(w, "Connection Error, could not read data", http.StatusServiceUnavailable)
		return
	}

	evolutionRequest, err := EvolutionRequestDeserializeJSON(requestBody)
	if err != nil {
		writeInfoResponse(w, "Bad request. Content-Type must be 'application/json'.", http.StatusBadRequest)
		return
	}

	generatedSchema, isGenerated, err := service.Evolve(*evolutionRequest)

	if err != nil {
		writeInfoResponse(w, "Schema dynamic generation error!", http.StatusInternalServerError)
		return
	}

	if !isGenerated {
		writeInfoResponse(w, "Schema couldn't be generated, dead-letter message", http.StatusOK)
	}

	response, err := service.UpdateSchema(r.Context(), id,
		&dto.SpectificationDTO{Specification: string(generatedSchema)}, true)

	log.Println("Sucessfully updated")
	if err != nil {
		writeInfoResponse(w, "Could not update schema", http.StatusInternalServerError)
		return
	}
	writeValidResponse(w, response, http.StatusOK)

}

// EvolutionRequestDeserializeJSON deserializes the request body from a byte array to the
// resultung dto.EvolutionDTO, an error is returned were the unmarshalling unsuccesful.
func EvolutionRequestDeserializeJSON(requestBody []byte) (*dto.EvolutionDTO, error) {
	evolutionRequestInfo := dto.EvolutionDTO{}
	err := json.Unmarshal(requestBody, &evolutionRequestInfo)
	return &evolutionRequestInfo, err
}
