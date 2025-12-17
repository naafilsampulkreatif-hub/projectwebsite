# 📑 Dokumentasi Teknis Proyek E-commerce Night Stalkers

**Target Audience:** Tim Developer  
**Status Proyek:** Perencanaan (Planning)  
**Versi Dokumen:** 1.0  
**Tanggal Revisi:** 16 Desember 2025  

---

## I. Deskripsi Fitur Secara Deskriptif

Dokumentasi ini mencakup fungsionalitas utama untuk platform e-commerce yang menjual komik dan novel dari IP Night Stalkers. Sistem ini menggunakan pendekatan autentikasi terpisah untuk Pelanggan (User) dan Administrator (Admin).

### 1. Sistem Autentikasi & Akun
* **Login Admin:** Akses khusus melalui portal terpisah (Dashboard) untuk mengelola operasional toko, inventaris, dan memantau transaksi.
* **Login/Registrasi User:** Memungkinkan pelanggan membuat akun untuk menyimpan riwayat pesanan, detail pengiriman, dan mempercepat proses transaksi di masa depan.
* **Logout:** Fitur keamanan untuk menghancurkan sesi aktif (session) atau token (JWT) pada sisi klien dan server guna melindungi akun.

### 2. Manajemen Toko & Produk
* **Katalog Produk:** Menampilkan daftar buku komik dan novel secara dinamis. Pengguna dapat melakukan pencarian dan filter berdasarkan genre atau harga.
* **Detail Produk:** Informasi komprehensif mencakup sinopsis, spesifikasi buku, harga, dan status stok real-time yang terhubung ke database MySQL.

### 3. Transaksi & Keranjang Belanja
* **Keranjang Belanja (Shopping Cart):** Fitur untuk menampung produk pilihan, mengubah kuantitas, atau menghapus item sebelum checkout.
* **Checkout & Pembayaran:** Proses finalisasi pembelian. Sistem akan memvalidasi data pengiriman, menghitung biaya, dan mengintegrasikan logika stok untuk memastikan ketersediaan barang saat pembayaran dilakukan.

---

## II. Alur Kerja (Flowchart)

Berikut adalah logika alur kerja untuk proses utama dalam sistem.

### A. Alur Login (User & Admin)
1.  **Mulai:** Pengguna mengakses halaman login.
2.  **Input:** Memasukkan Email dan Password.
3.  **Validasi (Backend):** * Sistem mencocokkan kredensial dengan database.
    * Jika valid: Generate Token (JWT) -> Arahkan ke Dashboard (Admin) atau Home (User).
    * Jika tidak valid: Tampilkan pesan kesalahan -> Kembali ke input login.
4.  **Selesai.**

### B. Alur Checkout Pesanan
1.  **Mulai:** Pengguna menekan tombol "Checkout" dari halaman keranjang.
2.  **Verifikasi Identitas:** * Jika sudah login: Data diri dimuat otomatis.
    * Jika belum login: Pengguna diminta mengisi formulir data diri (Nama, Email, Alamat).
3.  **Validasi Stok:** Backend melakukan pengecekan ketersediaan item di database MySQL.
4.  **Finalisasi Transaksi:** Pengguna memilih metode pembayaran dan melakukan konfirmasi.
5.  **Update Sistem:** Sistem mencatat transaksi baru dan secara otomatis mengurangi jumlah stok di database.
6.  **Selesai.**

---

## III. Full Stack & Spesifikasi Teknik

Spesifikasi berikut adalah standar wajib bagi tim pengembang untuk menjaga konsistensi dan performa aplikasi.

### 1. Teknologi Pengembangan (Stack)
* **Frontend:** Vue.js menggunakan build tool **Vite v7.3.0**.
* **Backend Utama:** **Go (Golang) versi 1.22+** (Digunakan untuk logika bisnis utama dan API berperforma tinggi).
* **Backend Pendukung:** **Node.js versi 18+** (Digunakan untuk microservices atau scripting tambahan jika diperlukan).
* **Database:** **MySQL versi 8.0+** (Database relasional untuk menjaga integritas data transaksi).

### 2. Spesifikasi Software & Tools
* **Code Editor:** Visual Studio Code.
* **API Testing:** Postman atau Insomnia.
* **Database Management:** MySQL Workbench atau DBeaver.
* **Version Control:** Git (GitHub/GitLab).

### 3. Spesifikasi Hardware (Minimal Server)
* **Sistem Operasi:** Ubuntu 22.04 LTS (Direkomendasikan untuk lingkungan produksi).
* **CPU:** Minimal 2 Cores.
* **RAM:** Minimal 4 GB.
* **Penyimpanan:** 40 GB SSD.

---
*Dokumen ini bersifat internal. Segala perubahan pada alur atau stack teknologi harus melalui persetujuan Lead Developer.*

<img width="2626" height="2384" alt="codetoflow" src="https://github.com/user-attachments/assets/62a524c6-595d-49ee-b851-36873d4f0e75" />
