package main

import (
	"log"
	"net/http"
)

func main() {
	InitDatabase()

	RegisterRoutes()

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}