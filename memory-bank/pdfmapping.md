# PDF Mapping Plan

This document outlines the phased approach to enhance the PDF generation in `backend/handlers.go` to closely match the `SARIPARI-PERTIWI-ABADI-b.pdf` layout and utilize all relevant data from `input.json`.

## Phase 1: Core Information and Basic Structure
*   **Objective:** Ensure all top-level header information and the primary corporate debtor details are correctly mapped and displayed in the PDF.
*   **Details:**
    *   Map all fields from `input.json -> data -> header` (e.g., `userReferenceCode`, `resultDate`, `inquiryId`, `inquiryDate`).
    *   Map all fields from `input.json -> data -> corporate` (e.g., `reportNumber`, `latestDataYearMonth`, `requestDate`, `corporateKeyWord`).
    *   For the first `corporateDebtor` in `input.json -> data -> corporate -> corporateDebtors`, map all its direct fields (e.g., `fullName`, `taxId`, `companyTypeDesc`, `estPlace`, `estCertNo`, `estCertDate`, `memberDesc`, `updatedDatetime`, `address`, `subDistrict`, `district`, `cityDesc`, `postalCode`, `countryDesc`, `economicSectorDesc`, `ratingDate`, `createdDatetime`, `goPublicFlag`).
    *   Focus on getting the data correctly extracted and placed, even if the layout isn't perfectly aligned with `SARIPARI-PERTIWI-ABADI-b.pdf` yet.

## Phase 2: Iterating Corporate Debtors and Basic Shareholder Information
*   **Objective:** Extend the PDF to include all corporate debtors and introduce the "Pemilik / Pengurus" section with their basic identifying information.
*   **Details:**
    *   Iterate through all `corporateDebtors` in `input.json` and display their mapped fields from Phase 1.
    *   For each `corporateDebtor`, if `officisSharehldrsGroups` exists, iterate through these groups.
    *   For each `officisSharehldrsGroup`, display its `memberDesc`.
    *   For each `officisSharehldrs` within a group, display their `identityNumberName` and `identityNumber`.

## Phase 3: Detailed Shareholder Information and Layout Refinements
*   **Objective:** Map all remaining detailed shareholder information and refine the overall PDF layout to closely match `SARIPARI-PERTIWI-ABADI-b.pdf`.
*   **Details:**
    *   For each `officisSharehldrs`, map all its remaining fields (e.g., `genderDesc`, `jobPositionDesc`, `shareOwnership`, `address`, `district`, `cityDesc`, `shareholderStatusDesc`, `subDistrict`).
    *   Implement more advanced `Maroto` features for precise positioning, tables, and multi-column layouts to match the visual structure of `SARIPARI-PERTIWI-ABADI-b.pdf`. This will involve careful adjustments to row heights, column widths, and text properties.
    *   Address pagination to ensure content flows correctly across multiple pages, similar to the 7-page sample.
