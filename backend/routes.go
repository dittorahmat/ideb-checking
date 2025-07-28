package main

import "net/http"

func RegisterRoutes(inputJSONPath string) {
	// Serve static files from the "frontend" directory
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api/login", CORSMiddleware(loginHandler))
	http.HandleFunc("/api/requests", CORSMiddleware(func(w http.ResponseWriter, r *http.Request) {
		createRequestHandler(w, r)
	}))
	http.HandleFunc("/api/getDebtorExactIndividual", CORSMiddleware(getDebtorExactIndividualHandler))
	http.HandleFunc("/api/getDebtorExactCorporate", CORSMiddleware(getDebtorExactCorporateHandler))
	http.HandleFunc("/api/generate-pdf", CORSMiddleware(generatePDFHandler))
}
