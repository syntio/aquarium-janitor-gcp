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
	"github.com/syntio/central-consumer/registry"
)

var csvValidatorURL = registry.Cfg.Functions.CsvValidatorURL
var csvRequestContentType = registry.Cfg.ContentType

// CsvValidator is a validator structure for CSV format.
type CsvValidator struct {
	validatorURL string
}

// Validate validates a CSV message with a schema. The validation is processed by a helper
// function over the network on csvValidatorURL. Function uses HTTP to request the validation and retrieve its response.
//
// Function returns the validation boolean result. An error is returned if any errors occur during the
// function execution.
func (cv *CsvValidator) Validate(message, schema []byte) (bool, error) {
	if cv.validatorURL == "" {
		return ValidateMessageByHTTP(message, schema, csvValidatorURL, csvRequestContentType)
	} else {
		return ValidateMessageByHTTP(message, schema, cv.validatorURL, csvRequestContentType)
	}
}
