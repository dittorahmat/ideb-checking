# Progress

## 1. What Works
- **Project Scaffolding:** `frontend` and `backend` directories are created.
- **Database:** An SQLite database (`ideb.db`) has been created with a `requests` table. The `schema.sql` file has been removed as GORM handles schema migration directly.
- **Backend API (Go):**
    - A basic Go web server is running.
    - Dummy login endpoint (`/api/login`) is functional.
    - API endpoint (`/api/requests`) for creating and listing requests is implemented using GORM and connected to the database.
    - **Refactoring:** `main.go` has been refactored into `database.go` (for DB initialization), `models.go` (for data structures), `handlers.go` (for HTTP handlers), and `routes.go` (for route registration) to improve modularity and maintainability. CORS middleware has been implemented, `createRequest` has been refactored, and request structs have been unified. Additionally, `handlers.go` has been refactored to extract JSON reading/unmarshalling into a helper function and PDF generation logic into a new `pdf.go` file, significantly improving modularity and maintainability. `readFileFunc` has been refactored to be a field of the `App` struct, and error handling has been improved with a custom `APIError` struct and generic error messages.
- **Frontend (HTML/JS/CSS):**
    - A basic user interface with a sidebar and content area is in place.
    - The default landing page after login is now `input-permintaan-badan-usaha.html`.
    - A dummy login form is functional.
    - The "Input Permintaan IDeb" form can be submitted to the backend, with "internal" search type triggering data ingestion.
    - The "Daftar Permintaan IDeb" pages (`debitur-individual.html` and `badan-usaha.html`) now fetch and display data from the backend, and the "Lihat Detail" link triggers PDF generation.
- **Integration:** The frontend is successfully communicating with the backend API for all implemented features.
- **Testing:** Comprehensive unit and integration tests have been added for the backend Go application, including new unit tests for helper functions. All backend tests are now passing. The testing strategy has been updated in `testing-strategy.md`.

## 2. What's Left to Build
- **PDF Generation (Phase 3 Complete):** The backend now generates a PDF report with a refined layout and all shareholder details mapped from `input.json`.
- **Asynchronous OJK Queries:** The mechanism for simulating a "live" asynchronous call to SLIK OJK has been implemented, including a simulated delay and updating request status to "Selesai" with dummy data storage in `get_idebs` table.

- **UI/UX Refinements:** The frontend is basic and can be improved with better styling, loading indicators, and user feedback.

## 3. Current Status
- **Status:** v0 Mockup - Backend Refactoring Complete. All tests passing.
- **Details:** The core functionality of the mockup is implemented. Significant refactoring has been done on the backend for dependency injection and error handling. All backend tests are now passing. The application is running locally with the backend server serving the frontend files.

## 4. Known Issues
- The `DropTable` linting error is a false positive from `golangci-lint` as the method does not return an error.