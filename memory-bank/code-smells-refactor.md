# Code Smells and Refactoring Plan

## Identified Code Smells:

1.  **Duplication in `handlers.go` (`handleInternalSearch` and `handleLiveSearch`):**
    Both `handleInternalSearch` and `handleLiveSearch` contained identical logic for reading `input.json` and unmarshalling it into an `InputJSON` struct. This has been resolved by extracting the logic into `readAndUnmarshalInputJSON`.

2.  **Long Function in `handlers.go` (`generatePDFHandler`):**
    The `generatePDFHandler` function was lengthy and performed multiple responsibilities. This has been resolved by extracting the PDF generation logic into a new `pdf.go` file and the `generateIdebPDF` function.

## Remaining Code Smells & Areas for Future Refactoring:

1.  **Inconsistent Naming/Usage of `Request` and `CorporateRequest` Models:**
    *   **Description:** The `Request` struct is used for both individual and corporate requests in `createRequestHandler`, but `Request` and `CorporateRequest` have distinct `TableName()` methods. This can lead to confusion or issues if the application grows and these models diverge more significantly.
    *   **Suggestion:** For v1, consider having separate request DTOs (Data Transfer Objects) for individual and corporate requests, and then mapping them to a unified internal `Request` model if the underlying database table is indeed unified. Alternatively, if they truly map to different tables, ensure the API endpoints and handling clearly distinguish between them from the start.

2.  **Global Variable `readFileFunc` in `handlers.go`:**
    *   **Description:** While it aids testing, `readFileFunc` is a global variable. Global state can make code harder to reason about and can introduce unexpected side effects in larger applications.
    *   **Suggestion:** For v1, consider dependency injection for `readFileFunc` (e.g., passing it as a parameter to `App` or to the handlers that need it).

3.  **Error Handling (General):**
    *   **Description:** The current error handling is basic (`respondWithError`, `log.Println`). For a production application (v1), a more robust and centralized error handling mechanism would be beneficial (e.g., custom error types, middleware for error logging and consistent error responses).
    *   **Suggestion:** Implement custom error types and a centralized error handling middleware for v1.

4.  **Configuration Loading in `main.go`:**
    *   **Description:** The `NewConfig` function directly uses `os.Getenv` and hardcoded default values. While functional, for a more complex application, a dedicated configuration library (e.g., `viper`) could provide more flexibility (e.g., reading from config files, environment variables, command-line flags).
    *   **Suggestion:** For v1, consider using a configuration library for more flexible and robust configuration management.

5.  **Static File Serving in `routes.go`:**
    *   **Description:** Serving static files directly from the Go backend is fine for a small mockup. In a production environment, it's often more efficient to serve static files via a dedicated web server (e.g., Nginx, Apache) or a CDN.
    *   **Suggestion:** For v1, consider serving static files through a dedicated web server or CDN.