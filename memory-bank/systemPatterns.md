# System Patterns

## 1. System Architecture
The application follows a simple two-tier client-server architecture:
- **Frontend:** A web-based user interface built with HTML, JS, and Bootstrap. It is responsible for all user interaction and rendering of data. It communicates with the backend via API calls.
- **Backend:** A Go application that serves the API. It contains all the business logic, interacts with the SQLite database using GORM, and communicates with the external SLIK OJK API.

For v1, this architecture will be maintained but deployed as a single unit within a monorepo to Azure App Service.

## 2. Data Flow

### New Request Data Flow
1.  **User** submits the "Input Permintaan IDeb" form on the frontend.
2.  **Frontend** sends a POST request containing the form data to the Go backend API.
3.  **Backend** receives the request and inserts a new record into the `requests` table in the SQLite database with a status of "Dalam Proses".
4.  If the request is a "live SLIK OJK" query, the backend initiates an asynchronous call to the external SLIK API in a separate goroutine.
5.  **Backend** immediately returns a success response to the frontend, confirming the request has been submitted.

### Status & Detail View Data Flow
1.  **User** navigates to the "Daftar Permintaan IDeb" page.
2.  **Frontend** sends a GET request to the backend API to fetch all requests for the user.
3.  **Backend** queries the SQLite database and returns the list of requests with their current statuses.
4.  **Frontend** renders the list.
5.  When the async OJK job is complete, it updates the corresponding request's status in the database to "Ready" and stores the received data.
6.  If a user clicks "Lihat Detail" on a "Ready" request, the **Frontend** requests the data from the **Backend**, which then generates a PDF to be displayed or downloaded.

## 3. Key Components
- **Authentication:** A dummy login system for v0 that accepts any credentials.
- **Request Input Module:** The HTML form and corresponding backend endpoint for submitting new IDEB requests.
- **Request List Module:** The UI and backend endpoint for displaying the status of all submitted requests.
- **PDF Generation Service:** A backend service that takes the detailed data of a completed request and formats it into a printable PDF.
- **SLIK OJK API Client:** A module within the backend responsible for making and managing asynchronous requests to the external OJK API.
