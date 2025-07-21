# Active Context

## 1. Current Work Focus
The current focus is on initializing the project and establishing the foundational structure for the v0 mockup of the Ideb Checking application. This involves setting up the memory bank, defining the initial architecture, and creating a clear plan for development.

## 2. Next Steps
The immediate next steps are to build the core components of the v0 mockup:
1.  **Project Scaffolding:** Create the directory structure for the monorepo (e.g., `/frontend`, `/backend`).
2.  **Database Setup:** Define the SQL schema for the tables based on the provided JSON sample and create the initial SQLite database file.
3.  **Backend API (Go):**
    - Create a basic web server using the `net/http` package.
    - Implement the dummy login endpoint.
    - Implement the endpoint to receive new IDEB requests and save them to the database.
    - Implement the endpoint to list all existing requests.
- **Frontend (HTML/Bootstrap):**
    - The default landing page after login has been changed to `input-permintaan-badan-usaha.html`.
4.  **Frontend (HTML/Bootstrap):**
    - Create the main `index.html` with the sidebar navigation.
    - Build the dummy login page.
    - Build the "Input Permintaan IDeb" form.
    - Build the "Daftar Permintaan IDeb" page to display request statuses.

## 3. Key Decisions & Considerations
- **Dummy Data:** For v0, we will use hardcoded or easily generated dummy data for responses. The provided JSON sample will be used to model the database structure.
- **PDF Generation:** A simple Go library for PDF generation will be chosen. The focus will be on functionality over complex formatting for v0.
- **Error Handling:** Basic error handling will be implemented for API calls, but comprehensive error management is deferred to v1.
