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
//
// Package validator is used for message validation.
package validator

import (
	"reflect"

	"github.com/syntio/central-consumer/validator/impl"
)

// Init function is used to initialize the registry.
func init() {
	registerType(&impl.AvroValidator{})
	registerType(&impl.CsvValidator{})
	registerType(&impl.JsonValidator{})
	registerType(&impl.ProtobufValidator{})
	registerType(&impl.XmlValidator{})
}

var typeRegistry = make(map[string]reflect.Type)

// registerType function registers supported message types.
func registerType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	typeRegistry[t.PkgPath()+"."+t.Name()] = t
}

// Validator interface provides method Validate. Method is used for message validation.
type Validator interface {
	Validate(message, schema []byte) (bool, error)
}

// ValidatorFactory is used to create a specific validator for each message based on a message type.
func ValidatorFactory(validator string) Validator {
	return reflect.New(typeRegistry[validator]).Interface().(Validator)
}
