# Product Context

## 1. Why This Project Exists
The primary goal is to create a version 0 (v0) mockup of an "Ideb Checking" application. This mockup will serve as a functional demo for presales activities, showcasing the core capabilities of the planned final product.

## 2. Problem It Solves
The application acts as a middleware to interact with the OJK's SLIK (Sistem Layanan Informasi Keuangan) API, which provides credit information for individuals and businesses. Accessing data from SLIK can be a slow process, sometimes taking up to an hour. This application aims to streamline that process by:
- Allowing users to query for credit information through a simplified interface.
- Storing previously fetched data locally, enabling instant access to historical reports.
- Managing the asynchronous nature of live data requests from OJK.

## 3. How It Should Work (User Experience)
A user logs into the application (using dummy credentials for v0). They can then perform two main actions:
1.  **Submit a new request:** The user fills out a form with required details (`nomor_referensi_pengguna`, `tujuan_penggunaan`, etc.) and chooses whether to search the internal database or initiate a live query to SLIK OJK.
2.  **View existing requests:** The user can navigate to a list of all submitted requests. This list displays the status of each request.
    - If a request is still processing, it shows "Dalam Proses".
    - If a request is complete, it provides a "Lihat Detail" link.
    - Clicking the link generates and displays a PDF report of the credit information.

The interface should be straightforward, with a collapsible sidebar for navigation between different sections like Input, Request List, Dashboard, and User Management.
