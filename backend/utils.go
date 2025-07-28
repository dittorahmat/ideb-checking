package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// respondWithError sends a JSON error response with the given message and status code.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, APIError{Message: message, Code: code})
}

// respondWithJSON sends a JSON response with the given payload and status code.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
			// Log the error, but don't return an error to the client as the response header has already been written.
			// In a real application, you might want to use a more robust logging mechanism.
			fmt.Printf("Error writing response: %v\n", err)
		}
}
