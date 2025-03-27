package api

import (
	"encoding/json"
	"net/http"
)

type RequestError struct {
	Code    int
	Message string
}

func throwError(w http.ResponseWriter, message string, statusCode int) {
	err := &RequestError{Code: statusCode, Message: message}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(err)
}

// Returns a user-friendly error status message
func ThrowRequestError(w http.ResponseWriter, message string, statusCode int) {
	throwError(w, message, statusCode)
}

// Returns generic error
func ThrowInternalError(w http.ResponseWriter) {
	throwError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
