# System Patterns

## 1. System Architecture
The application follows a simple two-tier client-server architecture:
- **Frontend:** A web-based user interface built with HTML, JS, and Bootstrap. It is responsible for all user interaction and rendering of data. It communicates with the backend via API calls.
- **Backend:** A Go application that serves the API. It contains all the business logic, interacts with the SQLite database using GORM, and communicates with the external SLIK OJK API.

For v1, this architecture will be maintained but deployed as a single unit within a monorepo to Azure App Service.

## 2. Data Flow

### New Request Data Flow
1.  **User** submits the "Input Permintaan IDeb" form on the frontend.
2.  **Frontend** sends a POST request containing the form data to the Go backend API (`/api/requests`).
3.  **Backend** receives the request and inserts a new record into the `requests` table in the SQLite database with a status of "Dalam Proses".
4.  If the request has `search_type` "internal", the backend reads data from `input.json`, extracts `userReferenceCode` and the first `taxId`, and stores this information along with the full JSON content into the `get_idebs` table. The status of the original request is then updated to "Selesai".
5.  If the request has `search_type` "live", the backend simulates an asynchronous call to the external SLIK OJK API (with a simulated delay). After the simulated call, the corresponding request's status in the `requests` table is updated to "Selesai", and dummy data is stored in the `get_idebs` table.
6.  **Backend** immediately returns a success response to the frontend, confirming the request has been submitted.

### Status & Detail View Data Flow
1.  **User** navigates to the "Daftar Permintaan IDeb" page.
2.  **Frontend** sends a GET request to the backend API to fetch all requests for the user.
3.  **Backend** queries the SQLite database and returns the list of requests with their current statuses.
4.  **Frontend** renders the list.
5.  If a user clicks "Lihat Detail" on a "Selesai" request, the **Frontend** makes a request to the `/api/generate-pdf` endpoint with the request ID.
6.  The **Backend** retrieves the corresponding data from the `get_idebs` table and generates a PDF report, which is then sent back to the frontend for display or download.

## 3. Key Components
- **Authentication:** A dummy login system for v0 that accepts any credentials.
- **Request Input Module:** The HTML forms (`input-permintaan-individual.html`, `input-permintaan-badan-usaha.html`) and corresponding backend endpoint (`/api/requests`) for submitting new IDEB requests.
- **Request List Module:** The UI (`debitur-individual.html`, `badan-usaha.html`) and backend endpoints (`/api/getDebtorExactIndividual`, `/api/getDebtorExactCorporate`) for displaying the status of all submitted requests.
- **PDF Generation Service:** A backend service (`generatePDFHandler`) that takes detailed data from the `get_idebs` table and formats it into a printable PDF using the `Maroto` library.
- **SLIK OJK API Client (Simulated):** A module within the backend responsible for simulating asynchronous requests to the external OJK API and storing dummy data in the `get_idebs` table.
- **Placeholder Pages:** Frontend pages for Dashboard (`dashboard.html`), Parameter (`parameter-user-api.html`, `parameter-valid-token.html`, `parameter-ldap.html`), and User Management (`user-management.html`).

## 4. Backend Go Application Structure (Refactored)
To improve modularity and maintainability, the Go backend application has been refactored into the following files:

- **`main.go`**
    - **Description:** The entry point of the application. It initializes the database and registers all HTTP routes.
    - **Key Functions:** `main()`

- **`database.go`**
    - **Description:** Handles the database connection and schema migration.
    - **Key Functions:** `InitDatabase()`

- **`models.go`**
    - **Description:** Defines the data structures (structs) for database models, such as `Request`, `CorporateRequest`, and `GetIdeb` (for storing SLIK OJK data).
    - **Key Functions:** `Request.TableName()`, `CorporateRequest.TableName()`, `GetIdeb.TableName()`

- **`handlers.go`**
    - **Description:** Contains the HTTP handler functions for various API endpoints.
    - **Key Functions:** `loginHandler()`, `getDebtorExactIndividualHandler()`, `getDebtorExactCorporateHandler()`, `createRequestHandler()`, `createRequest()`, `getRequests()`, `generatePDFHandler()`

- **`routes.go`**
    - **Description:** Registers all the HTTP routes and associates them with their respective handler functions.
    - **Key Functions:** `RegisterRoutes()`
    - **Registered Routes:** `/api/login`, `/api/requests`, `/api/getDebtorExactIndividual`, `/api/getDebtorExactCorporate`, `/api/generate-pdf`