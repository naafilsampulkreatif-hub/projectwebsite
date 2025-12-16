# Dokumentasi Teknis Proyek E-commerce Night Stalkers

**Target Audience:** Tim Developer
**Versi Dokumen:** 1.0
**Tanggal Revisi:** 16 Desember 2025

---

## I. Deskripsi Fitur Secara Deskriptif

Dokumen ini menjelaskan fungsionalitas utama website E-commerce Night Stalkers, yang melayani pembelian buku komik dan novel dengan dua jenis pengguna: Pelanggan (User) dan Administrator (Admin).

### A. Fitur Utama Pelanggan (User)

| Fitur | Deskripsi Fungsional | Tujuan Utama |
| :--- | :--- | :--- |
| **Login/Registrasi User** | Pelanggan dapat membuat akun baru atau masuk menggunakan kredensial email/password mereka. Akun diperlukan untuk menyimpan riwayat pesanan dan detail pengiriman. | Membangun basis pelanggan, memungkinkan pembelian berulang yang cepat, dan menyimpan data riwayat. |
| **Toko (Katalog Produk)** | Menampilkan daftar lengkap komik dan novel. Pengguna dapat mencari, memfilter (berdasarkan genre, harga), dan melihat detail setiap produk (sinopsis, spesifikasi, stok). | Memudahkan penemuan dan eksplorasi produk. |
| **Keranjang Belanja (Shopping Cart)** | Pengguna dapat menambahkan, menghapus, dan mengubah kuantitas produk sebelum melanjutkan ke proses pembayaran. Data keranjang terikat pada sesi atau akun pengguna. | Memungkinkan pengguna meninjau dan memfinalisasi pesanan. |
| **Checkout & Pembayaran** | Proses penyelesaian transaksi. Pengguna mengisi alamat pengiriman, memilih metode pengiriman, dan memilih metode pembayaran (transfer bank, e-wallet, dll.). Setelah sukses, sistem menghasilkan ID Pesanan unik. | Menyelesaikan proses konversi (pembelian). |
| **Logout User** | Mengakhiri sesi pengguna dari aplikasi. | Menjaga keamanan akun. |

### B. Fitur Utama Administrator (Admin)

| Fitur | Deskripsi Fungsional | Tujuan Utama |
| :--- | :--- | :--- |
| **Login Admin** | Akses ke antarmuka terpisah (Dashboard) dengan kredensial khusus Admin. | Mengamankan akses ke manajemen operasional dan data sensitif. |
| **Dashboard Admin** | Halaman pusat untuk mengelola seluruh operasional situs (CRUD: Create, Read, Update, Delete). | Mengelola produk, pesanan, dan pengguna secara efisien. |
| **Logout Admin** | Mengakhiri sesi Admin dari Dashboard. | Menjaga keamanan sistem. |

---

## II. Alur/Flow Chart

Berikut adalah deskripsi alur untuk dua fitur kritis (Login User dan Checkout) dalam format teks dan kode pseudo.

### A. Alur Fitur: Login User (Flowchart Teks)

1.  **START:** Pengguna mengakses halaman Login.
2.  **PROCESS:** Pengguna memasukkan Email dan Password.
3.  **DECISION (Kondisi 1):** Apakah Kredensial (Email/Password) Valid?
    * **IF YA:** Lanjut ke Poin 4.
    * **IF TIDAK:** Tampilkan pesan error ("Kredensial tidak valid"), kembali ke Poin 2.
4.  **PROCESS:** Sistem membuat Session/Token Autentikasi (JWT).
5.  **PROCESS:** Arahkan pengguna ke Halaman Utama (Homepage) atau Dashboard Profil.
6.  **END.**

<img width="2626" height="2384" alt="codetoflow" src="https://github.com/user-attachments/assets/62a524c6-595d-49ee-b851-36873d4f0e75" />

NightStalkers Ecommerce (Go + Vue)
Situs web e-commerce responsif yang dibangun dengan backend Golang dan frontend Vue.js, dilengkapi dengan dashboard admin multi-level, autentikasi JWT, dan fungsi keranjang belanja lengkap.

Prasyarat
Go 1.22+

Node.js 18+

MySQL 8.0+

Pengaturan (Setup)
Database
Buat database MySQL dengan nama ecommerce.

Jalankan skrip inisialisasi:

Bash

mysql -u root -p ecommerce < backend/schema.sql
Tambahkan pengguna admin secara manual ke dalam tabel users jika Anda ingin segera mengakses Admin Dashboard (setel role='admin').

Backend
Masuk ke direktori backend:

Bash

cd backend
Instal dependensi:

Bash

go mod tidy
Konfigurasi lingkungan (environment): Salin .env.example menjadi .env dan perbarui kata sandi MySQL Anda:

Bash

cp .env.example .env
# Edit .env dan setel DB_PASS=kata_sandi_asli_anda
Jalankan server:

Bash

go run main.go
Server akan berjalan di http://localhost:8080.

Frontend
Masuk ke direktori frontend:

Bash

cd frontend
Instal dependensi:

Bash

npm install
Jalankan server pengembangan:

Bash

npm run dev
Aplikasi akan berjalan di http://localhost:5173.

Fitur
Autentikasi: Registrasi/Login menggunakan JWT.

RBAC (Role-Based Access Control): Perbedaan peran antara Admin dan Pelanggan.

Admin Dashboard: Mengelola produk (CRUD: Tambah, Baca, Ubah, Hapus).

Belanja: Menjelajahi produk, Tambah ke Keranjang, dan Checkout.

Desain: Antarmuka pengguna (UI) responsif menggunakan TailwindCSS.
