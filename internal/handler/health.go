package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the response format of the health check endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health is the HTTP handler for the health check endpoint.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(HealthResponse{Status: "ok"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
