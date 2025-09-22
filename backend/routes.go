package main

import (
	"net/http"
)

func (a *App) RegisterRoutes() {
	// API routes
	a.Router.HandleFunc("/api/login", a.loginHandler).Methods("POST", "OPTIONS")
	a.Router.HandleFunc("/api/requests", a.createRequestHandler).Methods("POST", "OPTIONS")
	a.Router.HandleFunc("/api/requests/{type}", a.getRequestsHandler).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/api/generate-pdf", a.generatePDFHandler).Methods("GET", "OPTIONS")

	// Serve static files from the "frontend" directory directly from the root
	fs := http.FileServer(http.Dir("../frontend"))
	a.Router.PathPrefix("/").Handler(fs)

	a.Router.Use(CORSMiddleware)
	a.Router.Use(a.AuthMiddleware)
}