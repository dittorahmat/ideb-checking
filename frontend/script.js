function login() {
    fetch('/api/login', { method: 'POST' })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'success') {
                document.getElementById('login-view').style.display = 'none';
                document.getElementById('app-view').style.display = 'block';
                loadContent('input-permintaan-badan-usaha.html');
            } else {
                alert('Login failed');
            }
        });
}

function loadContent(page) {
    fetch(page)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.text();
        })
        .then(html => {
            document.getElementById('main-content').innerHTML = html;

            // After loading content, find the forms and attach listeners
            if (page === 'input-permintaan-badan-usaha.html') {
                // Handle both similar and exact match forms for corporate requests
                const similarForm = document.getElementById('ideb-request-form-badan-usaha-similar');
                const exactForm = document.getElementById('ideb-request-form-badan-usaha-exact');
                
                if (similarForm) {
                    similarForm.addEventListener('submit', function(event) {
                        event.preventDefault();
                        submitIdebRequest('ideb-request-form-badan-usaha-similar');
                    });
                }
                
                if (exactForm) {
                    exactForm.addEventListener('submit', function(event) {
                        event.preventDefault();
                        submitIdebRequest('ideb-request-form-badan-usaha-exact');
                    });
                }
            } else if (page === 'input-permintaan-individual.html') {
                // Handle both similar and exact match forms for individual requests
                const similarForm = document.getElementById('ideb-request-form-individual-similar');
                const exactForm = document.getElementById('ideb-request-form-individual-exact');
                
                if (similarForm) {
                    similarForm.addEventListener('submit', function(event) {
                        event.preventDefault();
                        submitIdebRequest('ideb-request-form-individual-similar');
                    });
                }
                
                if (exactForm) {
                    exactForm.addEventListener('submit', function(event) {
                        event.preventDefault();
                        submitIdebRequest('ideb-request-form-individual-exact');
                    });
                }
            } else if (page === 'debitur-individual.html' || page === 'badan-usaha.html') {
                loadRequests(page);
            }
        })
        .catch(error => {
            console.error('Error loading page: ', error);
            document.getElementById('main-content').innerHTML = '<p>Error loading content.</p>';
        });
}

function setupFormListener() {
    // This function is no longer needed as forms will directly call submitIdebRequest
}

function showLoadingPopup() {
    const loadingPopup = document.createElement('div');
    loadingPopup.id = 'loading-popup';
    loadingPopup.className = 'loading-overlay';
    loadingPopup.innerHTML = `
        <div class="loading-content">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
            <p>Loading data, please wait...</p>
        </div>
    `;
    document.body.appendChild(loadingPopup);
    loadingPopup.style.display = 'flex';
}

function hideLoadingPopup() {
    const loadingPopup = document.getElementById('loading-popup');
    if (loadingPopup) {
        document.body.removeChild(loadingPopup);
    }
}

function submitIdebRequest(formId) {
    const form = document.getElementById(formId);
    if (!form) {
        console.error(`Form with ID ${formId} not found.`);
        return;
    }

    // For Similar Match Search, handle differently
    if (formId.includes('similar')) {
        handleSimilarSearch(formId);
        return;
    }

    // Only proceed with the mockup functionality for Exact Match Search
    showLoadingPopup();

    let formData = {};
    
    // Handle different form structures
    if (formId === 'ideb-request-form-individual-similar' || formId === 'ideb-request-form-individual-exact') {
        // Individual forms (tabbed)
        const isSimilar = formId.includes('similar');
        const suffix = isSimilar ? '_individual_similar' : '_individual_exact';
        
        formData = {
            nomor_referensi_pengguna: form.querySelector('#nomor_referensi_pengguna' + suffix).value,
            tujuan_penggunaan: form.querySelector('#tujuan_penggunaan' + suffix).value,
            nomor_identitas: form.querySelector('#nomor_identitas' + suffix).value,
            permintaan_fasilitas_outstanding: form.querySelector('#permintaan_fasilitas_outstanding' + suffix) ? 
                form.querySelector('#permintaan_fasilitas_outstanding' + suffix).checked : false,
            search_type: form.querySelector('input[name="search_type' + suffix + '"]:checked').value,
            request_type: "individual",
            jenis_identitas: form.querySelector('#jenis_identitas' + suffix).value
        };
    } else if (formId === 'ideb-request-form-badan-usaha-similar' || formId === 'ideb-request-form-badan-usaha-exact') {
        // Corporate forms (tabbed)
        const isSimilar = formId.includes('similar');
        const suffix = isSimilar ? '_badan_usaha_similar' : '_badan_usaha_exact';
        
        formData = {
            nomor_referensi_pengguna: form.querySelector('#nomor_referensi_pengguna' + suffix).value,
            tujuan_penggunaan: form.querySelector('#tujuan_penggunaan' + suffix).value,
            nomor_identitas: form.querySelector('#nomor_identitas' + suffix).value,
            permintaan_fasilitas_outstanding: form.querySelector('#permintaan_fasilitas_outstanding' + suffix) ? 
                form.querySelector('#permintaan_fasilitas_outstanding' + suffix).checked : false,
            search_type: form.querySelector('input[name="search_type' + suffix + '"]:checked').value,
            request_type: "corporate"
        };
    }

    // Simulate a 3-second delay
    setTimeout(() => {
        fetch('/api/requests', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        })
        .then(response => {
            if (!response.ok) {
                // If response is not ok, throw an error to be caught by the .catch block
                return response.json().then(err => { throw new Error(err.error || 'Unknown server error') });
            }
            return response.json();
        })
        .then(data => {
            console.log('Success:', data);
            hideLoadingPopup();
            alert('Data ditemukan, silahkan cek di Daftar Permintaan'); // Show success message
            
            // Determine which list to refresh based on the form submitted
            if (formId.includes("individual")) {
                loadContent('debitur-individual.html');
            } else if (formId.includes("badan-usaha")) {
                loadContent('badan-usaha.html');
            }
        })
        .catch((error) => {
            console.error('Error:', error);
            hideLoadingPopup();
            alert('Error submitting request: ' + error.message);
        });
    }, 3000); // 3-second delay
}

