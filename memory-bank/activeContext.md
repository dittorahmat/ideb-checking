# Active Context

## 1. Current Work Focus
The current focus is on enhancing the PDF generation for the v0 mockup of the Ideb Checking application, currently in **Phase 2** of the plan outlined in `pdfmapping.md`. The backend has been refactored, is running in the background, and comprehensive backend tests are in place. Data ingestion from `input.json`, PDF generation, asynchronous OJK query simulation, and placeholder pages are implemented.

## 2. Next Steps
The immediate next steps are to build the core components of the v0 mockup:
1.  **Project Scaffolding:** Create the directory structure for the monorepo (e.g., `/frontend`, `/backend`).
2.  **Database Setup:** Define the SQL schema for the tables based on the provided JSON sample and create the initial SQLite database file, including the `get_idebs` table.
3.  **Backend API (Go):**
    - Create a basic web server using the `net/http` package.
    - Implement the dummy login endpoint.
    - Implement the endpoint to receive new IDEB requests and save them to the database using GORM, including ingestion of `input.json` data for "internal" search types and asynchronous simulation for "live" search types (currently has a bug in data population for "live" search types).
    - Implement the endpoint to list all existing requests using GORM.
    - Implement PDF generation from `get_idebs` table data, with Phase 1 (Core Information and Basic Structure) complete. Proceeding to Phase 2 (Iterating Corporate Debtors and Basic Shareholder Information).
    - **Refactoring:** `main.go` has been refactored into `database.go` (for DB initialization), `models.go` (for data structures), `handlers.go` (for HTTP handlers), and `routes.go` (for route registration) to improve modularity and maintainability.
4.  **Frontend (HTML/Bootstrap):**
    - Create the main `index.html` with the sidebar navigation.
    - Build the dummy login page.
    - Build the "Input Permintaan IDeb" form, with "internal" search type correctly submitting data.
    - Build the "Daftar Permintaan IDeb" page to display request statuses, with "Lihat Detail" linking to PDF generation.
    - Implement placeholder pages for Dashboard, Parameter, and User Management.

## 3. Key Decisions & Considerations
- **Dummy Data:** For v0, we will use hardcoded or easily generated dummy data for responses. The provided JSON sample will be used to model the database structure.
- **PDF Generation:** A simple Go library for PDF generation will be chosen. The focus will be on functionality over complex formatting for v0.
- **Error Handling:** Basic error handling will be implemented for API calls, but comprehensive error management is deferred to v1.