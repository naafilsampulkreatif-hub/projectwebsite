# 📑 Dokumentasi Teknis Proyek E-commerce Night Stalkers

**Target Audience:** Tim Developer  
**Status Proyek:** Perencanaan (Planning)  
**Versi Dokumen:** 1.0  

---

## I. Deskripsi Fitur Secara Deskriptif

Dokumentasi ini mencakup fungsionalitas utama untuk platform e-commerce yang menjual komik dan novel dari IP Night Stalkers. Sistem ini menggunakan pendekatan autentikasi ganda untuk User dan Admin.

### 1. Fitur Autentikasi (Login & Logout)
* **Login Admin:** Akses khusus melalui portal terpisah untuk mengelola konten toko dan memantau transaksi.
* **Login User:** Memungkinkan pelanggan masuk ke akun mereka untuk mengakses riwayat pembelian dan mempercepat proses checkout.
* **Logout:** Fitur keamanan untuk menghancurkan sesi (session) atau token (JWT) aktif baik pada sisi klien maupun server.

### 2. Fitur Toko (Katalog Produk)
* Menampilkan daftar buku komik dan novel secara dinamis.
* Dilengkapi dengan detail produk yang mencakup sinopsis, harga, dan status stok yang ditarik langsung dari database MySQL.

### 3. Fitur Keranjang (Shopping Cart)
* Memungkinkan User untuk menampung sementara produk yang ingin dibeli.
* Sistem menghitung total harga secara real-time sebelum berlanjut ke tahap pembayaran.

### 4. Fitur Checkout
* Proses finalisasi pembelian di mana User memasukkan data diri (jika belum login) dan alamat pengiriman.
* Integrasi logika stok untuk memastikan produk yang dibayar masih tersedia di gudang.

---

## II. Alur Kerja (Flowchart)

Berikut adalah logika alur fitur utama dalam representasi teks dan kode logika.

### A. Alur Login (User & Admin)
1. **Mulai:** User/Admin membuka halaman login.
2. **Input:** Memasukkan Email dan Password.
3. **Validasi (Backend):** - Jika data cocok dengan database -> Generate Token -> Masuk ke Dashboard/Home.
    - Jika data salah -> Tampilkan pesan error -> Kembali ke halaman login.
4. **Selesai.**

**Logika Kode (Pseudo-code):**
```pseudo
IF input_email AND input_password MATCH database_record:
    CREATE session_token
    IF role == "admin":
        REDIRECT to /admin-dashboard
    ELSE:
        REDIRECT to /home
ELSE:
    DISPLAY "Kredensial Salah"
    RELOAD login_page
B. Alur Checkout Pesanan
Mulai: User klik "Checkout" di halaman keranjang.

Cek Login: - Jika sudah Login: Ambil data pribadi dari database.

Jika belum Login: Tampilkan formulir data diri (Email, Nama, Alamat).

Validasi Stok: Backend mengecek ketersediaan buku di MySQL.

Pembayaran: User memilih metode pembayaran dan konfirmasi.

Update: Sistem mengurangi stok di database dan mencatat transaksi.

Selesai.

Logika Kode (Pseudo-code):

Cuplikan kode

FUNCTION process_checkout(user_data, cart_items):
    FOR item IN cart_items:
        IF item.qty > database.stock:
            RETURN "Stok Tidak Cukup"
    
    database.create_order(user_data, cart_items)
    database.update_stock(minus, cart_items)
    SEND confirmation_email
    RETURN "Checkout Berhasil"
III. Full Stack & Spesifikasi Teknik
Spesifikasi ini wajib dipenuhi oleh tim developer untuk menjaga konsistensi performa aplikasi.

1. Teknologi Pengembangan (Stack)
Frontend: Vue.js dengan build tool Vite v7.3.0.

Backend Utama: Go (Golang) versi 1.22+ (untuk performa tinggi dan konkurensi).

Backend Pendukung: Node.js versi 18+ (opsional untuk microservices atau scripting tambahan).

Database: MySQL versi 8.0+ (Relational database untuk integritas data transaksi).

### 2. Spesifikasi Software & Hardware
Development Software:

Code Editor: Visual Studio Code.

API Testing: Postman atau Insomnia.

Database Tool: MySQL Workbench atau DBeaver.

Hardware (Minimal Server):

CPU: 2 Cores.

RAM: 4 GB.

Storage: 40 GB SSD.

<img width="2626" height="2384" alt="codetoflow" src="https://github.com/user-attachments/assets/62a524c6-595d-49ee-b851-36873d4f0e75" />