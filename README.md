Dokumentasi Perencanaan Pembuatan Web Ecommerce Night Stalkers
==============================================================

1\. Gambaran Umum Proyek
------------------------

### 1.1 Deskripsi Proyek

**Night Stalkers** adalah platform e-commerce modern untuk penjualan buku komik dan novel yang mengutamakan kecepatan transaksi. Berbeda dengan e-commerce tradisional, sistem ini menggunakan pendekatan **Guest Checkout & Session-Based**. Pengunjung dapat langsung berbelanja tanpa perlu mendaftar akun (Login), namun data diri wajib diisi pada saat proses penyelesaian pesanan (Checkout).

### 1.2 Tujuan Pengembangan

*   Memberikan pengalaman belanja yang instan tanpa hambatan registrasi.
    
*   Mengelola keranjang belanja secara _real-time_ berbasis sesi pengunjung.
    
*   Menyediakan data transaksi yang lengkap bagi admin untuk keperluan evaluasi meskipun user tidak memiliki akun tetap.
    

### 1.3 Target Pengguna

*   **Pengunjung (Guest User):** Pengguna yang dapat langsung memilih produk dan melakukan checkout.
    
*   **Administrator:** Pengelola yang memiliki akses penuh untuk memantau data pesanan dan evaluasi performa toko.
    

2\. Alur Kerja / Flowchart (Sistem Per Sesi)
--------------------------------------------

### 2.1 Alur Belanja Pelanggan (Tanpa Login)

1.  **Start:** Pengunjung membuka website Night Stalkers.
    
2.  **Process:** Pengunjung menjelajahi katalog dan menambahkan produk ke **Keranjang Belanja (Session-based)**.
    
3.  **Process:** Pengunjung menekan tombol **Checkout**.
    
4.  **Input:** Pengunjung memasukkan data diri (Nama, Alamat, No. HP, Email) untuk pengiriman.
    
5.  **Validation:** Sistem memvalidasi kelengkapan data dan stok barang.
    
6.  **Decision:** \* **Jika Valid:** Data pesanan dan data diri pengunjung disimpan ke **Database MySQL** untuk evaluasi & proses kirim.
    
    *   **Jika Tidak:** Kembali ke form data diri.
        
7.  **Process:** Sesi belanja berakhir (data di browser akan hilang jika tab ditutup), namun data permanen sudah aman di database admin.
    
8.  **End.**
    

3\. Spesifikasi Teknis (Tech Stack)
-----------------------------------

### A. Perangkat Lunak (Software)

*   **Frontend:**
    
    *   **Framework:** Vue.js 3 (Vite).
        
    *   **State Management:** Pinia atau Vuex (untuk menyimpan data keranjang selama sesi aktif).
        
    *   **Styling:** TailwindCSS (Tema: Night Stalkers - Neon Green & Dark).
        
*   **Backend:**
    
    *   **Bahasa:** Go (Golang) 1.22+.
        
    *   **Framework:** Gin Web Framework.
        
    *   **Session Management:** Menggunakan UUID atau Cookie-based session untuk melacak keranjang pengunjung anonim.
        
*   **Database:**
    
    *   **Engine:** MySQL.
        
    *   **Struktur:** Tabel Orders menyimpan detail pembeli meskipun mereka tidak memiliki ID User tetap.
        

### B. Perangkat Keras Pengembangan (Hardware)

*   **Processor:** Intel(R) Core(TM) i5-10300H CPU @ 2.50GHz.
    
*   **RAM:** 16.0 GB.
    
*   **OS:** Windows 11 / Linux.
    

4\. Fitur Utama & Evaluasi Data
-------------------------------

**FiturLogika KerjaManfaat EvaluasiGuest Checkout**
User mengisi form identitas hanya saat akan membayar. Mendapatkan data profil pembeli tanpa memaksa registrasi.

**Session Cart**Keranjang belanja disimpan di memori browser/sesi backend. Menganalisis produk apa yang sering dimasukkan keranjang (meskipun tidak jadi beli).**Data Persistence

**Setelah checkout, data diri masuk ke tabel Guest\_Orders.Admin tetap bisa melakukan _follow-up_ atau analisis tren penjualan.

**Admin Dashboard**Panel kendali untuk melihat semua transaksi masuk.Memantau perputaran stok secara _real-time_.
