/// Copyright 2020 Syntio Inc.

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

import lib "github.com/xeipuuv/gojsonschema"

// JsonValidator is a validator structure for JSON format.
type JsonValidator struct{}

// Validate validates a JSON message with a schema.
//
// Function returns the validation boolean result. An error is returned if any errors occur during the
// function execution.
func (jv *JsonValidator) Validate(message, schema []byte) (bool, error) {
	var valid bool = false

	documentLoader := lib.NewBytesLoader(message)
	schemaLoader := lib.NewBytesLoader(schema)

	result, err := lib.Validate(schemaLoader, documentLoader)
	if err != nil {
		return valid, err
	}

	valid = result.Valid()
	return valid, err
}
