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

// Package rest contains the Schema registry REST Server configuration and start-up functions.
package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/syntio/schema-registry/business_logic"
	"github.com/syntio/schema-registry/util"
)

//
// GetSchemaByIdAndVersion is a GET function that expects parameters "id" and "version" for
// retrieving the schema from the underlying database.
//
// It currently writes back either:
//  - status 200 with a schema in JSON format, if the schema is registered
//  - status 404 with error message, if the schema is not registered.
//
func GetSchemaByIdAndVersion(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	schemaVersion := mux.Vars(r)["version"]

	version, err := util.StringToInt32(schemaVersion)
	if err != nil {
		writeInfoResponse(w, "Bad request. Version isn't a valid number.", http.StatusBadRequest)
		return
	}
	schemaInfo, found := service.GetSchema(r.Context(), id, version)
	if !found {
		writeInfoResponse(w, "Schema not found.", http.StatusNotFound)
		return
	}
	writeValidResponse(w, schemaInfo, http.StatusOK)
}
