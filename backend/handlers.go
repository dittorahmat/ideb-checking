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
		respondWithError(w, http.StatusBadRequest, "Request type is missing")
		return
	}

	var tableName string
	switch requestType {
	case "individual":
		tableName = "getDebtorExactIndividual"
	case "corporate":
		tableName = "getDebtorExactCorporate"
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid request type")
		return
	}

	a.getRequests(w, r, tableName)
}

func (a *App) createRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	switch req.SearchType {
	case SearchTypeInternal:
		if err := a.handleInternalSearch(&req); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal server error processing internal search")
			return
		}
	case SearchTypeLive:
		if err := a.handleLiveSearch(&req); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal server error processing live search")
			return
		}
	default:
		respondWithError(w, http.StatusBadRequest, "Invalid search type")
		return
	}

	respondWithJSON(w, http.StatusCreated, req)
}

func (a *App) handleInternalSearch(req *Request) error {
	inputData, err := a.readAndUnmarshalInputJSON(a.Config.InputJSONPath)
	if err != nil {
		return err
	}

	var nomorIdentitas string
	if len(inputData.Data.Corporate.CorporateDebtors) > 0 {
		nomorIdentitas = inputData.Data.Corporate.CorporateDebtors[0].TaxId
	}

	getIdebEntry := GetIdeb{
		NomorReferensiPengguna: inputData.Data.Header.UserReferenceCode,
		NomorIdentitas:         nomorIdentitas,
		Data:                   string(inputData.RawData),
	}

	if result := a.DB.Create(&getIdebEntry); result.Error != nil {
		return fmt.Errorf("Error saving to get_idebs table: %w", result.Error)
	}

	req.StatusAksi = "Selesai"
	if result := a.DB.Save(req); result.Error != nil {
		return fmt.Errorf("Error updating request status: %w", result.Error)
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

func (a *App) handleLiveSearch(req *Request) error {
	req.StatusAksi = "Dalam Proses"
	if result := a.DB.Create(req); result.Error != nil {
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
				NomorReferensiPengguna: inputData.Data.Header.UserReferenceCode,
				NomorIdentitas:         nomorIdentitas,
				Data:                   string(inputData.RawData),
			}
			a.DB.Create(&getIdebEntry)
		}
	}(req.ID)
	return nil
}

func (a *App) getRequests(w http.ResponseWriter, r *http.Request, tableName string) {
	var requests []Request
	if result := a.DB.Table(tableName).Find(&requests); result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving requests")
		return
	}

	respondWithJSON(w, http.StatusOK, requests)
}

func (a *App) generatePDFHandler(w http.ResponseWriter, r *http.Request) {
	// Get the request ID from the URL query parameter
	requestID := r.URL.Query().Get("id")
	if requestID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing request ID")
		return
	}

	// Retrieve the GetIdeb entry from the database
	var idebEntry GetIdeb
	if result := a.DB.Where("id = ?", requestID).First(&idebEntry); result.Error != nil {
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
