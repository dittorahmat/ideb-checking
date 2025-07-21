# Ideb Checking Application (v0 Mockup)

This repository contains a version 0 (v0) mockup of an "Ideb Checking" application, designed primarily for presales demonstrations. It showcases the core functionality of a system that acts as a middleware for the OJK's SLIK (Sistem Layanan Informasi Keuangan) API.

## Project Overview

The application addresses the challenge of slow data retrieval from the SLIK OJK API (which can take minutes to an hour) by providing a streamlined interface for users to query credit information. It stores previously fetched data locally for instant access and manages asynchronous live data requests.

## Features

### Implemented (v0 Mockup)

*   **User Authentication:** Dummy login functionality (any email/password works).
*   **Request Submission:** Users can submit new IDEB requests with details like reference number, purpose, identity type, and identity number.
*   **Search Options:** Choice between searching internal database or initiating a "live" query to SLIK OJK.
*   **Request Listing:** Displays a list of all submitted requests with their current status.
    *   "Dalam Proses" for pending requests.
    *   "Lihat Detail" link for completed requests (placeholder for PDF generation).
*   **Frontend Navigation:** Basic sidebar navigation with collapsible menus for different sections.
*   **Backend API:** RESTful API for login, creating requests, and listing requests.
*   **Database:** SQLite for local data storage.

### To Be Implemented (Future Versions)

*   **PDF Generation:** Generate detailed PDF reports for completed requests.
*   **Asynchronous OJK Queries:** Full implementation of the asynchronous call mechanism to SLIK OJK.
*   **Dashboard:** Functional dashboard page.
*   **Parameter Management:** Pages for User & Password API, Valid Token, and LDAP settings.
*   **User Management:** Functionality for managing application users.
*   **UI/UX Refinements:** Enhanced styling, loading indicators, and user feedback.

## Technology Stack

### Version 0 (Mockup)

*   **Frontend:** HTML, JavaScript, Bootstrap, CSS
*   **Backend:** Go (using `net/http` for simplicity)
*   **Database:** SQLite (`github.com/mattn/go-sqlite3`)

### Future Considerations (Version 1)

*   **Database:** PostgreSQL
*   **Backend Framework (Optional):** Gin or Echo for more complex routing/middleware.
*   **Caching:** Redis (for performance optimization).
*   **Asynchronous Processing:** Go's native goroutines and channels (no external message brokers needed for now).
*   **Deployment:** Azure App Service (monorepo setup).

## Project Structure

```
.
├── GEMINI.md
├── README.md
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   └── schema.sql
├── frontend/
│   ├── badan-usaha.html
│   ├── debitur-individual.html
│   ├── index.html
│   ├── input-permintaan-badan-usaha.html
│   ├── input-permintaan-individual.html
│   ├── script.js
│   └── style.css
└── memory-bank/
    ├── activeContext.md
    ├── productContext.md
    ├── progress.md
    ├── projectBrief.md
    ├── systemPatterns.md
    └── techContext.md
```

## Setup and Running the Application

1.  **Clone the repository:**
    ```bash
    git clone <repository_url>
    cd ideb/app
    ```

2.  **Backend Setup:**
    Navigate to the `backend` directory and download Go modules:
    ```bash
    cd backend
    go mod tidy
    ```
    Run the backend server:
    ```bash
    go run main.go
    ```
    The backend server will serve the frontend files and handle API requests.

3.  **Access the Frontend:**
    Open your web browser and navigate to `http://localhost:8080` (or the port specified by the Go application).

4.  **Login:**
    Use any email and password to log in (e.g., `test@example.com`, `password`). After successful login, you will be redirected to the "Input Permintaan IDeb - Badan Usaha" page.
