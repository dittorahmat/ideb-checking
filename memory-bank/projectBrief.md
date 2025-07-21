saya ingin membuat v0 (mockup) dari aplikasi ideb checking. tujuan dari mockup ini adalah untuk demo presales, jadi fungsionalnya cukup login dengan email password (dummy saja dulu, sembarang email password bisa masuk), bisa search data, menampilkan di data permintaan, apabila data sudah available data detailnya, bisa print as PDF. meskipun baru mockup, saya sudah siapkan struktur table data, yang kurang lebih 80% siap, yang akan digunakan untuk v1 (saya juga sudah siapkan sample input 1 record JSON, yang kemudian bisa diparsing ke table database).  tech stack yang saya pilih frontend html, javascript, bootstrap, css biasa. backend menggunakan go (saya tidak tahu, perlu framework atau tidak, tolong berikan saran). untuk database saya pilih SQLite dulu saja di v0, tapi v1 akan pakai PostgreSQL (saya tidak tahu, perlu caching redis atau tidak, tolong berikan saran). apa perlu rabbitmq / kafka juga ya untuk queueing / async processing, atau cukup pakai native go? tolong berikan saran. untuk v0, aplikasi ini cukup dijalankan di localhost saja dulu, tapi v1, rencananya mungkin akan deploy ke Azure App Service (monorepo frontend + backend)

aplikasi ini terdiri dari 2 bagian, tampilan frontend dan api / middleware, yang bisa diakses umum. ideb akan mengambil data dari API SLIK OJK, berupa semacam informasi kredit / pinjaman dari individual / badan usaha. OJK kemudian akan mengumpulkan data, secara realtime, dari berbagai lembaga keuangan / bank / pemberi kredit. 

ini bisa memakan waktu, bisa hitungan menit, bahkan bisa sampai 1 jam (info yang saya dapat seperti itu). jadi aplikasi IDEB yang saya buat ini akan menjadi middleware ke SLIK OJK. user yang query ke aplikasi IDEB ini bisa mencari data  informasi kredit yang sudah aplikasi kita simpan sebelumnya atau bisa mencari langsung OJK via aplikasi kita, kemudian data tersebut akan kita simpan ke database aplikasi IDEB. 

ini sidebar dari aplikasi Ideb checking, kalau ada indentasi itu artinya nested ya, bisa collapsible juga
Ideb Checking (main page)		
	Input Permintaan IDeb	
	Daftar Permintaan IDeb	
		Debitur Individual
		Badan Usaha
	Dashboard	
	Parameter	
		User & Password API
		valid Token
		LDAP
	User Management	

untuk input aplikasi , ada beberapa field yang diperlukan, semuanya wajib diisi. 
nomor_referensi_pengguna
tujuan_penggunaan
jenis_identitas
nomor_identitas
permintaan_fasilitas_outstanding
kemudian ada pilihan lagi untuk cari di internal database atau live SLIK OJK 

untuk output aplikasi, akan ada menu daftar permintaan. menu ini akan menampilkan semua permintaan yang telah diajukan user, baik itu sudah selesai atau belum (ingat, data dari API SLIK . sama dengan kolom input, dengan tambahan 1 kolom
status_aksi
di status aksi ini, apabila data belum ready, akan menampilkan "Dalam Proses" , apabila data sudah ready, akan menampikan link "Lihat Detail", yang apabila diklik akan print detail data dalam PDF.
