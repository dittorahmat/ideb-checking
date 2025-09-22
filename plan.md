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
    -   A `fetch` request will be sent to the backend API endpoint for similar searches.
    -   The JSON response (containing potentially multiple results, like in `memory-bank/similar.json`) will be processed.
    -   The `data` array from the response will be stored in the browser's `sessionStorage`.
    -   After successfully storing the data, the page will redirect to `frontend/badan-usaha.html`.

### b. `frontend/badan-usaha.html`

-   **Dynamic "Lihat Detail" Link:**
    -   On page load, the JavaScript will check if there are similar search results in `sessionStorage`.
    -   If results are found, the "Lihat Detail" link will be dynamically altered. Instead of a direct navigation link, it will trigger a JavaScript function to open a popup/modal.
-   **Popup/Modal:**
    -   The popup will fetch and display the content from a new file: `frontend/output-similar-badan-usaha.html`.

### c. `frontend/output-similar-badan-usaha.html` (New File)

-   **Purpose:** This file will serve as the template for the popup that displays the similar search results.
-   **Content:**
    -   A table to display the list of results.
    -   Table headers will include fields from the JSON data (e.g., "Kode Jenis Pelapor", "Kode Pelapor", "Nama Debitur", "Nomor Identitas") and a new "Aksi" column.
    -   The table body will be populated dynamically using JavaScript.

### d. `frontend/script.js` (or new embedded script)

-   **New Functions:**
    -   **`handleSimilarSearch(event)`**:
        -   Attached to the similar search form.
        -   Handles the `fetch` request and `sessionStorage` logic.
    -   **`displaySimilarResultsPopup()`**:
        -   Opens the popup/modal.
        -   Loads `output-similar-badan-usaha.html`.
        -   Calls a function to populate the results table.
    -   **`populateSimilarResultsTable()`**:
        -   Retrieves the search results from `sessionStorage`.
        -   Dynamically creates table rows in `output-similar-badan-usaha.html` for each result.
        -   Each row will include a "Lihat Status" link in the "Aksi" column. This link will have `data-*` attributes to store the unique identifiers for that row's data (e.g., `data-kode-pelapor`).
    -   **`handlePrintPdf(event)`**:
        -   Attached to the "Lihat Status" links in the popup.
        -   Reads the data from the `data-*` attributes of the clicked link.
        -   Makes a `fetch` request to the original backend endpoint responsible for generating the PDF, passing the specific data for the selected entry.
        -   Handles the response to trigger the PDF download/display.

## 3. User Flow Summary

1.  User navigates to `input-permintaan-badan-usaha.html` and fills out the "Similar Match Search" form.
2.  User clicks "Submit".
3.  The page shows a loading animation, sends the request, stores the results in `sessionStorage`, and redirects to `badan-usaha.html`.
4.  On `badan-usaha.html`, the user clicks the "Lihat Detail" link.
5.  A popup appears, showing a table of the similar match results from `output-similar-badan-usaha.html`.
6.  The user finds the desired result in the table and clicks the "Lihat Status" link for that row.
7.  The system initiates the original PDF generation process for that specific entity and the PDF is displayed or downloaded.

---
Once you approve this plan, I will proceed with creating and modifying the necessary files.
