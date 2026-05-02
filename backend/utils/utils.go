package utils

import (
    "encoding/json"
    "log"
    "net/http"
	"github.com/google/uuid"
)

// JSONResponse writes a JSON response with the given status code and data.
// It sets proper headers and handles encoding errors gracefully.
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    if err := json.NewEncoder(w).Encode(data); err != nil {
        // If encoding fails, log the error and send a generic response
        log.Printf("failed to encode JSON response: %v", err)
        http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
    }
}

// ErrorResponse logs the error and sends a JSON error response.
// It ensures consistent error formatting for clients.
func ErrorResponse(w http.ResponseWriter, status int, err error) {
    log.Printf("error: %v", err)

    response := map[string]string{
        "error": err.Error(),
    }

    JSONResponse(w, status, response)
}

// generateID creates and returns a new UUID string
func GenerateID() string {
    id := uuid.New()
    return id.String()
}
