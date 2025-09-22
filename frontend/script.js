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

    // For Similar Match Search, just show an alert and return
    if (formId.includes('similar')) {
        alert('Similar Match Search is not implemented yet.');
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
