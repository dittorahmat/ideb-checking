package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)



func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Dummy login handler
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (a *App) getRequestsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	if !ok {
		respondWithError(w, ErrInvalidInput.Code, ErrInvalidInput.Message)
		return
	}

	var tableName string
	switch requestType {
	case "individual":
		tableName = "getDebtorExactIndividual"
	case "corporate":
		tableName = "getDebtorExactCorporate"
	default:
		respondWithError(w, ErrInvalidInput.Code, ErrInvalidInput.Message)
		return
	}

	a.getRequests(w, r, tableName)
}

func (a *App) createRequestHandler(w http.ResponseWriter, r *http.Request) {
	var payload CombinedRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, ErrInvalidInput.Code, err.Error())
		return
	}

	var newReq interface{}

	switch payload.RequestType {
	case "individual":
		newReq = Request{
			BaseRequest: BaseRequest{
				NomorReferensiPengguna:         payload.NomorReferensiPengguna,
				TujuanPenggunaan:               payload.TujuanPenggunaan,
				NomorIdentitas:                 payload.NomorIdentitas,
				PermintaanFasilitasOutstanding: payload.PermintaanFasilitasOutstanding,
				SearchType:                     payload.SearchType,
			},
			JenisIdentitas: payload.JenisIdentitas,
		}
	case "corporate":
		newReq = CorporateRequest{
			BaseRequest: BaseRequest{
				NomorReferensiPengguna:         payload.NomorReferensiPengguna,
				TujuanPenggunaan:               payload.TujuanPenggunaan,
				NomorIdentitas:                 payload.NomorIdentitas,
				PermintaanFasilitasOutstanding: payload.PermintaanFasilitasOutstanding,
				SearchType:                     payload.SearchType,
			},
		}
	default:
		respondWithError(w, ErrInvalidInput.Code, "Invalid request_type")
		return
	}

	// Handle search type based on the actual request object
	switch v := newReq.(type) {
	case Request:
		switch v.SearchType {
		case SearchTypeInternal:
			if err := a.handleInternalSearch(&v); err != nil {
				respondWithError(w, ErrInternal.Code, err.Error())
				return
			}
		case SearchTypeLive:
			if err := a.handleLiveSearch(&v); err != nil {
				respondWithError(w, ErrInternal.Code, err.Error())
				return
			}
		default:
			respondWithError(w, ErrInvalidInput.Code, ErrInvalidInput.Message)
			return
		}
	case CorporateRequest:
		switch v.SearchType {
		case SearchTypeInternal:
			if err := a.handleInternalSearch(&v); err != nil {
				respondWithError(w, ErrInternal.Code, err.Error())
				return
			}
		case SearchTypeLive:
			if err := a.handleLiveSearch(&v); err != nil {
				respondWithError(w, ErrInternal.Code, err.Error())
				return
			}
		default:
			respondWithError(w, ErrInvalidInput.Code, ErrInvalidInput.Message)
			return
		}
	default:
		respondWithError(w, ErrInternal.Code, ErrInternal.Message)
		return
	}

	respondWithJSON(w, http.StatusCreated, newReq)
}

func (a *App) handleInternalSearch(req interface{}) error {
	inputData, err := a.readAndUnmarshalInputJSON(a.Config.InputJSONPath)
	if err != nil {
		return err
	}

	var nomorIdentitas string
	if len(inputData.Data.Corporate.CorporateDebtors) > 0 {
		nomorIdentitas = inputData.Data.Corporate.CorporateDebtors[0].TaxId
	}

	// Get the nomor_referensi_pengguna from the request object
	var nomorReferensiPengguna string
	switch r := req.(type) {
	case *Request:
		nomorReferensiPengguna = r.NomorReferensiPengguna
	case *CorporateRequest:
		nomorReferensiPengguna = r.NomorReferensiPengguna
	}

	getIdebEntry := GetIdeb{
		NomorReferensiPengguna: nomorReferensiPengguna, // Use the request's nomor_referensi_pengguna
		NomorIdentitas:         nomorIdentitas,
		Data:                   string(inputData.RawData),
	}

	if result := a.DB.Create(&getIdebEntry); result.Error != nil {
		return fmt.Errorf("Error saving to get_idebs table: %w", result.Error)
	}

	switch r := req.(type) {
	case *Request:
		r.StatusAksi = "Selesai"
		if result := a.DB.Create(r); result.Error != nil {
			return fmt.Errorf("Error creating individual request: %w", result.Error)
		}
	case *CorporateRequest:
		r.StatusAksi = "Selesai"
		if result := a.DB.Create(r); result.Error != nil {
			return fmt.Errorf("Error creating corporate request: %w", result.Error)
		}
	default:
		return fmt.Errorf("Unsupported request type for handleInternalSearch")
	}
	return nil
}

