# Plan: Implement Similar Match Search Flow for Badan Usaha

This document outlines the plan to modify the "Similar Match Search" functionality for "Badan Usaha" as per the user's request.

## 1. Project Overview

The goal is to change the behavior of the "Similar Match Search" in `input-permintaan-badan-usaha.html`. Instead of the current flow, a similar search will now potentially return multiple results. These results will be displayed in a popup table on the `badan-usaha.html` page. From this table, the user can then trigger the original PDF generation for a specific result.

The "Exact Match Search" flow will remain unchanged.

## 2. File Changes and Implementation Details

### a. `frontend/input-permintaan-badan-usaha.html`

-   **Modify Form Submission:**
    -   The `<form id="ideb-request-form-badan-usaha-similar">` will be updated to prevent the default browser submission.
    -   An event listener will be added to the "Submit" button.
-   **JavaScript Logic:**
    -   When the form is submitted, a loading indicator will be displayed.
    -   A `fetch` request will be sent to the new backend API endpoint for similar searches (`/api/similar-search`).
    -   The JSON response (containing potentially multiple results, like in `memory-bank/similar.json`) will be processed.
    -   If the response contains data, the `data` array from the response will be stored in the browser's `sessionStorage`.
    -   After successfully storing the data, the page will redirect to `frontend/badan-usaha.html`.
    -   If no results are found or an error occurs, a popup message "Data tidak ditemukan" will be displayed, with detailed error logging in the backend.

### b. `frontend/badan-usaha.html`

-   **Dynamic "Lihat Detail" Link:**
    -   On page load, the JavaScript will check if there are similar search results in `sessionStorage`.
    -   If results are found, the "Lihat Detail" link will be dynamically altered. Instead of a direct navigation link, it will trigger a JavaScript function to open a popup/modal.
-   **Popup/Modal:**
    -   The popup will use Bootstrap's modal component to display the similar search results.
    -   The modal will contain a table to display the list of results.

### c. `frontend/output-similar-badan-usaha.html` (New File)

-   **Purpose:** This file will serve as the template for the popup that displays the similar search results.
-   **Content:**
    -   A Bootstrap modal structure.
    -   A table to display the list of results.
    -   Table headers will include fields from the JSON data (e.g., "Kode Jenis Pelapor", "Kode Pelapor", "Nama Debitur", "Nomor Identitas") and a new "Aksi" column.
    -   The table body will be populated dynamically using JavaScript.

### d. Backend Changes (`backend/handlers.go`, `backend/routes.go`)

-   **New API Endpoint:**
    -   Create a new endpoint `/api/similar-search` to handle similar search requests.
    -   This endpoint will read and return the data from `memory-bank/similar.json`.
    -   Proper error handling with console logging when no results are found or other errors occur.

### e. `frontend/script.js` (or new embedded script)

-   **New Functions:**
    -   **`handleSimilarSearch(event)`**:
        -   Attached to the similar search form.
        -   Handles the `fetch` request to `/api/similar-search` and `sessionStorage` logic.
        -   Displays the same loading indicator as used for exact match search.
    -   **`displaySimilarResultsPopup()`**:
        -   Opens a Bootstrap modal.
        -   Dynamically creates the modal content with a table for results.
        -   Calls a function to populate the results table.
    -   **`populateSimilarResultsTable()`**:
        -   Retrieves the search results from `sessionStorage`.
        -   Dynamically creates table rows in the modal for each result.
        -   Each row will include a "Lihat Status" link in the "Aksi" column. This link will have `data-*` attributes to store the unique identifiers for that row's data.
    -   **`handlePrintPdf(event)`**:
        -   Attached to the "Lihat Status" links in the popup.
        -   Reads the data from the `data-*` attributes of the clicked link.
        -   Makes a `fetch` request to the original backend endpoint responsible for generating the PDF, passing the specific data for the selected entry.
        -   Handles the response to automatically trigger the PDF download.

## 3. User Flow Summary

1.  User navigates to `input-permintaan-badan-usaha.html` and fills out the "Similar Match Search" form.
2.  User clicks "Submit".
3.  The page shows the same loading animation as exact match search, sends the request to `/api/similar-search`, stores the results in `sessionStorage`, and redirects to `badan-usaha.html`.
4.  On `badan-usaha.html`, the user clicks the "Lihat Detail" link.
5.  A Bootstrap modal popup appears, showing a table of the similar match results.
6.  The user finds the desired result in the table and clicks the "Lihat Status" link for that row.
7.  The system initiates the original PDF generation process for that specific entity and the PDF is automatically downloaded.

## 4. Technical Implementation Details

-   **SessionStorage Management:** Data will persist in sessionStorage until the browser tab is closed.
-   **Error Handling:** 
    -   Frontend will display "Data tidak ditemukan" popup for no results or errors.
    -   Backend will log detailed error messages to console for debugging.
-   **UI/UX:** 
    -   Consistent loading indicator with exact match search.
    -   Bootstrap modal for popup implementation.
    -   "Lihat Status" buttons will automatically trigger PDF download.

## 5. Implementation Status: COMPLETED

All planned functionality has been successfully implemented as described above. The similar match search flow is fully operational with the following key features:

- Similar search form in `input-permintaan-badan-usaha.html` works correctly
- Results are fetched from the backend and displayed in a modal on the `badan-usaha.html` page
- PDF generation works for individual results from the similar search
- All error handling and user experience features are in place
