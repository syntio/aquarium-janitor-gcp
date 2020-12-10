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

// Package impl represents implementation of message validation process.
package impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ValidationRequest contains the message and schema which are used by the validator. The structure represents a HTTP
// request body.
type ValidationRequest struct {
	Data   string `json:"data"`
	Schema string `json:"schema"`
}

// ValidationResponse contains the validation result and an info message. The structure represents a HTTP response body.
type ValidationResponse struct {
	Validation bool   `json:"validation"`
	Info       string `json:"info"`
}

// Function returns HTTP request body as a byte array and error indicator. Error is set if the body can't be extracted.
func ValidationRequestSerializeJSON(requestBody *ValidationRequest) ([]byte, error) {
	jsonRequest, err := json.Marshal(requestBody)
	return jsonRequest, err
}

// Function returns ValidationResponse and an error ih the transfer wasn't successful.
func ValidationResponseDeserializeJSON(responseBody []byte) (*ValidationResponse, error) {
	validationResponse := ValidationResponse{}
	err := json.Unmarshal(responseBody, &validationResponse)
	return &validationResponse, err
}

// ValidateMessageByHTTP requests a message validation by HTTP. The function represents a generic HTTP validator request.
//
// Function returns the validation boolean result. An error is returned if any errors occur during the
// function execution.
func ValidateMessageByHTTP(message, schema []byte, validatorURL, requestContentType string) (bool, error) {
	var valid bool = false

	requestStructure := ValidationRequest{
		Data:   string(message),
		Schema: string(schema),
	}
	requestJSON, err := ValidationRequestSerializeJSON(&requestStructure)
	if err != nil {
		return valid, err
	}

	validatorResponse, err := http.Post(validatorURL, requestContentType, bytes.NewBuffer(requestJSON))
	if err != nil {
		return valid, err
	}
	defer validatorResponse.Body.Close()

	validatorBody, err := ioutil.ReadAll(validatorResponse.Body)
	if err != nil {
		return valid, err
	}

	validatorBodyStruct, err := ValidationResponseDeserializeJSON(validatorBody)
	if err != nil {
		return valid, err
	}

	switch validatorResponse.StatusCode {
	case http.StatusOK:
		valid = validatorBodyStruct.Validation
	case http.StatusBadRequest:
		err = fmt.Errorf("bad request: %s", validatorBodyStruct.Info)
	default:
		err = fmt.Errorf("error: status code [%v]", validatorResponse.StatusCode)
	}

	return valid, err
}