function loadRequests(page) {
    // Check if there are similar search results in sessionStorage
    if (sessionStorage.getItem('similarSearchResults')) {
        // If we're on the badan-usaha.html page, display the popup
        if (page === 'badan-usaha.html') {
            displaySimilarResultsPopup();
        }
        return;
    }
    
    let apiUrl = '';
    if (page === 'debitur-individual.html') {
        apiUrl = '/api/requests/individual';
    } else if (page === 'badan-usaha.html') {
        apiUrl = '/api/requests/corporate';
    }

    fetch(apiUrl)
        .then(response => response.json())
        .then(data => {
            const tbody = document.querySelector(`#${page.split('.')[0]} table tbody`);
            if (tbody) {
                tbody.innerHTML = ''; // Clear existing data
                data.forEach(req => {
                    let row = '';
                    if (page === 'debitur-individual.html') {
                        row = `<tr>
                            <td>${req.nomor_referensi_pengguna}</td>
                            <td>${req.tujuan_penggunaan}</td>
                            <td>${req.jenis_identitas}</td>
                            <td>${req.nomor_identitas}</td>
                            <td>${req.status_aksi === 'Dalam Proses' ? 'Dalam Proses' : `<a href="#" onclick="viewDetail(${req.ID})">Lihat Detail</a>`}</td>
                        </tr>`;
                    } else if (page === 'badan-usaha.html') {
                        row = `<tr>
                            <td>${req.nomor_referensi_pengguna}</td>
                            <td>${req.tujuan_penggunaan}</td>
                            <td>${req.nomor_identitas}</td>
                            <td>${req.status_aksi === 'Dalam Proses' ? 'Dalam Proses' : `<a href="#" onclick="viewDetail(${req.ID})">Lihat Detail</a>`}</td>
                        </tr>`;
                    }
                    tbody.innerHTML += row;
                });
            }
        });
}

function viewDetail(id) {
    window.open(`/api/generate-pdf?id=${id}`, '_blank');
}

