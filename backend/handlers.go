package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"gorm.io/gorm"
)

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

	if req.SearchType == "internal" {
		// Read input.json
		byteValue, err := os.ReadFile("../memory-bank/input.json")
		if err != nil {
			http.Error(w, "Error reading input.json: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var inputData InputJSON
		json.Unmarshal(byteValue, &inputData)

		// Extract data for GetIdeb table
		var nomorIdentitas string
		if len(inputData.Data.Corporate.CorporateDebtors) > 0 {
			nomorIdentitas = inputData.Data.Corporate.CorporateDebtors[0].TaxId
		}

		getIdebEntry := GetIdeb{
			NomorReferensiPengguna: inputData.Data.Header.UserReferenceCode,
			NomorIdentitas:         nomorIdentitas,
			Data:                   string(byteValue),
		}

		if result := DB.Create(&getIdebEntry); result.Error != nil {
			http.Error(w, "Error saving to get_idebs table: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		req.StatusAksi = "Selesai"
		if result := DB.Save(&req); result.Error != nil {
			http.Error(w, "Error updating request status: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

	} else if req.SearchType == "live" {
		// Set initial status to "Dalam Proses"
		req.StatusAksi = "Dalam Proses"
		if result := DB.Create(&req); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Simulate asynchronous call to SLIK OJK in a goroutine
		go func(requestID uint) {
			// Simulate a delay for the OJK API call
			time.Sleep(5 * time.Second) // Simulate 5 seconds delay

			// After the simulated call, update the request status to "Selesai"
			var updatedReq Request
			if result := DB.First(&updatedReq, requestID); result.Error == nil {
				updatedReq.StatusAksi = "Selesai"
				DB.Save(&updatedReq)

				// Read input.json for dummy data
				byteValue, err := os.ReadFile("../memory-bank/input.json")
				if err != nil {
					log.Println("Error reading input.json for live simulation: ", err)
					return
				}

				var inputData InputJSON
				json.Unmarshal(byteValue, &inputData)

				var nomorIdentitas string
				if len(inputData.Data.Corporate.CorporateDebtors) > 0 {
					nomorIdentitas = inputData.Data.Corporate.CorporateDebtors[0].TaxId
				}

				getIdebEntry := GetIdeb{
					NomorReferensiPengguna: inputData.Data.Header.UserReferenceCode,
					NomorIdentitas:         nomorIdentitas,
					Data:                   string(byteValue),
				}
				DB.Create(&getIdebEntry)
			}
		}(req.ID)
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

func generatePDFHandler(w http.ResponseWriter, r *http.Request) {
	// Get the request ID from the URL query parameter
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		http.Error(w, "Missing request ID", http.StatusBadRequest)
		return
	}

	// Retrieve the GetIdeb entry from the database
	var idebEntry GetIdeb
	if result := DB.Where("id = ?", requestID).First(&idebEntry); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Record not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse the JSON data
	var inputData InputJSON
	if err := json.Unmarshal([]byte(idebEntry.Data), &inputData); err != nil {
		http.Error(w, "Error parsing JSON data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new Maroto instance
	m := maroto.New()

	// Add content to the PDF based on inputData and SARIPARI-PERTIWI-ABADI-b.pdf layout
	m.AddRow(10, col.New(12).Add(text.New("Laporan Informasi Debitur (iDeb)", props.Text{Align: align.Center, Size: 16, Style: fontstyle.Bold})))
	m.AddRow(7, col.New(12).Add(text.New("Otoritas Jasa Keuangan", props.Text{Align: align.Center, Size: 12})))
	m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))

	// Header Information
	m.AddRow(5, col.New(3).Add(text.New("Nama", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.CorporateKeyWord.IdentityNumberName, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("NPWP", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.CorporateDebtors[0].TaxId, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Tempat Pendirian", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.CorporateKeyWord.TestPlace, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Tanggal Akte Pendirian", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.CorporateDebtors[0].EstCertDate, props.Text{Size: 8})))

	m.AddRow(5, col.New(3).Add(text.New("Kode Ref. Pengguna", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Header.UserReferenceCode, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Nomor Laporan", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.ReportNumber, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Posisi Data Terakhir", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.LatestDataYearMonth, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Tanggal Permintaan", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.RequestDate, props.Text{Size: 8})))

	m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))
	m.AddRow(7, col.New(12).Add(text.New("DATA POKOK DEBITUR", props.Text{Align: align.Center, Size: 12, Style: fontstyle.Bold})))
	m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))

	// Corporate Debtor Information
	if len(inputData.Data.Corporate.CorporateDebtors) > 0 {
		debtor := inputData.Data.Corporate.CorporateDebtors[0]
		m.AddRow(5, col.New(3).Add(text.New("Nama Lengkap", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.FullName, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("NPWP", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.TaxId, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Bentuk BU / Go Public", props.Text{Size: 8})), col.New(9).Add(text.New(fmt.Sprintf("%s / %s", debtor.CompanyTypeDesc, debtor.GoPublicFlag), props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Tempat Pendirian", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.EstPlace, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("No / Tanggal Akta Pendirian", props.Text{Size: 8})), col.New(9).Add(text.New(fmt.Sprintf("%s / %s", debtor.EstCertNo, debtor.EstCertDate), props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Alamat", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.Address, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Kelurahan", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.SubDistrict, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Kecamatan", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.District, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Kabupaten / Kota", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.CityDesc, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Kode Pos", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.PostalCode, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Negara", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.CountryDesc, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Bidang Usaha", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.EconomicSectorDesc, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Pelapor / Tanggal Update", props.Text{Size: 8})), col.New(9).Add(text.New(fmt.Sprintf("%s / %s", debtor.MemberDesc, debtor.UpdatedDatetime), props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Peringkat / Tgl Pemeringkatan", props.Text{Size: 8})), col.New(9).Add(text.New(fmt.Sprintf("%s / %s", debtor.RatingDate, debtor.RatingDate), props.Text{Size: 8})))
		m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))

		// Pemilik / Pengurus
		if len(debtor.OfficisSharehldrsGroups) > 0 {
			m.AddRow(7, col.New(12).Add(text.New("Pemilik / Pengurus", props.Text{Align: align.Center, Size: 12, Style: fontstyle.Bold})))
			m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))

			for j, group := range debtor.OfficisSharehldrsGroups {
				m.AddRow(5, col.New(12).Add(text.New(fmt.Sprintf("Group %d: %s", j+1, group.MemberDesc), props.Text{Size: 9, Style: fontstyle.Bold})))
				for k, shareholder := range group.OfficisSharehldrs {
					m.AddRow(5, col.New(12).Add(text.New(fmt.Sprintf("  Shareholder %d: %s (%s)", k+1, shareholder.IdentityNumberName, shareholder.IdentityNumber), props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("    Jenis Kelamin", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.GenderDesc, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("    Jabatan", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.JobPositionDesc, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("    Pangsa Kepemilikan", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.ShareOwnership, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("    Alamat", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.Address, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("    Kabupaten / Kota", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.CityDesc, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("    Status Pengurus/Pemilik", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.ShareholderStatusDesc, props.Text{Size: 8})))
					m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))
				}
			}
		}
	}

	// Output the PDF
	pdf, err := m.Generate()
	if err != nil {
		http.Error(w, "Error generating PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"ideb_report.pdf\"")
	w.Write(pdf.GetBytes())
}