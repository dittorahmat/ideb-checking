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

// handleSimilarSearch handles the similar search requests for corporate entities
func (a *App) handleSimilarSearch(w http.ResponseWriter, r *http.Request) {
	// Read the similar.json file
	similarData, err := a.ReadFileFunc("../memory-bank/similar.json")
	if err != nil {
		log.Printf("Error reading similar.json: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error reading similar data")
		return
	}

	// Parse the JSON data
	var inputData interface{}
	if err := json.Unmarshal(similarData, &inputData); err != nil {
		log.Printf("Error parsing similar.json: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error parsing similar data")
		return
	}

	// Respond with the data
	respondWithJSON(w, http.StatusOK, inputData)
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

// handleSimilarSearchPDF generates a PDF from the similar search result data directly
func (a *App) handleSimilarSearchPDF(w http.ResponseWriter, r *http.Request) {
	// Parse the similar search result from the request body
	var result SimilarSearchResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create InputJSON structure from the similar search result
	// For this mockup, we'll create a minimal InputJSON based on the result
	inputData := &InputJSON{
		Code:   "200",
		Status: "success",
		Remark: "Berhasil",
		Data: struct {
			Header    Header    `json:"header"`
			Corporate Corporate `json:"corporate"`
		}{
			Header: Header{
				UserReferenceCode:        result.NomorIdentitas, // Using nomor identitas as reference
				ResultDate:               time.Now().Format("20060102"),
				InquiryId:                result.KodePelapor,
				InquiryUserId:            "user123", // Default value
				InquiryCreatedBy:         "system",  // Default value
				InquiryMemberCode:        result.KodePelapor,
				InquiryOfficeCode:        "001", // Default value
				ReportRequestPurposeCode: "01",  // Default value
				InquiryDate:              time.Now().Format("20060102"),
				DataSetTotal:             "1",   // Since we're dealing with single results
				DataSetNumber:            "1",   // Default value
			},
			Corporate: Corporate{
				ReportNumber:        "RPT001", // Default value
				LatestDataYearMonth: time.Now().Format("200601"),
				RequestDate:         time.Now().Format("20060102"),
				CorporateKeyWord: struct {
					IdentityNumberName string `json:"identityNumberName"`
					TestPlace          string `json:"testPlace"`
					RecordStatusFlag   string `json:"recordStatusFlag"`
				}{
					IdentityNumberName: result.NamaDebitur,
					TestPlace:          result.TempatPendirian,
					RecordStatusFlag:   "A", // Default value
				},
				CorporateDebtors: []CorporateDebtor{
					{
						IdentityNumberName: result.NamaDebitur,
						FullName:           result.NamaDebitur,
						TaxId:              result.NomorIdentitas,
						CompanyType:        result.KodeJenisIdentitas, // Using this field as company type
						CompanyTypeDesc:    "PT",                    // Default value
						EstPlace:           result.TempatPendirian,
						EstCertNo:          "",                      // Default value
						EstCertDate:        result.TanggalPendirian,
						Member:             result.KodePelapor,
						MemberDesc:         "Default Member Desc", // Default value
						UpdatedDatetime:    time.Now().Format("2006-01-02 15:04:05"),
						Address:            result.AlamatDebitur,
						SubDistrict:        "", // Default value
						District:           "", // Default value
						City:               "", // Default value
						CityDesc:           "", // Default value
						PostalCode:         result.KodePos,
						Country:            "ID", // Default value
						CountryDesc:        "Indonesia", // Default value
						LatestAddCertNo:    "",          // Default value
						LatestAddCertDate:  result.TanggalPendirian,
						EconomicSector:     "", // Default value
						EconomicSectorDesc: "", // Default value
						RatingDate:         "", // Default value
						CreatedDatetime:    time.Now().Format("2006-01-02 15:04:05"),
						GoPublicFlag:       "0", // Default value
						OfficisSharehldrsGroups: []OfficisSharehldrsGroup{}, // Empty for now
					},
				},
			},
		},
	}

	// Marshal the structure to get the raw data
	rawData, err := json.Marshal(inputData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating input data")
		return
	}
	inputData.RawData = rawData

	// Generate PDF using the input data
	pdfBytes, err := generateIdebPDF(inputData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating PDF")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"similar_search_result.pdf\"")
	_, err = w.Write(pdfBytes)
	if err != nil {
		log.Printf("Error writing PDF to response: %v", err)
	}
}
