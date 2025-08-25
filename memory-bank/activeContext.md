# Active Context

## 1. Current Work Focus
The current focus is on improving the backend's configuration management. The application has been refactored to use the `viper` library for loading configuration from `.env` files and environment variables, replacing the previous method that relied on `os.Getenv` and hardcoded defaults in `main.go`. This enhances flexibility and robustness. Backend tests are passing, and UI/UX refinements for the frontend are planned next.

## 2. Next Steps
With the configuration refactor complete and backend tests passing, the next phase will focus on UI/UX refinements for the frontend.

## 3. Key Decisions & Considerations
- **Configuration Management:** Adopted `viper` for configuration loading. This provides a standardized way to manage configuration from multiple sources (`.env`, environment variables, defaults) and handles precedence automatically. The `config.yaml` file was removed.
- **Dependency Injection:** `readFileFunc` remains a dependency injected into the `App` struct, improving testability and reducing reliance on global state.
- **Structured Error Handling:** `APIError` is in place for consistent error responses, laying the groundwork for future robust error management.
- **Prioritization:** Focused on `readFileFunc` and error handling refactoring for v0. Other identified areas (model inconsistency, static file serving) are deferred to v1 due to their larger architectural impact.

## 4. Problems Faced
- **Resolved:** Persistent syntax errors in `handlers_test.go` have been resolved through manual cleanup and careful escaping of string literals. All backend tests are now passing.
- **Resolved:** Configuration loading was previously basic. It has now been improved using `viper`.

## 5. Immediate Plan to Resolve Problems
- **Resolved:** The plan to resolve persistent syntax errors in `handlers_test.go` has been successfully executed. All tests are now passing.
- **Resolved:** The plan to refactor configuration loading using `viper` has been successfully executed.