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
