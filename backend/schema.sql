CREATE TABLE IF NOT EXISTS requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor_referensi_pengguna TEXT NOT NULL,
    tujuan_penggunaan TEXT NOT NULL,
    jenis_identitas TEXT NOT NULL,
    nomor_identitas TEXT NOT NULL,
    permintaan_fasilitas_outstanding BOOLEAN NOT NULL,
    search_type TEXT NOT NULL, -- internal or live
    status_aksi TEXT NOT NULL DEFAULT 'Dalam Proses',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
