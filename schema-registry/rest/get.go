// Copyright Syntio d.o.o.
// All Rights Reserved
//
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
