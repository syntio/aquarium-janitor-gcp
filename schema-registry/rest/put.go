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