// Handle similar search functionality
function handleSimilarSearch(formId) {
    const form = document.getElementById(formId);
    if (!form) {
        console.error(`Form with ID ${formId} not found.`);
        return;
    }

    showLoadingPopup();

    let formData = {};
    
    // Handle different form structures for similar search
    if (formId === 'ideb-request-form-badan-usaha-similar') {
        // Get form values
        const nomorIdentitas = form.querySelector('#nomor_identitas_badan_usaha_similar').value;
        const namaBadanUsaha = form.querySelector('#nama_badan_usaha_similar').value;
        const tempatPendirian = form.querySelector('#tempat_pendirian_similar').value;
        const tanggalPendirian = form.querySelector('#tanggal_pendirian_similar').value;
        
        // Validate that Nama Badan Usaha is filled
        if (!namaBadanUsaha) {
            hideLoadingPopup();
            alert('Nama Badan Usaha wajib diisi');
            return;
        }
        
        formData = {
            nomor_referensi_pengguna: form.querySelector('#nomor_referensi_pengguna_badan_usaha_similar').value,
            tujuan_penggunaan: form.querySelector('#tujuan_penggunaan_badan_usaha_similar').value,
            nomor_identitas: nomorIdentitas,
            nama_badan_usaha: namaBadanUsaha,
            tempat_pendirian: tempatPendirian,
            tanggal_pendirian: tanggalPendirian,
            permintaan_fasilitas_outstanding: form.querySelector('#permintaan_fasilitas_outstanding_badan_usaha_similar') ? 
                form.querySelector('#permintaan_fasilitas_outstanding_badan_usaha_similar').checked : false,
            search_type: form.querySelector('input[name="search_type_badan_usaha_similar"]:checked').value,
            request_type: "corporate"
        };
    }

    // Send request to the similar search endpoint
    fetch('/api/similar-search', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
    .then(response => {
        if (!response.ok) {
            // If response is not ok, throw an error to be caught by the .catch block
            return response.json().then(err => { throw new Error(err.error || 'Unknown server error') });
        }
        return response.json();
    })
    .then(data => {
        hideLoadingPopup();
        
        // Check if data is available
        if (data && data.data && data.data.length > 0) {
            // Store the results in sessionStorage
            sessionStorage.setItem('similarSearchResults', JSON.stringify(data));
            // Redirect to badan-usaha.html
            loadContent('badan-usaha.html');
        } else {
            // No results found
            alert('Data tidak ditemukan');
        }
    })
    .catch((error) => {
        console.error('Error:', error);
        hideLoadingPopup();
        alert('Data tidak ditemukan');
    });
}

// Display similar search results in a popup
function displaySimilarResultsPopup() {
    // Check if there are similar search results in sessionStorage
    const results = sessionStorage.getItem('similarSearchResults');
    if (!results) {
        return;
    }

    const data = JSON.parse(results);

    // Create the modal HTML
    const modalHTML = `
        <div class="modal fade" id="similarResultsModal" tabindex="-1" aria-labelledby="similarResultsModalLabel" aria-hidden="true">
            <div class="modal-dialog modal-xl">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title" id="similarResultsModalLabel">Hasil Pencarian Similar Match</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                    </div>
                    <div class="modal-body">
                        <table class="table table-striped">
                            <thead>
                                <tr>
                                    <th>Kode Jenis Pelapor</th>
                                    <th>Kode Pelapor</th>
                                    <th>Nama Debitur</th>
                                    <th>Nomor Identitas</th>
                                    <th>Aksi</th>
                                </tr>
                            </thead>
                            <tbody id="similarResultsTableBody">
                                <!-- Rows will be populated dynamically -->
                            </tbody>
                        </table>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Tutup</button>
                    </div>
                </div>
            </div>
        </div>
    `;

    // Add the modal to the document
    const modalContainer = document.createElement('div');
    modalContainer.innerHTML = modalHTML;
    document.body.appendChild(modalContainer);

    // Populate the table with results
    populateSimilarResultsTable(data.data);

    // Show the modal
    const modal = new bootstrap.Modal(document.getElementById('similarResultsModal'));
    modal.show();

    // Clean up the modal when it's closed
    document.getElementById('similarResultsModal').addEventListener('hidden.bs.modal', function () {
        document.body.removeChild(modalContainer);
    });
}

// Populate the similar results table
function populateSimilarResultsTable(results) {
    const tableBody = document.getElementById('similarResultsTableBody');
    if (!tableBody) return;

    tableBody.innerHTML = '';

    results.forEach((result, index) => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${result.kode_jenis_pelapor}</td>
            <td>${result.kode_pelapor}</td>
            <td>${result.nama_debitur}</td>
            <td>${result.nomor_identitas}</td>
            <td>
                <button class="btn btn-primary btn-sm lihat-status-btn" 
                        data-kode-jenis-pelapor="${result.kode_jenis_pelapor}"
                        data-kode-pelapor="${result.kode_pelapor}"
                        data-nama-debitur="${result.nama_debitur}"
                        data-nomor-identitas="${result.nomor_identitas}">
                    Lihat Status
                </button>
            </td>
        `;
        tableBody.appendChild(row);
    });

    // Add event listeners to the "Lihat Status" buttons
    document.querySelectorAll('.lihat-status-btn').forEach(button => {
        button.addEventListener('click', handlePrintPdf);
    });
}

// Handle PDF generation for a specific result
function handlePrintPdf(event) {
    const button = event.target;
    
    // Get data from the button's data attributes
    const data = {
        kode_jenis_pelapor: button.dataset.kodeJenisPelapor,
        kode_pelapor: button.dataset.kodePelapor,
        nama_debitur: button.dataset.namaDebitur,
        nomor_identitas: button.dataset.nomorIdentitas
    };

    // Show loading indicator
    showLoadingPopup();
    
    // Send request to the backend to generate PDF
    fetch('/api/generate-similar-pdf', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.blob();
    })
    .then(blob => {
        // Create a temporary URL for the PDF blob
        const url = window.URL.createObjectURL(blob);
        // Create a temporary link element to trigger the download
        const a = document.createElement('a');
        a.href = url;
        a.download = 'similar_search_result.pdf';
        document.body.appendChild(a);
        a.click();
        // Clean up
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    })
    .catch(error => {
        console.error('Error generating PDF:', error);
        alert('Error generating PDF. Please try again.');
    })
    .finally(() => {
        // Hide loading indicator
        hideLoadingPopup();
        // Close the modal
        const modal = bootstrap.Modal.getInstance(document.getElementById('similarResultsModal'));
        modal.hide();
    });
}
