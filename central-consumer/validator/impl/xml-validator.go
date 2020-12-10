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

var xmlValidatorURL = registry.Cfg.Functions.XmlValidatorURL
var xmlRequestContentType = registry.Cfg.ContentType

// XmlValidator is a validator structure for xml format.
type XmlValidator struct {
	validatorURL string
}

// Validate validates a XML message with a schema. The validation is processed by a helper
// function over the network on xmlValidatorURL. Function uses HTTP to request the validation and retrieve its response.
//
// Function returns the validation boolean result. An error is returned if any errors occur during the
// function execution.
func (xv *XmlValidator) Validate(message, schema []byte) (bool, error) {
	if xv.validatorURL == "" {
		return ValidateMessageByHTTP(message, schema, xmlValidatorURL, xmlRequestContentType)
	} else {
		return ValidateMessageByHTTP(message, schema, xv.validatorURL, xmlRequestContentType)
	}
}
