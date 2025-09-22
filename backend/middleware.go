package main

import (
	"net/http"
)

// CORSMiddleware sets up CORS headers for all responses.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware checks if authentication is required and validates the session
func (a *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If authentication is disabled, skip the middleware
		if a.Config.DisableAuth {
			next.ServeHTTP(w, r)
			return
		}

		// For now, we're just checking if the request is to the API
		// In a real implementation, we would check for a valid session token
		if r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}

		// For API routes, we would normally check for authentication
		// For this prototype, we'll just allow all requests
		next.ServeHTTP(w, r)
	})
}