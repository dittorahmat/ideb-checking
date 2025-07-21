package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var frontendPath string

type Request struct {
	ID                           int       `json:"id"`
	NomorReferensiPengguna       string    `json:"nomor_referensi_pengguna"`
	TujuanPenggunaan             string    `json:"tujuan_penggunaan"`
	JenisIdentitas               string    `json:"jenis_identitas"`
	NomorIdentitas               string    `json:"nomor_identitas"`
	PermintaanFasilitasOutstanding bool      `json:"permintaan_fasilitas_outstanding"`
	SearchType                   string    `json:"search_type"`
	StatusAksi                   string    `json:"status_aksi"`
	CreatedAt                    time.Time `json:"created_at"`
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./ideb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()

	// Serve static files from the "frontend" directory
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/requests", requestsHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func createTable() {
	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS requests (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        nomor_referensi_pengguna TEXT NOT NULL,
        tujuan_penggunaan TEXT NOT NULL,
        jenis_identitas TEXT NOT NULL,
        nomor_identitas TEXT NOT NULL,
        permintaan_fasilitas_outstanding BOOLEAN NOT NULL,
        search_type TEXT NOT NULL,
        status_aksi TEXT NOT NULL DEFAULT 'Dalam Proses',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Dummy login handler
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"success"}`)
}

func requestsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

	switch r.Method {
	case "GET":
		getRequests(w, r)
	case "POST":
		createRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createRequest(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statement, err := db.Prepare(`INSERT INTO requests 
		(nomor_referensi_pengguna, tujuan_penggunaan, jenis_identitas, nomor_identitas, permintaan_fasilitas_outstanding, search_type) 
		VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = statement.Exec(req.NomorReferensiPengguna, req.TujuanPenggunaan, req.JenisIdentitas, req.NomorIdentitas, req.PermintaanFasilitasOutstanding, req.SearchType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func getRequests(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nomor_referensi_pengguna, tujuan_penggunaan, jenis_identitas, nomor_identitas, permintaan_fasilitas_outstanding, search_type, status_aksi, created_at FROM requests")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []Request
	for rows.Next() {
		var req Request
		if err := rows.Scan(&req.ID, &req.NomorReferensiPengguna, &req.TujuanPenggunaan, &req.JenisIdentitas, &req.NomorIdentitas, &req.PermintaanFasilitasOutstanding, &req.SearchType, &req.StatusAksi, &req.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		requests = append(requests, req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
