# Code Smells and Refactoring Plan

## Identified Code Smells:

1.  **Duplication in `handlers.go` (`handleInternalSearch` and `handleLiveSearch`):**
    Both `handleInternalSearch` and `handleLiveSearch` contained identical logic for reading `input.json` and unmarshalling it into an `InputJSON` struct. This has been resolved by extracting the logic into `readAndUnmarshalInputJSON`.

2.  **Long Function in `handlers.go` (`generatePDFHandler`):**
    The `generatePDFHandler` function was lengthy and performed multiple responsibilities. This has been resolved by extracting the PDF generation logic into a new `pdf.go` file and the `generateIdebPDF` function.

## Resolved Code Smells:

1.  **Inconsistent Naming/Usage of `Request` and `CorporateRequest` Models:**
    *   **Description:** Previously, the `Request` struct was used for both individual and corporate requests in `createRequestHandler`, despite `Request` and `CorporateRequest` having distinct `TableName()` methods. This led to potential confusion and data inconsistencies.
    *   **Resolution:** Introduced separate DTOs (`IndividualRequestPayload` and `CorporateRequestPayload`) for individual and corporate requests. `createRequestHandler` now unmarshals incoming requests into the appropriate DTO based on a `request_type` field in the payload. The `handleInternalSearch` and `handleLiveSearch` functions were updated to accept `interface{}` and handle both `*Request` and `*CorporateRequest` types, ensuring correct data handling and persistence.

## Resolved Code Smells:

2.  **Global Variable `readFileFunc` in `handlers.go`:**
    *   **Description:** Previously, `readFileFunc` was a global variable, which could lead to unexpected side effects and make testing harder.
    *   **Resolution:** `readFileFunc` has been refactored from a global variable to a field within the `App` struct. This improves testability and reduces reliance on global state by enabling dependency injection.

## Resolved Code Smells:

3.  **Error Handling (General):**
    *   **Description:** Previously, error handling was basic, relying on generic `respondWithError` calls.
    *   **Resolution:** Custom error types (`Error` struct and predefined error variables like `ErrNotFound`, `ErrInvalidInput`, etc.) have been introduced in `utils.go`. Handlers now return these specific error types, allowing for more granular error management and consistent API responses.

## Resolved Code Smells:

4.  **Configuration Loading in `main.go`:**
    *   **Description:** Previously, the `NewConfig` function directly used `os.Getenv` and hardcoded default values, limiting flexibility.
    *   **Resolution:** Integrated the `viper` library for configuration management. `NewConfig` now uses `viper` to load configuration values from `.env` files (if present) and environment variables, with sensible defaults. This provides a more flexible and robust configuration system.

5.  **Static File Serving in `routes.go`:**
    *   **Description:** Serving static files directly from the Go backend is fine for a small mockup. In a production environment, it's often more efficient to serve static files via a dedicated web server (e.g., Nginx, Apache) or a CDN.
    *   **Suggestion:** For v1, consider serving static files through a dedicated web server or CDN.