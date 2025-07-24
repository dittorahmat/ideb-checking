document.addEventListener('DOMContentLoaded', (event) => {
    // App starts with login view
});

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
            if (page === 'input-permintaan.html') {
                setupFormListener();
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

function submitIdebRequest(formId) {
    const form = document.getElementById(formId);
    if (!form) {
        console.error(`Form with ID ${formId} not found.`);
        return;
    }

    const formData = {
        nomor_referensi_pengguna: form.querySelector('[id^="nomor_referensi_pengguna"]').value,
        tujuan_penggunaan: form.querySelector('[id^="tujuan_penggunaan"]').value,
        jenis_identitas: form.querySelector('[id^="jenis_identitas"]').value,
        nomor_identitas: form.querySelector('[id^="nomor_identitas"]').value,
        permintaan_fasilitas_outstanding: form.querySelector('[id^="permintaan_fasilitas_outstanding"]').checked,
        search_type: form.querySelector('input[name^="search_type"]:checked').value
    };

    fetch('/api/requests', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        alert('Request submitted successfully!');
        // Determine which list to refresh based on the form submitted
        if (formId.includes("individual")) {
            loadContent('debitur-individual.html');
        } else if (formId.includes("badan-usaha")) {
            loadContent('badan-usaha.html');
        }
    })
    .catch((error) => {
        console.error('Error:', error);
        alert('Error submitting request.');
    });
}

function loadRequests(page) {
    let apiUrl = '';
    if (page === 'debitur-individual.html') {
        apiUrl = '/api/getDebtorExactIndividual';
    } else if (page === 'badan-usaha.html') {
        apiUrl = '/api/getDebtorExactCorporate';
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
                            <td>${req.status_aksi === 'Dalam Proses' ? 'Dalam Proses' : `<a href="#" onclick="viewDetail(${req.id})">Lihat Detail</a>`}</td>
                        </tr>`;
                    } else if (page === 'badan-usaha.html') {
                        row = `<tr>
                            <td>${req.nomor_referensi_pengguna}</td>
                            <td>${req.tujuan_penggunaan}</td>
                            <td>${req.nomor_identitas}</td>
                            <td>${req.status_aksi === 'Dalam Proses' ? 'Dalam Proses' : `<a href="#" onclick="viewDetail(${req.id})">Lihat Detail</a>`}</td>
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
