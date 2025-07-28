package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Drop existing tables to ensure a clean state for each test
	db.Migrator().DropTable(&Request{}, &CorporateRequest{}, &GetIdeb{})

	if err := db.AutoMigrate(&Request{}, &CorporateRequest{}, &GetIdeb{}); err != nil {
		panic(fmt.Sprintf("Failed to auto migrate test database: %v", err))
	}

	return db
}

func TestReadAndUnmarshalInputJSON(t *testing.T) {
	// Test case 1: Valid JSON file
	dummyInputJSON := []byte(`{"code": "000","status": "SUCCESS","remark": "SUCCESS","data": {"header": {"userReferenceCode": "TESTREF123"},"corporate": {"corporateDebtors": [{"taxId": "123456789"}]}}}`)
	tempFile, err := os.CreateTemp("", "input_valid_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	_, err = tempFile.Write(dummyInputJSON)
	assert.NoError(t, err)
	tempFile.Close()

	mockReadFileFunc := func(name string) ([]byte, error) {
		return dummyInputJSON, nil
	}

	app := &App{ReadFileFunc: mockReadFileFunc}
	inputData, err := app.readAndUnmarshalInputJSON(tempFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, inputData)
	assert.Equal(t, "TESTREF123", inputData.Data.Header.UserReferenceCode)
	assert.Equal(t, "123456789", inputData.Data.Corporate.CorporateDebtors[0].TaxId)
	assert.Equal(t, dummyInputJSON, inputData.RawData)

	// Test case 2: Invalid JSON file
	invalidJSON := []byte(`invalid json`)
	tempFileInvalid, err := os.CreateTemp("", "input_invalid_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFileInvalid.Name())
	_, err = tempFileInvalid.Write(invalidJSON)
	assert.NoError(t, err)
	tempFileInvalid.Close()

	mockReadFileFunc = func(name string) ([]byte, error) {
		return invalidJSON, nil
	}
	app.ReadFileFunc = mockReadFileFunc

	inputData, err = app.readAndUnmarshalInputJSON(tempFileInvalid.Name())
	assert.Error(t, err)
	assert.Nil(t, inputData)
	assert.Contains(t, err.Error(), "Error unmarshalling input.json")

	// Test case 3: File not found
	mockReadFileFunc = func(name string) ([]byte, error) {
		return nil, os.ErrNotExist
	}
	app.ReadFileFunc = mockReadFileFunc

	inputData, err = app.readAndUnmarshalInputJSON("non_existent_file.json")
	assert.Error(t, err)
	assert.Nil(t, inputData)
	assert.Contains(t, err.Error(), "Error reading input.json")
}

func TestHandleInternalSearch(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create a dummy input.json content
	dummyInputJSON := []byte(`{"code": "000","status": "SUCCESS","remark": "SUCCESS","data": {"header": {"userReferenceCode": "UNITTESTREF123"},"corporate": {"corporateDebtors": [{"taxId": "987654321"}]}}}`)

	// Mock readFileFunc
	mockReadFileFunc := func(name string) ([]byte, error) {
		return dummyInputJSON, nil
	}

	// Create App instance with test DB and config
	app := &App{
		DB: db,
		Config: &Config{
			InputJSONPath: "dummy_path.json", // Path doesn't matter as readFileFunc is mocked
		},
		ReadFileFunc: mockReadFileFunc,
	}

	// Create a Request object
	req := &Request{
		BaseRequest: BaseRequest{
			NomorReferensiPengguna: "UNITTESTREF123",
			SearchType:             SearchTypeInternal,
		},
	}

	// Call the function under test
	err := app.handleInternalSearch(req)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, "Selesai", req.StatusAksi)

	var savedRequest Request
	db.First(&savedRequest, req.ID)
	assert.Equal(t, "Selesai", savedRequest.StatusAksi)

	var savedGetIdeb GetIdeb
	db.Where("nomor_referensi_pengguna = ?", "UNITTESTREF123").First(&savedGetIdeb)
	assert.Equal(t, "987654321", savedGetIdeb.NomorIdentitas)
	assert.Equal(t, string(dummyInputJSON), savedGetIdeb.Data)
}

func TestCreateRequestHandler_InternalSearch(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create a dummy input.json file
	tempFile, err := os.CreateTemp("", "input_internal_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up the dummy file

	dummyInputJSON := []byte(`{"code": "000","status": "SUCCESS","remark": "SUCCESS","data": {"header": {"userReferenceCode": "TESTREF123"},"corporate": {"corporateDebtors": [{"taxId": "123456789"}]}}}`)
	_, err = tempFile.Write(dummyInputJSON)
	assert.NoError(t, err)
	tempFile.Close()

	inputJSONPath := tempFile.Name()

	// Mock readFileFunc
	mockReadFileFunc := func(name string) ([]byte, error) {
		if name == inputJSONPath {
			return dummyInputJSON, nil
		}
		return os.ReadFile(name) // Fallback to actual file read for other cases if needed
	}

	// Create App instance with test DB and config
	app := &App{
		DB: db,
		Config: &Config{
			InputJSONPath: inputJSONPath,
		},
		ReadFileFunc: mockReadFileFunc,
	}

	// Create a request body
	requestBody := map[string]interface{}{
		"nomor_referensi_pengguna":         "TESTREF123",
		"tujuan_penggunaan":                "Test Tujuan",
		"jenis_identitas":                  "KTP",
		"nomor_identitas":                  "12345",
		"permintaan_fasilitas_outstanding": true,
		"search_type":                      SearchTypeInternal,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Create a request
	req, err := http.NewRequest("POST", "/api/requests", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.createRequestHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response Request
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Selesai", response.StatusAksi)

	// Verify database entry
	var savedRequest Request
	db.First(&savedRequest, response.ID)
	assert.Equal(t, "Selesai", savedRequest.StatusAksi)

	var savedGetIdeb GetIdeb
	db.Where("nomor_referensi_pengguna = ?", "TESTREF123").First(&savedGetIdeb)
	assert.Equal(t, "123456789", savedGetIdeb.NomorIdentitas)
}

func TestHandleLiveSearch(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create a dummy input.json content
	dummyInputJSON := []byte(`{"code": "000","status": "SUCCESS","remark": "SUCCESS","data": {"header": {"userReferenceCode": "LIVETESTREF456"},"corporate": {"corporateDebtors": [{"taxId": "112233445"}]}}}`)

	// Mock readFileFunc
	mockReadFileFunc := func(name string) ([]byte, error) {
		return dummyInputJSON, nil
	}

	// Create App instance with test DB and config
	app := &App{
		DB: db,
		Config: &Config{
			InputJSONPath: "dummy_path.json", // Path doesn't matter as readFileFunc is mocked
		},
		ReadFileFunc: mockReadFileFunc,
	}

	// Create a Request object
	req := &Request{
		BaseRequest: BaseRequest{
			NomorReferensiPengguna: "LIVETESTREF456",
			SearchType:             SearchTypeLive,
		},
	}

	// Call the function under test
	err := app.handleLiveSearch(req)
	assert.NoError(t, err)

	// Assertions immediately after calling handleLiveSearch
	assert.Equal(t, "Dalam Proses", req.StatusAksi)

	// Verify database entry immediately after initial request
	var savedRequest Request
	db.First(&savedRequest, req.ID)
	assert.Equal(t, "Dalam Proses", savedRequest.StatusAksi)

	// Wait for the goroutine to complete (simulated delay + DB update)
	// In a real scenario, you might use a channel or a more sophisticated wait mechanism
	// For this simple test, a short sleep is acceptable.
	time.Sleep(6 * time.Second)

	// Verify database entry after goroutine completes
	db.First(&savedRequest, req.ID)
	assert.Equal(t, "Selesai", savedRequest.StatusAksi)

	var savedGetIdeb GetIdeb
	db.Where("nomor_referensi_pengguna = ?", "LIVETESTREF456").First(&savedGetIdeb)
	assert.Equal(t, "112233445", savedGetIdeb.NomorIdentitas)
	assert.Equal(t, string(dummyInputJSON), savedGetIdeb.Data)
}

func TestCreateRequestHandler_LiveSearch(t *testing.T) {
	// Setup test database
	db := setupTestDB()

	// Create a dummy input.json file
	tempFile, err := os.CreateTemp("", "input_live_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up the dummy file

	dummyInputJSON := []byte(`{"code": "000","status": "SUCCESS","remark": "SUCCESS","data": {"header": {"userReferenceCode": "LIVETEST123"},"corporate": {"corporateDebtors": [{"taxId": "67890"}]}}}`)
	_, err = tempFile.Write(dummyInputJSON)
	assert.NoError(t, err)
	tempFile.Close()

	inputJSONPath := tempFile.Name()

	// Mock readFileFunc
	mockReadFileFunc := func(name string) ([]byte, error) {
		if name == inputJSONPath {
			return dummyInputJSON, nil
		}
		return os.ReadFile(name) // Fallback to actual file read for other cases if needed
	}

	// Create App instance with test DB and config
	app := &App{
		DB: db,
		Config: &Config{
			InputJSONPath: inputJSONPath,
		},
		ReadFileFunc: mockReadFileFunc,
	}

	// Create a request body
	requestBody := map[string]interface{}{
		"nomor_referensi_pengguna":         "LIVETEST123",
		"tujuan_penggunaan":                "Live Test Tujuan",
		"jenis_identitas":                  "KTP",
		"nomor_identitas":                  "67890",
		"permintaan_fasilitas_outstanding": false,
		"search_type":                      SearchTypeLive,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Create a request
	req, err := http.NewRequest("POST", "/api/requests", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.createRequestHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response Request
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Dalam Proses", response.StatusAksi)

	// Verify database entry immediately after initial request
	var savedRequest Request
	db.First(&savedRequest, response.ID)
	assert.Equal(t, "Dalam Proses", savedRequest.StatusAksi)

	// Wait for the goroutine to complete (simulated delay + DB update)
	// In a real scenario, you might use a channel or a more sophisticated wait mechanism
	// For this simple test, a short sleep is acceptable.
	time.Sleep(6 * time.Second)

	// Verify database entry after goroutine completes
	db.First(&savedRequest, response.ID)
	assert.Equal(t, "Selesai", savedRequest.StatusAksi)

	var savedGetIdeb GetIdeb
	db.Where("nomor_referensi_pengguna = ?", "LIVETEST123").First(&savedGetIdeb)
	assert.Equal(t, "67890", savedGetIdeb.NomorIdentitas)
}

func TestRespondWithError(t *testing.T) {
	rr := httptest.NewRecorder()
	respondWithError(rr, http.StatusBadRequest, "Test Error Message")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var apiError APIError
	err := json.Unmarshal(rr.Body.Bytes(), &apiError)
	assert.NoError(t, err)
	assert.Equal(t, "Test Error Message", apiError.Message)
	assert.Equal(t, http.StatusBadRequest, apiError.Code)
}

func TestRespondWithJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]string{"message": "Success"}
	respondWithJSON(rr, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, `{"message":"Success"}`, rr.Body.String())
}

func TestLoginHandler(t *testing.T) {
	app := &App{}
	req, err := http.NewRequest("POST", "/api/login", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.loginHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"status":"success"}`, rr.Body.String())
}

func TestGetDebtorExactIndividualHandler(t *testing.T) {
	db := setupTestDB()
	app := &App{DB: db}

	// Insert dummy data
	db.Create(&Request{
		BaseRequest: BaseRequest{
			NomorReferensiPengguna:         "INDIVIDUAL1",
			TujuanPenggunaan:               "Kredit",
			NomorIdentitas:                 "12345",
			PermintaanFasilitasOutstanding: true,
			SearchType:                     SearchTypeInternal,
			StatusAksi:                     "Selesai",
		},
		JenisIdentitas: "KTP",
	})
	db.Create(&Request{
		BaseRequest: BaseRequest{
			NomorReferensiPengguna:         "INDIVIDUAL2",
			TujuanPenggunaan:               "Kredit",
			NomorIdentitas:                 "67890",
			PermintaanFasilitasOutstanding: false,
			SearchType:                     SearchTypeLive,
			StatusAksi:                     "Dalam Proses",
		},
		JenisIdentitas: "SIM",
	})

	req, err := http.NewRequest("GET", "/api/requests/individual", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/requests/{type}", app.getRequestsHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var requests []Request
	err = json.Unmarshal(rr.Body.Bytes(), &requests)
	assert.NoError(t, err)
	assert.Len(t, requests, 2)
	assert.Equal(t, "INDIVIDUAL1", requests[0].NomorReferensiPengguna)
	assert.Equal(t, "INDIVIDUAL2", requests[1].NomorReferensiPengguna)
}

func TestGetDebtorExactCorporateHandler(t *testing.T) {
	db := setupTestDB()
	app := &App{DB: db}

	// Insert dummy data
	db.Create(&CorporateRequest{
		BaseRequest: BaseRequest{
			NomorReferensiPengguna:         "CORPORATE1",
			TujuanPenggunaan:               "Investasi",
			NomorIdentitas:                 "112233",
			PermintaanFasilitasOutstanding: true,
			SearchType:                     SearchTypeInternal,
			StatusAksi:                     "Selesai",
		},
	})
	db.Create(&CorporateRequest{
		BaseRequest: BaseRequest{
			NomorReferensiPengguna:         "CORPORATE2",
			TujuanPenggunaan:               "Modal Kerja",
			NomorIdentitas:                 "445566",
			PermintaanFasilitasOutstanding: false,
			SearchType:                     SearchTypeLive,
			StatusAksi:                     "Dalam Proses",
		},
	})

	req, err := http.NewRequest("GET", "/api/requests/corporate", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/requests/{type}", app.getRequestsHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var requests []CorporateRequest
	err = json.Unmarshal(rr.Body.Bytes(), &requests)
	assert.NoError(t, err)
	assert.Len(t, requests, 2)
	assert.Equal(t, "CORPORATE1", requests[0].NomorReferensiPengguna)
	assert.Equal(t, "CORPORATE2", requests[1].NomorReferensiPengguna)
}

func TestGeneratePDFHandler(t *testing.T) {
	db := setupTestDB()
	app := &App{DB: db}

	// Insert a dummy GetIdeb entry
	dummyJSONData := `{"code": "000","status": "SUCCESS","remark": "SUCCESS","data": {"header": {"userReferenceCode": "PDFTEST123"},"corporate": {"reportNumber": "REPORT123","latestDataYearMonth": "2023-12","requestDate": "2023-01-01","corporateKeyWord": {"identityNumberName": "PT TEST JAYA","testPlace": "Jakarta"},"corporateDebtors": [{"fullName": "John Doe","taxId": "1234567890","companyTypeDesc": "PT","goPublicFlag": "T","estPlace": "Bandung","estCertNo": "ACT123","estCertDate": "2020-01-01","address": "Jl. Contoh No. 1","subDistrict": "Kel. Contoh","district": "Kec. Contoh","cityDesc": "Jakarta","postalCode": "12345","countryDesc": "Indonesia","memberDesc": "Bank A","updatedDatetime": "2023-01-01","economicSectorDesc": "Perdagangan","ratingDate": "2022-01-01","officisSharehldrsGroups": [{"memberDesc": "Group A","officisSharehldrs": [{"identityNumberName": "Shareholder 1","identityNumber": "SH123","genderDesc": "Laki-laki","jobPositionDesc": "Direktur","shareOwnership": "50","address": "Jl. Saham No. 1","subDistrict": "Kel. Saham","district": "Kec. Saham","cityDesc": "Surabaya","shareholderStatusDesc": "Aktif"}]}]}]}}}`

	getIdebEntry := GetIdeb{
		NomorReferensiPengguna: "PDFTEST123",
		NomorIdentitas:         "1234567890",
		Data:                   dummyJSONData,
	}
	db.Create(&getIdebEntry)

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/generate-pdf?id=%d", getIdebEntry.ID), nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/generate-pdf", app.generatePDFHandler).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/pdf", rr.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=\"ideb_report.pdf\"", rr.Header().Get("Content-Disposition"))
}