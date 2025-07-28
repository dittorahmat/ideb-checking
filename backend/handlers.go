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
	var genericPayload GenericRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&genericPayload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Re-read the body for specific unmarshalling
	r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	var newReq interface{}

	switch genericPayload.RequestType {
	case "individual":
		var payload IndividualRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid individual request payload")
			return
		}
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
		var payload CorporateRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid corporate request payload")
			return
		}
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
		respondWithError(w, http.StatusBadRequest, "Invalid request type")
		return
	}

	// Handle search type based on the actual request object
	switch v := newReq.(type) {
	case Request:
		switch v.SearchType {
		case SearchTypeInternal:
			if err := a.handleInternalSearch(&v); err != nil {
				respondWithError(w, http.StatusInternalServerError, "Internal server error processing individual internal search")
				return
			}
		case SearchTypeLive:
			if err := a.handleLiveSearch(&v); err != nil {
				respondWithError(w, http.StatusInternalServerError, "Internal server error processing individual live search")
				return
			}
		default:
			respondWithError(w, http.StatusBadRequest, "Invalid search type for individual request")
			return
		}
	case CorporateRequest:
		switch v.SearchType {
		case SearchTypeInternal:
			if err := a.handleInternalSearch(&v); err != nil {
				respondWithError(w, http.StatusInternalServerError, "Internal server error processing corporate internal search")
				return
			}
		case SearchTypeLive:
			if err := a.handleLiveSearch(&v); err != nil {
				respondWithError(w, http.StatusInternalServerError, "Internal server error processing corporate live search")
				return
			}
		default:
			respondWithError(w, http.StatusBadRequest, "Invalid search type for corporate request")
			return
		}
	default:
		respondWithError(w, http.StatusInternalServerError, "Unknown request type")
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

	getIdebEntry := GetIdeb{
		NomorReferensiPengguna: inputData.Data.Header.UserReferenceCode,
		NomorIdentitas:         nomorIdentitas,
		Data:                   string(inputData.RawData),
	}

	if result := a.DB.Create(&getIdebEntry); result.Error != nil {
		return fmt.Errorf("Error saving to get_idebs table: %w", result.Error)
	}

	switch r := req.(type) {
	case *Request:
		r.StatusAksi = "Selesai"
		if result := a.DB.Save(r); result.Error != nil {
			return fmt.Errorf("Error updating individual request status: %w", result.Error)
		}
	case *CorporateRequest:
		r.StatusAksi = "Selesai"
		if result := a.DB.Save(r); result.Error != nil {
			return fmt.Errorf("Error updating corporate request status: %w", result.Error)
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
					NomorReferensiPengguna: inputData.Data.Header.UserReferenceCode,
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
					NomorReferensiPengguna: inputData.Data.Header.UserReferenceCode,
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
