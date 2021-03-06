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

package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	service "github.com/syntio/schema-registry/business_logic"
)

// BackwardResolver retrieves a list of schemas that have the same ID.
func BackwardResolver(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		writeInfoResponse(w, "Id non existent", 400)
	}
	schemaVersions, err := service.ListSchemas(r.Context(), id)
	if err != nil {
		writeInfoResponse(w, "Server storage error while getting schema versions.", http.StatusInternalServerError)
		return
	}
	writeValidResponse(w, schemaVersions, http.StatusOK)
}
