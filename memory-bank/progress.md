# Progress

## 1. What Works
- **Project Scaffolding:** `frontend` and `backend` directories are created.
- **Database:** An SQLite database (`ideb.db`) has been created with a `requests` table.
- **Backend API (Go):**
    - A basic Go web server is running.
    - Dummy login endpoint (`/api/login`) is functional.
    - API endpoint (`/api/requests`) for creating and listing requests is implemented using GORM and connected to the database.
    - **Refactoring:** `main.go` has been refactored into `database.go` (for DB initialization), `models.go` (for data structures), `handlers.go` (for HTTP handlers), and `routes.go` (for route registration) to improve modularity and maintainability.
- **Frontend (HTML/JS/CSS):**
    - A basic user interface with a sidebar and content area is in place.
    - The default landing page after login is now `input-permintaan-badan-usaha.html`.
    - A dummy login form is functional.
    - The "Input Permintaan IDeb" form can be submitted to the backend, with "internal" search type triggering data ingestion.
    - The "Daftar Permintaan IDeb" pages (`debitur-individual.html` and `badan-usaha.html`) now fetch and display data from the backend, and the "Lihat Detail" link triggers PDF generation.
- **Integration:** The frontend is successfully communicating with the backend API for all implemented features.

## 2. What's Left to Build
- **PDF Generation:** The "Lihat Detail" link is a placeholder. The backend needs to be able to generate a PDF report for a completed request.
- **Asynchronous OJK Queries:** The mechanism for simulating a "live" asynchronous call to SLIK OJK has been implemented, including a simulated delay and updating request status to "Selesai" with dummy data storage in `get_idebs` table.

- **UI/UX Refinements:** The frontend is basic and can be improved with better styling, loading indicators, and user feedback.

## 3. Current Status
- **Status:** v0 Mockup - Initial Implementation Complete.
- **Details:** The core functionality of the mockup is in place. A user can log in, submit a request, and see the list of requests. The application is running locally with the backend server serving the frontend files.

## 4. Known Issues
- No real error handling is implemented on the frontend or backend beyond basic console logs and alerts.
- Comprehensive unit and integration tests have been added for the backend Go application.