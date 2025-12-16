![nightstalkers_flowchart_svg](https://github.com/user-attachments/assets/1a318aba-5b87-4a4a-ae6e-b40bccccc9de)Dokumentasi Perencanaan Pembuatan Web Ecommerce Night Stalkers
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
    

![Uploading nightstal<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 2800" style="background: white;">
  <defs>
    <marker id="arrowhead" markerWidth="10" markerHeight="10" refX="9" refY="3" orient="auto">
      <polygon points="0 0, 10 3, 0 6" fill="#000" />
    </marker>
    <style>
      .box { fill: white; stroke: black; stroke-width: 3; }
      .diamond { fill: white; stroke: black; stroke-width: 3; }
      .text { font-family: Arial, sans-serif; font-size: 14px; fill: black; text-anchor: middle; }
      .line { stroke: black; stroke-width: 2; fill: none; marker-end: url(#arrowhead); }
      .label { font-family: Arial, sans-serif; font-size: 12px; fill: black; }
    </style>
  </defs>

  <!-- Start -->
  <rect x="450" y="20" width="300" height="60" rx="30" class="box"/>
  <text x="600" y="55" class="text">Start: Akses Night Stalkers</text>
  <line x1="600" y1="80" x2="600" y2="110" class="line"/>

  <!-- Jelajahi Katalog -->
  <rect x="450" y="110" width="300" height="60" class="box"/>
  <text x="600" y="145" class="text">Jelajahi Katalog Buku/Komik</text>
  <line x1="600" y1="170" x2="600" y2="200" class="line"/>

  <!-- Pilih Produk (Diamond) -->
  <polygon points="600,200 750,260 600,320 450,260" class="diamond"/>
  <text x="600" y="265" class="text">Pilih Produk</text>
  
  <!-- Tambah ke Keranjang -->
  <line x1="750" y1="260" x2="900" y2="260" class="line"/>
  <text x="810" y="250" class="label">Tambah</text>
  <rect x="900" y="230" width="250" height="60" class="box"/>
  <text x="1025" y="265" class="text">Simpan ke Keranjang Sesi</text>
  <line x1="1025" y1="230" x2="1025" y2="140" class="line"/>
  <line x1="1025" y1="140" x2="750" y2="140" class="line"/>
  <line x1="750" y1="140" x2="750" y2="170" class="line"/>

  <!-- Checkout -->
  <line x1="600" y1="320" x2="600" y2="360" class="line"/>
  <text x="650" y="345" class="label">Checkout</text>
  <rect x="450" y="360" width="300" height="60" class="box"/>
  <text x="600" y="395" class="text">Cek Stok Produk di Keranjang</text>
  <line x1="600" y1="420" x2="600" y2="450" class="line"/>

  <!-- Stok Tersedia? -->
  <polygon points="600,450 750,510 600,570 450,510" class="diamond"/>
  <text x="600" y="515" class="text">Stok Tersedia?</text>

  <!-- Tidak Tersedia -->
  <line x1="450" y1="510" x2="250" y2="510" class="line"/>
  <text x="370" y="500" class="label">Tidak</text>
  <rect x="50" y="480" width="200" height="60" class="box"/>
  <text x="150" y="510" class="text">Tampilkan Pesan</text>
  <text x="150" y="525" class="text">Error Stok</text>
  <line x1="150" y1="480" x2="150" y2="400" class="line"/>
  <rect x="50" y="340" width="200" height="60" class="box"/>
  <text x="150" y="365" class="text">Hapus Item</text>
  <text x="150" y="380" class="text">Tidak Tersedia</text>
  <line x1="150" y1="340" x2="150" y2="140" class="line"/>
  <line x1="150" y1="140" x2="450" y2="140" class="line"/>

  <!-- Ya - Reserve Stok -->
  <line x1="600" y1="570" x2="600" y2="610" class="line"/>
  <text x="650" y="595" class="label">Ya</text>
  <rect x="450" y="610" width="300" height="60" class="box"/>
  <text x="600" y="645" class="text">Reserve Stok Sementara</text>
  <line x1="600" y1="670" x2="600" y2="700" class="line"/>

  <!-- Form Data Diri -->
  <rect x="450" y="700" width="300" height="60" class="box"/>
  <text x="600" y="725" class="text">Form Data Diri &</text>
  <text x="600" y="745" class="text">Alamat Pengiriman</text>
  <line x1="600" y1="760" x2="600" y2="790" class="line"/>

  <!-- Validasi Data -->
  <polygon points="600,790 750,850 600,910 450,850" class="diamond"/>
  <text x="600" y="855" class="text">Validasi Data</text>

  <!-- Tidak Valid -->
  <line x1="750" y1="850" x2="900" y2="850" class="line"/>
  <text x="810" y="840" class="label">Tidak Valid</text>
  <rect x="900" y="820" width="200" height="60" class="box"/>
  <text x="1000" y="845" class="text">Tampilkan Error</text>
  <text x="1000" y="860" class="text">Validasi</text>
  <line x1="1000" y1="820" x2="1000" y2="730" class="line"/>
  <line x1="1000" y1="730" x2="750" y2="730" class="line"/>

  <!-- Valid -->
  <line x1="600" y1="910" x2="600" y2="950" class="line"/>
  <text x="650" y="935" class="label">Valid</text>
  <rect x="450" y="950" width="300" height="60" class="box"/>
  <text x="600" y="985" class="text">Pilih Metode Pembayaran</text>
  <line x1="600" y1="1010" x2="600" y2="1040" class="line"/>

  <!-- COD -->
  <rect x="450" y="1040" width="300" height="60" class="box"/>
  <text x="600" y="1075" class="text">COD - Cash on Delivery</text>
  <line x1="600" y1="1100" x2="600" y2="1130" class="line"/>

  <!-- Konfirmasi Pesanan -->
  <rect x="450" y="1130" width="300" height="60" class="box"/>
  <text x="600" y="1165" class="text">Konfirmasi Pesanan</text>
  <line x1="600" y1="1190" x2="600" y2="1220" class="line"/>

  <!-- Proses Pesanan -->
  <polygon points="600,1220 750,1280 600,1340 450,1280" class="diamond"/>
  <text x="600" y="1285" class="text">Proses Pesanan</text>

  <!-- Error -->
  <line x1="450" y1="1280" x2="280" y2="1280" class="line"/>
  <text x="380" y="1270" class="label">Error</text>
  <rect x="130" y="1250" width="150" height="60" class="box"/>
  <text x="205" y="1285" class="text">Retry Mechanism</text>
  <line x1="205" y1="1310" x2="205" y2="1340" class="line"/>

  <!-- Retry Berhasil? -->
  <polygon points="205,1340 330,1400 205,1460 80,1400" class="diamond"/>
  <text x="205" y="1400" class="text">Retry</text>
  <text x="205" y="1415" class="text">Berhasil?</text>

  <!-- Retry Tidak -->
  <line x1="80" y1="1400" x2="20" y2="1400" class="line"/>
  <line x1="20" y1="1400" x2="20" y2="1500" class="line"/>
  <text x="50" y="1390" class="label">Tidak</text>
  <rect x="20" y="1500" width="150" height="60" class="box"/>
  <text x="95" y="1535" class="text">Kembalikan Stok</text>
  <line x1="95" y1="1560" x2="95" y2="1590" class="line"/>
  
  <rect x="20" y="1590" width="200" height="60" class="box"/>
  <text x="120" y="1615" class="text">Tampilkan Error &</text>
  <text x="120" y="1635" class="text">Opsi Coba Lagi</text>
  <line x1="120" y1="1650" x2="120" y2="1680" class="line"/>

  <!-- User Action -->
  <polygon points="120,1680 245,1740 120,1800 -5,1740" class="diamond"/>
  <text x="120" y="1745" class="text">User Action</text>

  <!-- Coba Lagi -->
  <line x1="120" y1="1680" x2="120" y2="1160" class="line"/>
  <line x1="120" y1="1160" x2="450" y2="1160" class="line"/>
  <text x="200" y="1670" class="label">Coba Lagi</text>

  <!-- Batal -->
  <line x1="120" y1="1800" x2="120" y2="1840" class="line"/>
  <text x="70" y="1825" class="label">Batal</text>
  <rect x="20" y="1840" width="200" height="60" class="box"/>
  <text x="120" y="1870" class="text">Update Status: Dibatalkan</text>
  <line x1="120" y1="1900" x2="120" y2="1930" class="line"/>
  
  <rect x="20" y="1930" width="200" height="60" class="box"/>
  <text x="120" y="1965" class="text">Hapus Data Sesi</text>
  <line x1="120" y1="1990" x2="120" y2="2750" class="line"/>
  <line x1="120" y1="2750" x2="550" y2="2750" class="line"/>

  <!-- Retry Ya -->
  <line x1="330" y1="1400" x2="380" y2="1400" class="line"/>
  <text x="350" y="1390" class="label">Ya</text>
  <line x1="380" y1="1400" x2="380" y2="1520" class="line"/>
  <line x1="380" y1="1520" x2="600" y2="1520" class="line"/>

  <!-- Sukses -->
  <line x1="600" y1="1340" x2="600" y2="1520" class="line"/>
  <text x="650" y="1420" class="label">Sukses</text>

  <!-- Simpan Data Pesanan -->
  <rect x="450" y="1520" width="300" height="60" class="box"/>
  <text x="600" y="1545" class="text">Simpan Data Pesanan</text>
  <text x="600" y="1565" class="text">ke MySQL</text>
  <line x1="600" y1="1580" x2="600" y2="1610" class="line"/>

  <!-- Kurangi Stok -->
  <rect x="450" y="1610" width="300" height="60" class="box"/>
  <text x="600" y="1645" class="text">Kurangi Stok Definitif</text>
  <line x1="600" y1="1670" x2="600" y2="1700" class="line"/>

  <!-- Generate Nota -->
  <rect x="450" y="1700" width="300" height="60" class="box"/>
  <text x="600" y="1735" class="text">Generate Nota/Invoice</text>
  <line x1="600" y1="1760" x2="600" y2="1790" class="line"/>

  <!-- Kirim Notifikasi Admin -->
  <rect x="450" y="1790" width="300" height="60" class="box"/>
  <text x="600" y="1825" class="text">Kirim Notifikasi ke Admin</text>
  <line x1="600" y1="1850" x2="600" y2="1880" class="line"/>

  <!-- Tampilkan Status -->
  <rect x="425" y="1880" width="350" height="80" class="box"/>
  <text x="600" y="1910" class="text">Tampilkan Status Pembayaran:</text>
  <text x="600" y="1930" class="text">COD Menunggu Konfirmasi</text>
  <line x1="600" y1="1960" x2="600" y2="1990" class="line"/>

  <!-- Tampilkan Nota -->
  <rect x="450" y="1990" width="300" height="60" class="box"/>
  <text x="600" y="2025" class="text">Tampilkan Nota di Website</text>
  <line x1="600" y1="2050" x2="600" y2="2080" class="line"/>

  <!-- Hapus Sesi -->
  <rect x="450" y="2080" width="300" height="60" class="box"/>
  <text x="600" y="2115" class="text">Hapus Data Sesi Keranjang</text>
  <line x1="600" y1="2140" x2="600" y2="2170" class="line"/>

  <!-- End -->
  <rect x="450" y="2170" width="300" height="60" rx="30" class="box"/>
  <text x="600" y="2205" class="text">End</text>

</svg>kers_flowchart_svg.svg…]()


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
