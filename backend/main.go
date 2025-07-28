package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Determine the absolute path to input.json
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	executableDir := filepath.Dir(executablePath)
	inputJSONPath := filepath.Join(executableDir, "..", "memory-bank", "input.json")

	InitDatabase()

	RegisterRoutes(inputJSONPath)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