func (a *App) readAndUnmarshalInputJSON(filePath string) (*InputJSON, error) {
	byteValue, err := a.ReadFileFunc(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading input.json: %w", err)
	}

	var inputData InputJSON
	if err := json.Unmarshal(byteValue, &inputData); err != nil {
		return nil, fmt.Errorf("Error unmarshalling input.json: %w", err)
	}
	inputData.RawData = byteValue // Store the raw data

	return &inputData, nil
}

func (a *App) handleLiveSearch(req interface{}) error {
	switch r := req.(type) {
	case *Request:
		r.StatusAksi = "Dalam Proses"
		if result := a.DB.Create(r); result.Error != nil {
			return result.Error
		}

		go func(requestID uint) {
			time.Sleep(5 * time.Second)

			var updatedReq Request
			if result := a.DB.First(&updatedReq, requestID); result.Error == nil {
				updatedReq.StatusAksi = "Selesai"
				a.DB.Save(&updatedReq)

				inputData, err := a.readAndUnmarshalInputJSON(a.Config.InputJSONPath)
				if err != nil {
					log.Println("Error reading and unmarshalling input.json for live simulation: ", err)
					return
				}

				var nomorIdentitas string
				if len(inputData.Data.Corporate.CorporateDebtors) > 0 {
					nomorIdentitas = inputData.Data.Corporate.CorporateDebtors[0].TaxId
				}

				getIdebEntry := GetIdeb{
					NomorReferensiPengguna: updatedReq.NomorReferensiPengguna, // Use the request's nomor_referensi_pengguna
					NomorIdentitas:         nomorIdentitas,
					Data:                   string(inputData.RawData),
				}
				a.DB.Create(&getIdebEntry)
			}
		}(r.ID)
	case *CorporateRequest:
		r.StatusAksi = "Dalam Proses"
		if result := a.DB.Create(r); result.Error != nil {
			return result.Error
		}

		go func(requestID uint) {
			time.Sleep(5 * time.Second)

			var updatedReq CorporateRequest
			if result := a.DB.First(&updatedReq, requestID); result.Error == nil {
				updatedReq.StatusAksi = "Selesai"
				a.DB.Save(&updatedReq)

				inputData, err := a.readAndUnmarshalInputJSON(a.Config.InputJSONPath)
				if err != nil {
					log.Println("Error reading and unmarshalling input.json for live simulation: ", err)
					return
				}

				var nomorIdentitas string
				if len(inputData.Data.Corporate.CorporateDebtors) > 0 {
					nomorIdentitas = inputData.Data.Corporate.CorporateDebtors[0].TaxId
				}

				getIdebEntry := GetIdeb{
					NomorReferensiPengguna: updatedReq.NomorReferensiPengguna, // Use the request's nomor_referensi_pengguna
					NomorIdentitas:         nomorIdentitas,
					Data:                   string(inputData.RawData),
				}
				a.DB.Create(&getIdebEntry)
			}
		}(r.ID)
	default:
		return fmt.Errorf("Unsupported request type for handleLiveSearch")
	}
	return nil
}

func (a *App) getRequests(w http.ResponseWriter, r *http.Request, tableName string) {
	var results interface{}
	var err error

	switch tableName {
	case "getDebtorExactIndividual":
		var requests []Request
		if result := a.DB.Table(tableName).Find(&requests); result.Error != nil {
			err = result.Error
		}
		results = requests
	case "getDebtorExactCorporate":
		var requests []CorporateRequest
		if result := a.DB.Table(tableName).Find(&requests); result.Error != nil {
			err = result.Error
		}
		results = requests
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid table name")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving requests: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

func (a *App) generatePDFHandler(w http.ResponseWriter, r *http.Request) {
	// Get the request ID from the URL query parameter
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing request ID")
		return
	}

	var request CorporateRequest
	if result := a.DB.First(&request, requestID); result.Error != nil {
		respondWithError(w, http.StatusNotFound, "Request not found")
		return
	}

	// Retrieve the GetIdeb entry from the database
	var idebEntry GetIdeb
	if result := a.DB.Where("nomor_referensi_pengguna = ?", request.NomorReferensiPengguna).First(&idebEntry); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			respondWithError(w, http.StatusNotFound, "Record not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Database error: "+result.Error.Error())
		}
		return
	}

	// Parse the JSON data
	var inputData InputJSON
	if err := json.Unmarshal([]byte(idebEntry.Data), &inputData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing JSON data")
		return
	}

	// Create a new Maroto instance
	pdfBytes, err := generateIdebPDF(&inputData)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Error generating PDF")
				return
			}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"ideb_report.pdf\"")
	_, err = w.Write(pdfBytes)
	if err != nil {
		log.Printf("Error writing PDF to response: %v", err)
	}
}
