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

package dto

// EvolutionDTO represents a schema evolution request. Data defines the input message from which a new schema is
// registered (evolved).
type EvolutionDTO struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

// Struct containing information needed to register a schema.
// It is used as INPUT in the PostSchema function.
type SchemaDTO struct {
	Description   string `json:"description"`
	Specification string `json:"specification"`
	Name          string `json:"name"`
	SchemaType    string `json:"schema-type"`
}

// Structure ReportDTO is a simple wrapper of the system's message for the user.
// It enables easier parsing to a JSON object.
type ReportDTO struct {
	Message string `json:"message"`
}

// SpectificationDTO represents the schema versioning request. Specification defines the new version of the schema.
type SpectificationDTO struct {
	Specification string `json:"specification"`
}

// InsertInfoDTO represents a schema registry/evolution response for methods other than GET.
type InsertInfoDTO struct {
	Id      string `json:"identification"`
	Version int32  `json:"version"`
	Message string `json:"message"`
}
