package main

import "net/http"

func RegisterRoutes() {
	// Serve static files from the "frontend" directory
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/requests", createRequestHandler)
	http.HandleFunc("/api/getDebtorExactIndividual", getDebtorExactIndividualHandler)
	http.HandleFunc("/api/getDebtorExactCorporate", getDebtorExactCorporateHandler)
	http.HandleFunc("/api/generate-pdf", generatePDFHandler)
}