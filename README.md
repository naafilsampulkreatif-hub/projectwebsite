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


# Night Stalkers Ecommerce (Go + Vue)

Sebuah situs web e-commerce responsif yang dibangun dengan backend Golang dan frontend Vue.js, menampilkan dasbor admin multi-level, otentikasi JWT, dan fungsionalitas keranjang belanja lengkap.

## Prasyarat

- Go 1.22+
- Node.js 18+
- MySQL 8.0+

## Setup

### Database

1. Create a MySQL database named `ecommerce`.
2. Run the initialization script:
   ```bash
   mysql -u root -p ecommerce < backend/schema.sql
   ```
3. Masukkan pengguna admin secara manual ke dalam `users` tabel jika Anda ingin mengakses Dasbor Admin segera (role='admin').

### Backend

1. Masuk ke direktori backend:
   ```bash
   cd backend
   ```
2. Instal dependensi:
   ```bash
   go mod tidy
   ```
3. Konfigurasi lingkungan:
   Copy `.env.example` to `.env` dan perbarui kata sandi MySQL Anda:
   ```bash
   cp .env.example .env
   # Edit .env and set DB_PASS=your_real_password
   ```

4. Jalankan servernya:
   ```bash
   go run main.go
   ```
   Server berjalan `http://localhost:8080`.

### Frontend

1. Masuk ke direktori frontend:
   ```bash
   cd frontend
   ```
2. Instal dependensi:
   ```bash
   npm install
   ```
3. Jalankan server pengembangan:
   ```bash
   npm run dev
   ```
   Aplikasi berjalan `http://localhost:5173`.

## Fitur-fitur

- **Otentikasi**: Register/Login dengan JWT.
- **RBAC**: Peran Admin vs Pelanggan.
- **Admin Dashboard**: Kelola produk (CRUD).
- **belanja**: Telusuri produk, Tambahkan ke Keranjang, Lakukan Pembayaran.
- **desain**: Responsive UI dengan TailwindCSS.
