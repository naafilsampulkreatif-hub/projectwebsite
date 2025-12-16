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









# Molla Ecommerce (Go + Vue)

A responsive ecommerce website built with a Golang backend and Vue.js frontend, featuring a multi-level admin dashboard, JWT authentication, and full shopping cart functionality.

## Prerequisites

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
3. Insert an admin user manually into the `users` table if you want to access the Admin Dashboard immediately (role='admin').

### Backend

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Configure environment:
   Copy `.env.example` to `.env` and update your MySQL password:
   ```bash
   cp .env.example .env
   # Edit .env and set DB_PASS=your_real_password
   ```

4. Run the server:
   ```bash
   go run main.go
   ```
   Server runs on `http://localhost:8080`.

### Frontend

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the development server:
   ```bash
   npm run dev
   ```
   App runs on `http://localhost:5173`.

## Features

- **Authentication**: Register/Login with JWT.
- **RBAC**: Admin vs Customer roles.
- **Admin Dashboard**: Manage products (CRUD).
- **Shopping**: Browse products, Add to Cart, Checkout.
- **Design**: Responsive UI with TailwindCSS.
