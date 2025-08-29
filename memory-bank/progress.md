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
    - The "Input Permintaan IDeb" form submission now correctly triggers backend API calls.
    - **Event handling for dynamic content has been completely refactored and fixed.**
    - **A mock loading popup with a 3-second delay has been added to the form submission process for better UX.**
    - The "Daftar Permintaan IDeb" pages (`debitur-individual.html` and `badan-usaha.html`) now fetch and display data from the backend, and the "Lihat Detail" link triggers PDF generation.
- **Integration:** The frontend is successfully communicating with the backend API for all implemented features.
- **Testing:** Comprehensive unit and integration tests have been added for the backend Go application, including new unit tests for helper functions. All backend tests are now passing. The testing strategy has been updated in `testing-strategy.md`.

## 2. What's Left to Build
- **PDF Generation (Phase 3 Complete):** The backend now generates a PDF report with a refined layout and all shareholder details mapped from `input.json`.
- **Asynchronous OJK Queries:** The mechanism for simulating a "live" asynchronous call to SLIK OJK has been implemented, including a simulated delay and updating request status to "Selesai" with dummy data storage in `get_idebs` table.
- **UI/UX Refinements:** The frontend is basic and can be improved with better styling, loading indicators, and user feedback.

## 3. Current Status
- **Status:** Blocked by critical backend bug.
- **Details:** The application is currently non-functional for its primary purpose of creating new requests. While significant frontend fixes have been implemented, the backend crashes with a 500 Internal Server Error upon form submission. The immediate priority is to diagnose and fix this backend panic.

## 4. Known Issues
- **Critical:** The backend returns a `500 Internal Server Error` when the `/api/requests` endpoint is called. The exact cause is unknown as of now because the error message has not been successfully captured from the logs.
- The `DropTable` linting error is a false positive from `golangci-lint` as the method does not return an error.
