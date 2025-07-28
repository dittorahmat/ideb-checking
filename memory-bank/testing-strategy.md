The current testing strategy is extremely limited. Based on the `backend/models_test.go` file, the only tests present are basic unit tests verifying the `TableName()` method for the `Request`, `CorporateRequest`, and `GetIdeb` structs. These tests use the `github.com/stretchr/testify/assert` library.

**Current Testing Strategy:**
*   **Scope:** Only covers the `TableName()` methods of a few data models.
*   **Type:** Basic unit tests.
*   **Tools:** Go's built-in `testing` package and `testify/assert`.

**How to Improve the Testing Strategy:**

The current strategy provides almost no confidence in the application's correctness. Here's a comprehensive plan for improvement, starting with the backend:

## Backend Testing Strategy (Go)

### 1. Backend Unit Tests

**Purpose:** To test individual functions or components in isolation, ensuring they work as expected without external dependencies.

*   **Handlers (`handlers.go`):**
    *   Write unit tests for each handler function (e.g., `loginHandler`, `createRequestHandler`, `generatePDFHandler`).
    *   **Mocking:** These tests should mock external dependencies such as the database (`DB`), file system operations (`os.ReadFile`), and any external API calls (e.g., simulated SLIK OJK). This ensures that only the handler's logic is being tested.
    *   **Tools:** `net/http/httptest` for simulating HTTP requests and responses, `testify/assert` for assertions, and `sqlmock` or `go-sqlmock` for mocking database interactions.

*   **Helper Functions (`handlers.go`, `utils.go`, `pdf.go`):**
    *   Test any significant helper functions (e.g., `handleInternalSearch`, `handleLiveSearch`, `readAndUnmarshalInputJSON`, `generateIdebPDF`, `respondWithError`, `respondWithJSON`) in isolation.

*   **Models (`models.go`):**
    *   While `TableName()` methods are covered, consider adding tests for any custom methods, validation logic, or data transformations within the model structs if they are introduced.

### 2. Backend Integration Tests

**Purpose:** To test the interaction between multiple components or modules, ensuring they work together correctly.

*   **API Endpoints:**
    *   Create tests that make actual HTTP requests to the API endpoints (e.g., `/api/login`, `/api/requests`, `/api/generate-pdf`).
    *   **Assertions:** Assert on the HTTP responses (status codes, headers, and body content).
    *   **Test Database:** Use a dedicated test database (e.g., an in-memory SQLite database initialized with `setupTestDB()`) to ensure that database interactions are working correctly and tests are isolated from development data.
    *   **Tools:** `net/http/httptest` for making real HTTP requests, `testify/assert` for assertions, and GORM's in-memory SQLite for database interactions.

*   **Database Interactions:**
    *   Verify that data is correctly saved, retrieved, updated, and deleted through the API endpoints.

*   **PDF Generation:**
    *   Test the `generatePDFHandler` by making a request and verifying that a PDF is returned. While detailed content verification of the PDF might be complex without a dedicated PDF testing library, you can at least verify the content type and basic structure.

## Frontend Testing Strategy (JavaScript)

*(To be expanded upon later, focusing on backend first as requested.)*

*   **Unit Tests (JavaScript):** Introduce a JavaScript testing framework (e.g., Jest, Mocha with Chai) to test individual JavaScript functions in `script.js` (e.g., `login`, `loadContent`, `submitIdebRequest`, `loadRequests`). Mock DOM interactions and API calls.

*   **End-to-End (E2E) Tests:** Use a tool like Cypress or Playwright to simulate user interactions in a real browser. These tests would cover the entire user flow, from logging in to submitting requests and viewing PDFs, ensuring that the frontend and backend integrate seamlessly.

## Recommended Tools (Summary)

*   **Go Backend:**
    *   `net/http/httptest`: For testing HTTP handlers and making HTTP requests.
    *   `testify/assert`: For rich assertions.
    *   `sqlmock` or `go-sqlmock`: For mocking database interactions in unit tests.
    *   GORM's in-memory SQLite: For integration tests.

*   **JavaScript Frontend:**
    *   **Unit Testing:** Jest (popular, good for mocking).
    *   **E2E Testing:** Cypress or Playwright (both excellent for browser automation).

## Implementation Steps

1.  **Start with Backend Integration Tests:** These provide the most immediate value by verifying the core API functionality and interactions between components.
2.  **Add Backend Unit Tests:** Once integration tests are stable, add unit tests for more granular control and faster feedback during development.
3.  **Introduce Frontend Unit Tests:** Begin testing individual JavaScript components.
4.  **Implement Frontend E2E Tests:** Build confidence in the complete user experience.

This phased approach will gradually build a robust testing suite, significantly improving the reliability and maintainability of the application.