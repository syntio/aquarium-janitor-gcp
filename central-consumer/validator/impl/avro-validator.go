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


import lib "github.com/hamba/avro"

// AvroValidator is a validator structure for Avro format.
type AvroValidator struct{}

// ValidateMessageWithAVROSchema validates an Avro message with a schema.
//
// Function returns the validation boolean result. An error is returned if any errors occur during the
// function execution.
func (av *AvroValidator) Validate(message, schema []byte) (bool, error) {
	var valid bool = false

	libSchema, err := lib.Parse(string(schema))
	if err != nil {
		return valid, err
	}

	var data interface{}
	if err = lib.Unmarshal(libSchema, message, &data); err != nil {
		return valid, err
	}

	return !valid, err
}
