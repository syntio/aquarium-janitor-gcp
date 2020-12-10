// Copyright Syntio d.o.o.
// All Rights Reserved
//
// Package rest contains the Schema registry REST Server configuration and start-up functions.
package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	service "github.com/syntio/schema-registry/business_logic"
	"github.com/syntio/schema-registry/model/dto"
)

//
// PostSchema is is a POST function that registers the received schema to the underlying database
// and returns it's given id and version.
//
// The expected input schema JSON should contain following fields:
// - Description   string
// - Type          int32
// - Specification string
// - Name          string
//
// Function writes back a JSON with fields:
// - Identification int64
// - Version        int32
// - Message        string
//
func PostSchema(w http.ResponseWriter, r *http.Request) {

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeInfoResponse(w, "Connection Error, could not read data", http.StatusServiceUnavailable)
		return
	}

	postRequest, err := PostRequestDeserializeJSON(requestBody)
	if err != nil {
		writeInfoResponse(w, "Bad request. Content-Type must be 'application/json'.", http.StatusBadRequest)
		return
	}

	response, err := service.CreateSchema(r.Context(), *postRequest)

	if err != nil {
		writeInfoResponse(w, "Server storage error! Schema was not registered.", http.StatusInternalServerError)
		return
	}

	writeValidResponse(w, response, http.StatusCreated)
}

// PostRequestDeserializeJSON deserializes the request body from a byte array to the
// resultung dto.SchemaDTO, an error is returned were the unmarshalling unsuccessful.
func PostRequestDeserializeJSON(requestBody []byte) (*dto.SchemaDTO, error) {
	postRequestInfo := dto.SchemaDTO{}
	err := json.Unmarshal(requestBody, &postRequestInfo)
	return &postRequestInfo, err
}
