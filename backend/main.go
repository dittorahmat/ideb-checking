package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var frontendPath string

type Request struct {
	gorm.Model
	NomorReferensiPengguna       string    `json:"nomor_referensi_pengguna"`
	TujuanPenggunaan             string    `json:"tujuan_penggunaan"`
	JenisIdentitas               string    `json:"jenis_identitas"`
	NomorIdentitas               string    `json:"nomor_identitas"`
	PermintaanFasilitasOutstanding bool      `json:"permintaan_fasilitas_outstanding"`
	SearchType                   string    `json:"search_type"`
	StatusAksi                   string    `json:"status_aksi" gorm:"default:'Dalam Proses'"`
}

func (Request) TableName() string {
	return "getDebtorExactIndividual"
}

type CorporateRequest struct {
	gorm.Model
	NomorReferensiPengguna       string    `json:"nomor_referensi_pengguna"`
	TujuanPenggunaan             string    `json:"tujuan_penggunaan"`
	NomorIdentitas               string    `json:"nomor_identitas"`
	PermintaanFasilitasOutstanding bool      `json:"permintaan_fasilitas_outstanding"`
	SearchType                   string    `json:"search_type"`
	StatusAksi                   string    `json:"status_aksi" gorm:"default:'Dalam Proses'"`
}

func (CorporateRequest) TableName() string {
	return "getDebtorExactCorporate"
}

func main() {
	var err error
	DB, err = gorm.Open(sqlite.Open("./ideb.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&Request{}, &CorporateRequest{})

	// Serve static files from the "frontend" directory
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/requests", createRequestHandler)
	http.HandleFunc("/api/getDebtorExactIndividual", getDebtorExactIndividualHandler)
	http.HandleFunc("/api/getDebtorExactCorporate", getDebtorExactCorporateHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}



func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Dummy login handler
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"success"}`)
}

func getDebtorExactIndividualHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

	if r.Method == "GET" {
		getRequests(w, r, "getDebtorExactIndividual")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getDebtorExactCorporateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

	if r.Method == "GET" {
		getRequests(w, r, "getDebtorExactCorporate")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

	if r.Method == "POST" {
		createRequest(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createRequest(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := DB.Create(&req); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func getRequests(w http.ResponseWriter, r *http.Request, tableName string) {
	var requests []Request
	if result := DB.Table(tableName).Find(&requests); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
