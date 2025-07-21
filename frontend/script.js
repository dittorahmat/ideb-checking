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
    const form = document.getElementById('ideb-request-form');
    if (form) {
        form.addEventListener('submit', function(event) {
            event.preventDefault();
            submitIdebRequest();
        });
    }
}

function submitIdebRequest() {
    const formData = {
        nomor_referensi_pengguna: document.getElementById('nomor_referensi_pengguna').value,
        tujuan_penggunaan: document.getElementById('tujuan_penggunaan').value,
        jenis_identitas: document.getElementById('jenis_identitas').value,
        nomor_identitas: document.getElementById('nomor_identitas').value,
        permintaan_fasilitas_outstanding: document.getElementById('permintaan_fasilitas_outstanding').checked,
        search_type: document.querySelector('input[name="search_type"]:checked').value
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
        loadContent('debitur-individual.html'); // Refresh the list
    })
    .catch((error) => {
        console.error('Error:', error);
        alert('Error submitting request.');
    });
}

function loadRequests(page) {
    fetch('/api/requests')
        .then(response => response.json())
        .then(data => {
            const tbody = document.querySelector(`#${page.split('.')[0]} table tbody`);
            if (tbody) {
                tbody.innerHTML = ''; // Clear existing data
                data.forEach(req => {
                    const row = `<tr>
                        <td>${req.nomor_referensi_pengguna}</td>
                        <td>${req.tujuan_penggunaan}</td>
                        <td>${req.jenis_identitas}</td>
                        <td>${req.nomor_identitas}</td>
                        <td>${req.status_aksi === 'Dalam Proses' ? 'Dalam Proses' : `<a href="#" onclick="viewDetail(${req.id})">Lihat Detail</a>`}</td>
                    </tr>`;
                    tbody.innerHTML += row;
                });
            }
        });
}

function viewDetail(id) {
    // For v0, we can just show an alert.
    // In a real app, this would fetch the details and generate a PDF.
    alert(`Viewing details for request ID: ${id}`);
}
