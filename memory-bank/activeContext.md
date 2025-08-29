# Active Context

## 1. Current Work Focus
The current focus is debugging a critical backend error that occurs during form submission. The frontend has been significantly improved and is now correctly handling user input and events, but the backend is returning a `500 Internal Server Error`.

## 2. Next Steps
The immediate next step is to reliably capture the panic message from the Go backend. The current logging to a file has been unsuccessful, likely due to log buffering. The plan is to run the server in the foreground to directly capture the `stdout`/`stderr` output when the crash is triggered.

## 3. Key Decisions & Considerations
- **Frontend Event Handling:** The method for handling events on dynamically loaded content has been completely refactored. Logic has been moved from non-functional inline scripts into the main `script.js`, with a centralized `loadContent` function that now correctly attaches event listeners after injecting HTML. This is a more robust and secure approach.
- **Debugging Strategy:** Standard logging has failed. The strategy has shifted to direct output capture by running the server process in the foreground. This is a more aggressive but necessary step to diagnose the stubborn bug.

## 4. Problems Faced
- **Solved:** The form submission process was completely broken due to a series of frontend issues, including a misunderstanding of how browsers handle scripts in dynamically loaded `innerHTML`. This was a complex issue to diagnose and has been fully resolved.
- **Active:** The Go backend panics when the `/api/requests` endpoint is hit. The root cause is unknown because the error message has not been successfully captured yet.

## 5. Immediate Plan to Resolve Problems
- The plan is to run the Go application in the foreground and have the user trigger the 500 error. This should print the panic message directly to the console, allowing for analysis and a fix.
