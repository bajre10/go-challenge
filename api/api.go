package api

import (
	"encoding/json"
	"net/http"
)

func SendJson(w http.ResponseWriter, statusCode int, item any) {
	w.Header().Set("Content-Type", "application/json")
	SendStatusCode(w, statusCode)

	err := json.NewEncoder(w).Encode(item)
	if err != nil {
		ThrowInternalError(w)
	}
}

func SendStatusCode(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
