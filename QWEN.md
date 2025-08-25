# Ideb Checking Application (v0 Mockup) - Context for Qwen Code

## Project Overview

This project is a version 0 (v0) mockup of an "Ideb Checking" application. Its primary purpose is for presales demonstrations. The application acts as a middleware for the OJK's SLIK (Sistem Layanan Informasi Keuangan) API.

The core challenge it addresses is the slow response time of the SLIK OJK API (minutes to hours) by providing a user interface to query credit information, storing previously fetched data locally for instant retrieval, and managing asynchronous requests to the live API.

### Key Features (v0 Mockup)
*   **User Authentication:** Dummy login (any email/password works).
*   **Request Submission:** Submit new IDEB requests with details (reference number, purpose, identity type/number, etc.).
*   **Search Options:** Search internal database or initiate a "live" query to SLIK OJK.
*   **Request Listing:** View a list of submitted requests with status ("Dalam Proses" or "Lihat Detail").
*   **PDF Generation (Placeholder):** "Lihat Detail" link is a placeholder intended to generate a PDF report.
*   **Frontend Navigation:** Basic sidebar navigation.
*   **Backend API:** RESTful API for login, creating requests, and listing requests.
*   **Database:** Local SQLite storage.

### Technology Stack (v0)
*   **Frontend:** HTML, JavaScript, Bootstrap, CSS
*   **Backend:** Go (`net/http`, `gorilla/mux`)
*   **Database:** SQLite (`gorm.io/driver/sqlite`)

### Future Technology Considerations (v1)
*   **Database:** PostgreSQL
*   **Backend Framework (Optional):** Gin or Echo
*   **Caching:** Redis
*   **Asynchronous Processing:** Native Go goroutines/channels (no RabbitMQ/Kafka planned for now).
*   **Deployment:** Azure App Service (monorepo).

## Project Structure

```
.
├── README.md
├── Dockerfile
├── .dockerignore
├── docker-compose.yml
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── database.go
│   ├── handlers.go
│   ├── middleware.go
│   ├── models.go
│   ├── pdf.go
│   ├── routes.go
│   └── utils.go
├── frontend/
│   ├── badan-usaha.html
│   ├── dashboard.html
│   ├── debitur-individual.html
│   ├── index.html
│   ├── input-permintaan-badan-usaha.html
│   ├── input-permintaan-individual.html
│   ├── parameter-ldap.html
│   ├── parameter-user-api.html
│   ├── parameter-valid-token.html
│   ├── script.js
│   ├── style.css
│   └── user-management.html
└── memory-bank/ (Contains project context, notes, and sample data)
```

## Building and Running the Application

### Running with Go directly

1.  **Setup Backend:**
    *   Navigate to the `backend` directory: `cd backend`
    *   Download Go modules: `go mod tidy`
2.  **Run Backend Server:**
    *   From the `backend` directory, execute: `go run .`
    *   This command starts the Go server, which serves the frontend files and handles API requests. It typically runs on `localhost:8080`.
3.  **Access Frontend:**
    *   Open a web browser and go to `http://localhost:8080`.
4.  **Login:**
    *   Use any email and password (e.g., `test@example.com`/`password`) on the login page.

### Running with Docker

1.  **Build and run with Docker:**
    *   Build the Docker image: `docker build -t ideb-app .`
    *   Run the container: `docker run -d -p 8080:8080 --name ideb-container ideb-app`

2.  **Or use Docker Compose (recommended):**
    *   Build and run with docker-compose: `docker-compose up -d`

3.  **Access Frontend:**
    *   Open a web browser and go to `http://localhost:8080`.

4.  **Login:**
    *   Use any email and password (e.g., `test@example.com`/`password`) on the login page.

### Stopping the Docker Container

*   If using Docker directly: `docker stop ideb-container`
*   If using Docker Compose: `docker-compose down`

## Development Conventions

*   **Backend:** Go code is structured with `main.go` initializing the application, `handlers.go` containing HTTP request logic, `models.go` defining data structures and database tables, `database.go` for DB initialization, and `routes.go` for API endpoint registration.
*   **Frontend:** Uses basic HTML/CSS/JS with Bootstrap. `script.js` handles frontend logic like API calls, form submissions, and dynamic content loading.
*   **Data Handling:** Requests and fetched IDEB data are stored in SQLite tables using GORM. Sample data is loaded from `memory-bank/input.json` for the internal search mockup.
*   **API:** RESTful API endpoints under `/api/`.