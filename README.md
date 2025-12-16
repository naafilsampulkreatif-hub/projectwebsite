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

#### Pseudo Code:

```pseudo-code
FUNCTION HandleUserLogin(email, password):
    // 1. Verifikasi Kredensial
    user = DB.FindUserByEmail(email)
    IF user IS NOT FOUND OR NOT Password.Verify(password, user.hashed_password):
        RETURN ERROR "Kredensial tidak valid"
    END IF

    // 2. Buat Sesi
    token = Auth.GenerateJWT(user.id)
    Session.SaveToken(token)

    // 3. Arahkan
    RETURN SUCCESS, token, RedirectTo('/profile')
END FUNCTION
B. Alur Fitur: Checkout (Flowchart Teks)START: Pengguna menekan tombol "Checkout" dari Keranjang Belanja.DECISION (Kondisi 1): Apakah User Sudah Login?IF YA: Isi otomatis data diri dari profil, Lanjut ke Poin 3.IF TIDAK: Pengguna diwajibkan Login atau Registrasi. DECISION (Kondisi 2): Apakah User Berhasil Login/Daftar?IF YA: Lanjut ke Poin 3.IF TIDAK: Kembali ke Poin 2.PROCESS: Pengguna mengisi/memverifikasi Alamat Pengiriman dan memilih Metode Pengiriman (API Pengiriman dihitung).PROCESS: Pengguna memilih Metode Pembayaran.DECISION (Kondisi 3): Apakah Pengguna Menekan "Bayar Sekarang"?IF YA: Lanjut ke Poin 6.IF TIDAK: Kembali ke Keranjang/Batal.PROCESS: Sistem memvalidasi data, mengurangi stok produk, membuat entri Pesanan baru di DB dengan Status = 'Pending Payment'.PROCESS: Kirim Email Konfirmasi Pesanan dan Detail Pembayaran.PROCESS: Arahkan pengguna ke Halaman Terima Kasih / Instruksi Pembayaran.END.Pseudo Code:Cuplikan kodeFUNCTION HandleCheckout(cart_data, user_id, shipping_address, payment_method):
    // 1. Validasi Keranjang
    IF Cart.IsEmpty(cart_data):
        RETURN ERROR "Keranjang kosong"

    // 2. Hitung Total Biaya
    total_price = Cart.CalculateTotal(cart_data) + Shipping.CalculateCost(shipping_address)

    // 3. Kurangi Stok
    IF Product.DecreaseStock(cart_data) IS FAILED:
        RETURN ERROR "Stok tidak mencukupi"

    // 4. Buat Pesanan
    order = DB.CreateOrder(user_id, total_price, shipping_address, payment_method, status='PENDING')

    // 5. Notifikasi
    Email.SendOrderConfirmation(order.id, user.email)

    RETURN SUCCESS, order.id, RedirectTo('/order/thank-you')
END FUNCTION
III. Full Stack/Teknik yang Akan DigunakanA. Tumpukan Teknologi (Full Stack)LapisanTeknologi UtamaBahasa/VersiFramework/Library UtamaFrontend (UI)Vue.jsJavaScript/TypeScriptVue.js 3, Vite v7.3.0Backend (API)Go (Golang)Go 1.22+Native Go (atau Gin/Echo jika diperlukan)DatabaseMySQLMySQL 8.0+SQL Driver (misalnya, go-sql-driver/mysql)Secondary BackendNode.jsNode.js 18+Express.js (dapat digunakan untuk layanan mikro non-kritis atau worker)B. Spesifikasi Software dan ToolsSistem Operasi Pengembangan: Linux/macOS/Windows Subsystem for Linux (WSL).Version Control: Git & GitHub/GitLab.Editor: Visual Studio Code (VS Code) dengan ekstensi Go dan Vue.Package Manager:Go: go modNode.js/Vue: npm atau yarnTesting: Jest (Frontend), Go Testing Package (Backend).C. Spesifikasi Hardware (Lingkungan Produksi)KomponenSpesifikasi Minimum AwalCatatanApplication Server (Go)2 vCPU, 4 GB RAM, 80 GB SSDFokus pada efisiensi Go untuk kecepatan pemrosesan request API.Database Server (MySQL)4 vCPU, 8 GB RAM, 200 GB SSD (diutamakan SSD performa tinggi)MySQL 8.0 membutuhkan sumber daya yang memadai, terutama untuk optimasi kueri kompleks.CDN & DNSCloudflare / AWS CloudFrontWajib untuk caching assets (gambar produk, JS/CSS) dan performa global.Sistem Operasi ServerUbuntu Server LTS (Long Term Support)Pilihan umum karena stabilitas dan dukungan komunitas.
