package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kalshiev/lucien/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// RespondJSON responds with JSON
func RespondJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(data)
}

// RespondError responds with an error
func RespondError(w http.ResponseWriter, code int, msg string) {
	RespondJSON(w, code, models.ErrorResponse{Error: msg})
}

// RespondMessage responds with a message
func RespondMessage(w http.ResponseWriter, code int, msg string) {
	RespondJSON(w, code, models.MessageResponse{Message: msg})
}

// ParseUUID parses and validates a UUID from a string
func ParseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, models.BadRequest("Invalid UUID format")
	}
	return id, nil
}

// GetPathParam gets a path parameter from chi
func GetPathParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
