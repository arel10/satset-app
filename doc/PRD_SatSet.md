# PRD — SatSet
## Aplikasi Manajemen Laundry Berbasis Web dengan AI Assistant

| Atribut | Nilai |
|---|---|
| Versi | 1.1 |
| Tanggal | 14 Oktober 2026 |
| Status | Draft — Revisi |
| Nama Aplikasi | SatSet |
| Jenis Produk | Aplikasi internal untuk satu bisnis laundry |
| Platform | Web responsif |
| Bahasa UI | Bahasa Indonesia |
| Backend | Go + Gin |
| ORM | GORM |
| Frontend | Next.js + React + TypeScript |
| Database | MySQL 8 |
| AI Provider | Google Gemini API |
| Target Pengguna | Owner, Kasir, Operator |

---

## 1. Ringkasan Produk

**SatSet** adalah aplikasi web untuk membantu pemilik usaha laundry mengelola seluruh aktivitas operasional dalam satu tempat, mulai dari pencatatan pelanggan, pembuatan pesanan, pengelolaan layanan dan harga, pembayaran, proses pengerjaan, pengeluaran, hingga laporan bisnis.

SatSet dilengkapi AI Assistant berbasis Google Gemini yang memungkinkan Owner berinteraksi dengan data bisnis menggunakan bahasa sehari-hari.

Contoh pertanyaan:

- "Berapa omzet laundry saya bulan ini?"
- "Berapa banyak pesanan yang belum selesai?"
- "Siapa saja pelanggan yang belum mengambil cucian?"
- "Layanan apa yang paling banyak dipesan minggu ini?"
- "Berapa total pengeluaran untuk listrik dan deterjen bulan ini?"
- "Analisis kondisi bisnis laundry saya dan berikan rekomendasi."
- "Buatkan strategi promosi berdasarkan layanan yang paling laku."

AI Assistant dirancang untuk memberikan jawaban berdasarkan data aktual dari database, bukan sekadar menghasilkan jawaban berdasarkan pengetahuan umum model.

### Prinsip utama

1. Aplikasi tetap berfungsi normal meskipun layanan AI tidak tersedia.
2. Semua perhitungan transaksi dan laporan dilakukan oleh backend.
3. Gemini hanya dapat mengakses data melalui tools backend yang telah ditentukan.
4. AI Assistant hanya tersedia untuk Owner.
5. Semua transaksi dan perubahan penting memiliki pencatatan yang jelas.
6. Sistem dibuat untuk satu bisnis laundry, bukan SaaS multi-tenant.
7. Aplikasi dirancang agar mudah dikembangkan seiring pertumbuhan bisnis.

---

## 2. Tujuan Produk

1. Menggantikan pencatatan manual menggunakan buku, spreadsheet, atau chat.
2. Mempercepat proses pembuatan pesanan oleh Kasir.
3. Memudahkan pemantauan status pengerjaan setiap pesanan.
4. Mencegah kesalahan perhitungan harga, pembayaran, dan laporan.
5. Membantu Owner memahami kondisi bisnis melalui dashboard dan AI.
6. Menyediakan riwayat pelanggan dan transaksi yang mudah ditelusuri.
7. Menyederhanakan pekerjaan Operator melalui daftar pekerjaan terpusat.
8. Menghasilkan sistem operasional yang dapat digunakan dalam jangka panjang.

### 2.1 Indikator keberhasilan

| Area | Target |
|---|---|
| Pembuatan pesanan | Kurang dari 1 menit untuk transaksi normal |
| Pencarian | Pesanan dapat ditemukan berdasarkan nomor nota, nama, atau telepon |
| Akurasi transaksi | Total dihitung dan divalidasi oleh backend |
| Pembayaran | Nilai pembayaran konsisten dengan total tagihan |
| Laporan | Angka laporan sesuai dengan transaksi di database |
| AI Assistant | Jawaban angka bersumber dari hasil tools backend |
| Keamanan | Pengguna hanya dapat mengakses fitur sesuai perannya |
| Keandalan | Error AI tidak mengganggu transaksi utama |
| Responsivitas | Bisa digunakan melalui laptop, tablet, dan HP |

---

## 3. Pengguna dan Hak Akses

SatSet memiliki tiga peran pengguna dengan hak akses yang telah ditentukan.

### 3.1 Owner

Owner merupakan pemilik bisnis dan memiliki akses penuh terhadap seluruh fitur.

Hak akses:

- Melihat dashboard bisnis secara lengkap.
- Mengakses AI Assistant.
- Membuat dan melihat pesanan.
- Mengelola pelanggan.
- Mengelola layanan dan harga.
- Melihat riwayat pembayaran.
- Melihat, mengoreksi, atau membatalkan pembayaran sesuai aturan.
- Mencatat dan mengelola pengeluaran.
- Melihat laporan bisnis.
- Mengelola akun pengguna.
- Mengubah pengaturan bisnis.
- Membatalkan pesanan dengan alasan yang tercatat.
- Melihat audit log dan riwayat aktivitas penting.

### 3.2 Kasir

Kasir bertanggung jawab terhadap transaksi pelanggan.

Hak akses:

- Melihat dashboard operasional sederhana.
- Membuat pesanan baru.
- Mencari dan menambahkan pelanggan.
- Melihat daftar serta detail pesanan.
- Mencatat pembayaran penuh.
- Mencetak nota.
- Memantau status pesanan.
- Melakukan perubahan status yang diizinkan.
- Membantu proses pengambilan pesanan.

Batasan:

- Tidak dapat menggunakan AI Assistant.
- Tidak dapat mengakses pengeluaran dan laporan keuangan khusus Owner.
- Tidak dapat mengelola pengguna atau pengaturan bisnis.
- Tidak dapat mengoreksi atau membatalkan pembayaran.
- Tidak dapat menghapus data transaksi secara permanen.

### 3.3 Operator

Operator bertanggung jawab terhadap proses pengerjaan laundry.

Hak akses:

- Melihat halaman Daftar Pekerjaan.
- Melihat detail pekerjaan yang diperlukan.
- Melihat jenis layanan dan jumlah atau berat cucian.
- Melihat catatan khusus pengerjaan.
- Mengubah status pengerjaan sesuai alur yang diizinkan.

Batasan:

- Tidak dapat mencatat pembayaran.
- Tidak dapat mengelola pelanggan atau harga.
- Tidak dapat melihat pengeluaran dan laporan keuangan.
- Tidak dapat mengakses AI Assistant.
- Tidak dapat mengubah pengaturan bisnis.

**Ketentuan keamanan:** seluruh hak akses harus diperiksa di backend pada setiap endpoint. Menyembunyikan menu di frontend saja tidak dianggap sebagai kontrol keamanan.

---

## 4. Ruang Lingkup Produk

### 4.1 Fitur MVP

Fitur yang wajib tersedia pada versi pertama:

1. Login dan logout.
2. Manajemen pengguna dan peran.
3. Pengaturan profil bisnis.
4. Manajemen layanan dan harga.
5. Manajemen pelanggan.
6. Pembuatan pesanan melalui POS.
7. Penomoran pesanan otomatis.
8. Pencetakan nota.
9. Pelacakan status pengerjaan.
10. Pencatatan pembayaran.
11. Pembatalan pesanan dan audit log.
12. Pencatatan pengeluaran.
13. Dashboard Owner dan dashboard operasional.
14. Laporan bisnis dasar.
15. AI Assistant berbasis Gemini.
16. Riwayat percakapan AI.
17. Validasi, keamanan, dan backup database.

### 4.2 Fitur pengembangan berikutnya

Fitur berikut tidak menjadi persyaratan MVP:

- Notifikasi WhatsApp otomatis.
- Inventaris dan stok perlengkapan laundry.
- Poin loyalitas dan voucher pelanggan.
- Ekspor laporan ke Excel, CSV, dan PDF.
- Foto cucian atau dokumentasi kondisi pakaian.
- Integrasi payment gateway dan verifikasi QRIS otomatis.
- Mode offline untuk Kasir.
- Prediksi permintaan dan analisis tren menggunakan data historis.
- Pengingat otomatis untuk pesanan yang belum diambil.

### 4.3 Di luar ruang lingkup

- Sistem multi-bisnis atau multi-tenant.
- Manajemen banyak cabang.
- Aplikasi native Android dan iOS.
- Sistem akuntansi dan perpajakan lengkap.
- AI yang mengubah atau menghapus data transaksi.
- Pengiriman pesan pelanggan secara otomatis oleh AI.
- Sistem pembayaran otomatis melalui AI.

---

## 5. Kebutuhan Fungsional

### 5.1 Login dan Manajemen Pengguna

Sistem menyediakan autentikasi untuk Owner, Kasir, dan Operator.

Persyaratan:

- Login menggunakan username atau email dan password.
- Password disimpan menggunakan bcrypt atau Argon2id.
- Sistem menyediakan logout dan pengelolaan sesi yang aman.
- Gunakan session cookie yang aman atau access token berumur pendek dengan refresh token yang dikelola secara aman.
- Cookie autentikasi menggunakan HttpOnly, Secure pada HTTPS, dan kebijakan SameSite yang sesuai.
- Owner dapat menambahkan, mengubah peran, menonaktifkan, dan mengaktifkan kembali pengguna.
- Sistem membatasi percobaan login yang gagal berulang.
- Pengguna yang dinonaktifkan tidak dapat menggunakan aplikasi.
- Sistem mencatat aktivitas autentikasi penting.

### 5.2 Pengaturan Bisnis

Owner dapat mengelola informasi bisnis untuk digunakan pada dashboard dan nota.

Data pengaturan:

- Nama usaha.
- Alamat usaha.
- Nomor telepon.
- Logo usaha.
- Teks footer nota.
- Pengaturan identitas dasar aplikasi.

Ketentuan:

- Mata uang menggunakan Rupiah (IDR).
- Zona waktu menggunakan Asia/Jakarta.
- Pengaturan bisnis disimpan di database.
- Informasi pada nota mengikuti pengaturan terbaru.
- Perubahan pengaturan tidak mengubah riwayat transaksi yang sudah tercatat.

### 5.3 Layanan dan Harga

Owner dapat menambah, mengubah, dan menonaktifkan layanan laundry.

Atribut layanan:

- Nama layanan.
- Kategori layanan, jika diperlukan.
- Satuan: kg, pcs, atau pasang.
- Harga per satuan.
- Estimasi waktu pengerjaan.
- Status aktif atau nonaktif.
- Deskripsi opsional.

Contoh layanan:

- Cuci Kiloan.
- Cuci + Setrika.
- Setrika Saja.
- Bedcover.
- Sepatu.
- Dry Clean.

Daftar di atas merupakan contoh awal, bukan daftar layanan final.

Ketentuan:

- Layanan aktif dapat dipilih saat membuat pesanan.
- Layanan nonaktif tidak dapat digunakan untuk pesanan baru.
- Layanan yang sudah digunakan dalam transaksi tetap tercatat pada riwayat pesanan.
- Harga dan nama layanan disimpan sebagai snapshot pada item pesanan.
- Perubahan harga tidak mengubah transaksi lama.
- Harga menggunakan bilangan bulat Rupiah.
- Satuan berat menggunakan nilai desimal dengan presisi yang sesuai.

### 5.4 Manajemen Pelanggan

Sistem menyimpan informasi pelanggan agar transaksi mudah dilakukan dan riwayatnya dapat dilihat.

Data pelanggan:

- Nama pelanggan.
- Nomor telepon.
- Alamat opsional.
- Catatan opsional.
- Tanggal dibuat.
- Status pelanggan jika diperlukan.

Persyaratan:

- Nomor telepon dinormalisasi sebelum disimpan.
- Nomor telepon harus unik setelah normalisasi.
- Kasir dapat mencari pelanggan saat membuat pesanan.
- Kasir dapat menambahkan pelanggan tanpa meninggalkan alur pembuatan pesanan.
- Detail pelanggan menampilkan riwayat pesanan.
- Pencarian pelanggan mendukung nama dan nomor telepon.
- Data pelanggan tidak dihapus permanen jika masih diperlukan oleh riwayat transaksi.
- Perubahan informasi pelanggan tidak mengubah snapshot data transaksi yang sudah tersimpan.

### 5.5 Pesanan dan POS

Halaman POS menjadi salah satu fitur utama SatSet.

Kasir harus dapat membuat pesanan melalui alur yang sederhana dan cepat.

Data pesanan:

- Nomor pesanan.
- Pelanggan.
- Tanggal dan waktu masuk.
- Estimasi selesai.
- Daftar item layanan.
- Berat atau jumlah setiap item.
- Harga satuan.
- Subtotal.
- Diskon jika diizinkan.
- Total tagihan.
- Catatan khusus.
- Status pengerjaan.
- Status pembayaran.
- Pengguna yang membuat pesanan.

#### 5.5.1 Nomor pesanan

Format nomor pesanan:

`LDR-YYMMDD-001`

Contoh:

`LDR-261014-001`

Ketentuan:

- Nomor pesanan dibuat otomatis oleh backend.
- Nomor pesanan harus unik.
- Urutan dimulai dari 001 setiap hari.
- Pembuatan pesanan bersamaan tidak boleh menghasilkan nomor yang sama.
- Kegagalan transaksi tidak boleh merusak urutan nomor berikutnya.
- Penanganan nomor pesanan dilakukan secara aman melalui transaksi database atau mekanisme penguncian yang sesuai.

#### 5.5.2 Perhitungan pesanan

Rumus dasar:

`Subtotal item = jumlah atau berat × harga satuan`

`Total pesanan = jumlah subtotal item − diskon`

Ketentuan:

- Backend menjadi sumber kebenaran untuk semua perhitungan.
- Frontend hanya menampilkan estimasi perhitungan.
- Backend memvalidasi harga, jumlah, diskon, dan total.
- Diskon bersifat opsional dan hanya diberikan sesuai kewenangan.
- Semua nilai uang disimpan sebagai bilangan bulat Rupiah, bukan FLOAT.
- Setiap pesanan dapat memiliki lebih dari satu item layanan.

#### 5.5.3 Daftar dan detail pesanan

Fitur pencarian dan filter:

- Nomor pesanan.
- Nama pelanggan.
- Nomor telepon.
- Status pengerjaan.
- Tanggal masuk.
- Status pembayaran.
- Pesanan terlambat.

Daftar pesanan harus menggunakan pagination agar tetap responsif saat data bertambah.

#### 5.5.4 Pembatalan pesanan

- Pesanan tidak boleh dihapus permanen.
- Pembatalan menggunakan status DIBATALKAN.
- Pembatalan harus menyimpan alasan, waktu, dan pengguna yang melakukan tindakan.
- Pembatalan dilakukan oleh Owner atau Kasir sesuai kebijakan yang ditentukan.
- Pesanan yang sudah selesai atau sudah memiliki pembayaran harus ditangani berdasarkan aturan pembatalan yang aman.
- Jika uang sudah diterima, pembatalan tidak otomatis berarti refund.
- Refund atau koreksi pembayaran harus melalui proses terpisah yang berwenang.

#### 5.5.5 Pencegahan transaksi ganda

- Tombol simpan dinonaktifkan selama permintaan berlangsung.
- Backend mendukung `Idempotency-Key`.
- Permintaan ulang dengan kunci yang sama tidak boleh membuat pesanan baru.
- Pembuatan pesanan, item, dan pembayaran harus konsisten dalam transaksi database.

### 5.6 Pencetakan Nota

Setiap pesanan dapat dicetak menjadi nota untuk pelanggan.

Informasi nota:

- Identitas usaha.
- Nomor pesanan.
- Tanggal transaksi.
- Nama pelanggan.
- Daftar layanan.
- Jumlah atau berat.
- Harga satuan.
- Subtotal dan diskon.
- Total tagihan.
- Status pembayaran.
- Estimasi selesai.
- Catatan khusus yang relevan.
- Teks footer usaha.

Ketentuan:

- Nota dapat dicetak melalui browser.
- Layout nota harus ringkas dan mudah dibaca.
- Format printer thermal 58 mm, 80 mm, atau A4 belum diputuskan.
- Implementasi awal harus fleksibel agar format dapat disesuaikan setelah perangkat printer dikonfirmasi.
- Nota dapat dicetak ulang tanpa membuat transaksi baru.
- Dukungan penyimpanan PDF dapat disediakan melalui fitur cetak browser pada MVP, dengan pengembangan lebih lanjut jika diperlukan.

### 5.7 Status Pengerjaan

Alur standar:

`DITERIMA → DICUCI → DIKERINGKAN → DISETRIKA → DIKEMAS → SIAP_DIAMBIL → SELESAI`

Pesanan juga dapat berstatus:

`DIBATALKAN`

Ketentuan:

- Perubahan status dilakukan melalui backend.
- Sistem menolak perpindahan status yang tidak diizinkan.
- Owner dapat menangani pengecualian tertentu sesuai aturan dan audit log.
- Tahapan yang tidak relevan dapat dilewati sesuai jenis layanan.
- Setiap perubahan status mencatat waktu dan pengguna.
- Sistem menghitung status terlambat berdasarkan estimasi selesai.
- Pesanan terlambat adalah pesanan yang melewati estimasi selesai dan belum mencapai status SIAP_DIAMBIL atau SELESAI.
- Pesanan yang dibatalkan tidak dihitung sebagai pesanan aktif.
- Operator melihat pekerjaan yang perlu diproses dan dapat mengurutkannya berdasarkan estimasi selesai.

### 5.8 Pembayaran

**Kebijakan pembayaran SatSet: satu pesanan wajib dibayar lunas. Tidak ada pembayaran DP atau cicilan.**

Ketentuan:

- Pembayaran dilakukan sebesar total tagihan.
- Metode pembayaran: Tunai, Transfer, QRIS, dan Lainnya.
- Transfer dan QRIS dicatat secara manual pada MVP.
- Tidak ada verifikasi otomatis dari bank atau penyedia pembayaran.
- Sistem memvalidasi jumlah pembayaran terhadap total tagihan.
- Pesanan baru hanya dianggap LUNAS jika pembayaran valid telah menutupi seluruh tagihan.
- Pesanan tidak boleh dianggap lunas hanya karena Kasir memilih metode pembayaran.
- Status pembayaran menggunakan BELUM_LUNAS dan LUNAS.
- Status SEBAGIAN tidak digunakan dalam alur normal MVP.
- Pesanan yang belum lunas dapat tetap dicatat sesuai kebutuhan operasional, tetapi penyelesaiannya harus mengikuti kebijakan bisnis.
- Koreksi atau pembatalan pembayaran hanya dapat dilakukan Owner dan wajib tercatat di audit log.
- Pembayaran tidak boleh dihapus permanen tanpa jejak.
- Jika transaksi dibatalkan setelah pembayaran, sistem harus mempertahankan riwayat pembayaran dan mencatat tindak lanjutnya.

### 5.9 Pengeluaran

Fitur pengeluaran hanya dapat diakses Owner.

Data pengeluaran:

- Tanggal pengeluaran.
- Kategori.
- Nominal.
- Keterangan.
- Pengguna yang mencatat.
- Waktu pencatatan.

Contoh kategori:

- Deterjen.
- Pewangi.
- Listrik.
- Air.
- Gaji.
- Sewa.
- Perawatan mesin.
- Perlengkapan.
- Lainnya.

Ketentuan:

- Owner dapat menambah dan mengelola kategori.
- Nominal menggunakan bilangan bulat Rupiah.
- Pengeluaran dapat difilter berdasarkan tanggal dan kategori.
- Perubahan atau pembatalan catatan pengeluaran dicatat.
- Pengeluaran tidak boleh diakses Kasir dan Operator.
- Pengeluaran yang dibatalkan tidak boleh masuk perhitungan laporan.

### 5.10 Dashboard

#### Dashboard Owner

Menampilkan:

- Jumlah pesanan berdasarkan periode.
- Omzet pesanan.
- Kas yang diterima.
- Total pengeluaran.
- Estimasi laba operasional sederhana.
- Pesanan aktif.
- Pesanan terlambat.
- Pesanan siap diambil.
- Pesanan yang belum lunas.
- Layanan terlaris.
- Grafik omzet harian.
- Aktivitas operasional terbaru.
- AI Assistant.

#### Dashboard Kasir

Menampilkan:

- Ringkasan pesanan hari ini.
- Pesanan baru.
- Pesanan yang siap diambil.
- Pesanan yang belum lunas.
- Akses cepat ke POS dan daftar pesanan.

#### Dashboard Operator

Operator menggunakan halaman Daftar Pekerjaan sebagai halaman utama, dengan informasi pekerjaan yang relevan dan prioritas berdasarkan estimasi selesai.

#### Definisi angka dashboard

- **Omzet:** total nilai pesanan berdasarkan tanggal pesanan, tidak termasuk pesanan yang dibatalkan.
- **Kas diterima:** total pembayaran valid yang diterima berdasarkan tanggal pembayaran.
- **Pengeluaran:** total pengeluaran bisnis yang valid pada periode yang dipilih.
- **Estimasi laba:** kas diterima dikurangi pengeluaran pada periode yang sama. Ini merupakan indikator sederhana, bukan laba akuntansi.
- **Pesanan aktif:** pesanan yang belum selesai dan belum dibatalkan.
- **Pesanan terlambat:** pesanan melewati estimasi selesai dan belum siap diambil atau selesai.

Dashboard harus membedakan omzet, kas diterima, pengeluaran, dan estimasi laba secara jelas.

### 5.11 Laporan

Owner dapat memilih periode:

- Hari ini.
- Tujuh hari terakhir.
- Bulan ini.
- Rentang tanggal khusus.

Laporan MVP mencakup:

1. Ringkasan transaksi.
2. Omzet berdasarkan tanggal.
3. Kas diterima berdasarkan tanggal pembayaran.
4. Pengeluaran berdasarkan tanggal.
5. Pengeluaran per kategori.
6. Jumlah pesanan berdasarkan status.
7. Layanan paling banyak dipesan.
8. Pesanan yang belum lunas.
9. Pesanan terlambat.
10. Ringkasan performa bisnis.

Semua laporan harus berasal dari query database yang konsisten. Periode waktu dan dasar perhitungan harus dijelaskan dengan jelas.

---

## 6. AI Assistant — Gemini

### 6.1 Tujuan

AI Assistant membantu Owner memahami data operasional dan kondisi bisnis laundry melalui percakapan berbahasa Indonesia.

AI harus mampu menjawab pertanyaan tentang data aktual, memberikan ringkasan, membandingkan periode, dan menyusun rekomendasi berdasarkan data yang tersedia.

AI bukan pengganti laporan resmi. Angka transaksi tetap dihitung oleh backend.

### 6.2 Hak akses AI

- Hanya Owner yang dapat menggunakan AI Assistant.
- Kasir dan Operator tidak dapat membuka fitur AI.
- Backend wajib memeriksa autentikasi dan peran Owner pada setiap permintaan AI.
- Gemini tidak menerima akses langsung ke database.
- Gemini tidak menerima kredensial database atau API key internal.
- Semua akses data dilakukan melalui tools backend yang diizinkan.

### 6.3 Contoh pertanyaan

Pertanyaan operasional:

- "Berapa omzet hari ini?"
- "Bandingkan omzet minggu ini dengan minggu lalu."
- "Berapa pesanan yang sedang dicuci?"
- "Tampilkan pesanan yang terlambat."
- "Siapa pelanggan yang belum mengambil cucian?"
- "Berapa banyak pesanan yang selesai bulan ini?"

Pertanyaan keuangan:

- "Berapa total kas masuk bulan ini?"
- "Berapa pengeluaran untuk listrik dan air?"
- "Kategori pengeluaran mana yang paling besar?"
- "Berapa estimasi laba operasional bulan ini?"
- "Apakah omzet naik dibanding bulan sebelumnya?"

Pertanyaan pelanggan dan layanan:

- "Siapa pelanggan yang paling sering menggunakan layanan kita?"
- "Layanan apa yang paling laku minggu ini?"
- "Berapa rata-rata nilai pesanan?"
- "Layanan mana yang penjualannya menurun?"

Pertanyaan analisis:

- "Analisis kondisi bisnis laundry saya bulan ini."
- "Apa masalah operasional yang paling terlihat dari data?"
- "Berikan rekomendasi untuk meningkatkan omzet berdasarkan data yang tersedia."
- "Buatkan strategi promosi berdasarkan layanan terlaris."
- "Buatkan draft pesan WhatsApp untuk pelanggan yang belum mengambil cucian."

Pertanyaan yang memerlukan informasi di luar database boleh dijawab menggunakan pengetahuan umum Gemini, tetapi harus dibedakan dari fakta aktual bisnis.

### 6.4 Akses data Gemini

Gemini diperbolehkan memiliki akses luas terhadap informasi bisnis yang relevan, meliputi:

- Ringkasan omzet.
- Data pembayaran dan kas masuk.
- Daftar pesanan dan status pengerjaan.
- Informasi pelanggan yang diperlukan untuk menjawab pertanyaan.
- Riwayat transaksi pelanggan.
- Daftar layanan dan harga.
- Statistik layanan terlaris.
- Data pengeluaran.
- Ringkasan laporan bisnis.
- Metadata aktivitas operasional yang relevan.

Akses luas tersebut diberikan melalui kumpulan tools backend yang terkontrol, bukan akses bebas ke database.

Informasi pribadi seperti nama, nomor telepon, dan alamat hanya dikirim ke Gemini jika diperlukan untuk menjawab pertanyaan Owner. Jika jawaban cukup menggunakan agregasi atau identitas samaran, data pribadi tidak perlu dikirim.

### 6.5 Arsitektur AI

Alur kerja:

1. Owner memasukkan pertanyaan.
2. Backend memvalidasi sesi, peran, panjang pertanyaan, dan batas penggunaan.
3. Backend mengirim pertanyaan serta definisi tools yang tersedia ke Gemini.
4. Gemini memilih tools yang dibutuhkan.
5. Backend memvalidasi nama tools dan argumennya.
6. Backend menjalankan fungsi yang diizinkan menggunakan GORM dan query yang telah ditentukan.
7. Hasil dikembalikan ke Gemini.
8. Gemini menyusun jawaban bahasa Indonesia berdasarkan hasil tersebut.
9. Backend mengembalikan jawaban dan periode data yang digunakan kepada frontend.

Alur arsitektur:

`Frontend → Backend Go/Gin → Gemini API`

`Backend Go/Gin → AI Tools → Service/Repository → MySQL`

Gemini tidak boleh menghasilkan dan mengeksekusi SQL secara langsung.

### 6.6 Daftar AI Tools

| Tool | Fungsi |
|---|---|
| `get_revenue_summary` | Mengambil ringkasan omzet dan kas diterima |
| `get_order_summary` | Menghitung jumlah pesanan per status |
| `search_orders` | Mencari pesanan berdasarkan nomor, pelanggan, tanggal, atau status |
| `get_overdue_orders` | Mengambil daftar pesanan terlambat |
| `get_unpaid_orders` | Mengambil daftar tagihan yang belum lunas |
| `get_top_services` | Mengambil layanan paling banyak dipesan |
| `get_expense_summary` | Mengambil ringkasan pengeluaran |
| `get_customer_summary` | Mengambil statistik pelanggan |
| `get_customer_order_history` | Mengambil riwayat pesanan pelanggan |
| `get_payment_summary` | Mengambil ringkasan pembayaran berdasarkan periode |
| `get_business_settings` | Mengambil informasi bisnis yang relevan |
| `get_operational_summary` | Mengambil ringkasan kondisi operasional |
| `draft_customer_message` | Menyusun draft pesan pelanggan tanpa mengirimkannya |

Tools dapat ditambahkan seiring kebutuhan aplikasi, selama setiap tool memiliki tujuan jelas, validasi argumen, dan kontrol akses.

### 6.7 Aturan penggunaan tools

- Semua tool yang membaca database harus bersifat read-only.
- Model tidak boleh membuat atau menjalankan SQL.
- Backend hanya menerima nama tool yang terdaftar.
- Semua argumen divalidasi sebelum query dijalankan.
- Query harus menggunakan parameter yang aman.
- Hasil pencarian dibatasi, misalnya maksimal 50 baris.
- Permintaan data sensitif harus melalui pemeriksaan kewenangan.
- Jika pertanyaan membutuhkan data yang tidak tersedia, AI harus menyatakannya dengan jujur.
- Angka transaksi berasal dari hasil backend, bukan perkiraan model.
- Jawaban harus menyebutkan periode yang digunakan jika berkaitan dengan waktu.
- Backend menyimpan metadata tools yang dipanggil dan hasil sukses atau gagal.
- AI tidak boleh melakukan transaksi, mengubah harga, membatalkan pesanan, mengoreksi pembayaran, menghapus pelanggan, atau mengubah pengeluaran.

### 6.8 Analisis dan rekomendasi bisnis

Gemini boleh memberikan rekomendasi berdasarkan data yang tersedia, seperti:

- Layanan yang memiliki volume pesanan tertinggi.
- Perbandingan omzet antarperiode.
- Tren jumlah pesanan.
- Pola pengeluaran.
- Pelanggan yang sering menggunakan layanan.
- Peluang promosi.
- Saran perbaikan proses operasional.

Rekomendasi harus membedakan fakta dari interpretasi. AI tidak boleh mengklaim adanya tren atau penyebab tertentu jika data tidak cukup mendukung kesimpulan tersebut.

### 6.9 Keamanan dan privasi AI

- API key Gemini hanya disimpan di environment variable backend.
- Data pelanggan diminimalkan sesuai kebutuhan pertanyaan.
- Isi catatan pelanggan diperlakukan sebagai data, bukan instruksi untuk AI.
- Backend tidak mempercayai perintah yang muncul dalam isi database.
- Jumlah pemanggilan tools dibatasi, misalnya maksimal 5 kali per pertanyaan.
- Penggunaan AI dibatasi per pengguna, misalnya 30 pertanyaan per hari.
- Panjang pertanyaan dan ukuran hasil tools dibatasi.
- Request AI memiliki timeout, misalnya 30 detik.
- Data yang dikirim ke Gemini harus mengikuti kebijakan privasi dan ketentuan layanan Google yang berlaku.
- Sebelum penggunaan produksi, Owner harus mengevaluasi ketentuan penggunaan data pada tier Gemini yang dipilih.

### 6.10 Penanganan kegagalan AI

Jika Gemini mengalami gangguan, kuota habis, atau timeout:

- Tampilkan pesan bahwa AI Assistant sedang tidak tersedia.
- Sediakan akses kembali ke dashboard dan laporan biasa.
- Jangan membatalkan atau mengganggu transaksi yang sudah berhasil.
- Jangan menampilkan angka yang tidak didukung data.
- Catat metadata kegagalan untuk keperluan troubleshooting.

### 6.11 Riwayat percakapan

- Riwayat percakapan disimpan per pengguna.
- Owner dapat melihat percakapan miliknya sendiri.
- Owner dapat menghapus riwayat percakapan.
- Simpan waktu, pertanyaan, jawaban, metadata tools, dan status pemrosesan sesuai kebutuhan.
- Jangan menyimpan password, token autentikasi, atau API key dalam riwayat chat.
- Riwayat chat mengikuti kebijakan retensi data aplikasi.

---

## 7. Aturan Bisnis

1. Setiap pesanan mempunyai nomor unik.
2. Harga layanan disalin sebagai snapshot ketika pesanan dibuat.
3. Total pesanan dihitung dan divalidasi oleh backend.
4. Pembayaran normal wajib melunasi seluruh tagihan.
5. Pembayaran parsial atau cicilan tidak termasuk alur normal MVP.
6. Status pembayaran ditentukan dari tagihan dan pembayaran valid.
7. Pesanan tidak dihapus permanen.
8. Pembatalan pesanan harus memiliki catatan alasan dan pelaku.
9. Status pengerjaan mengikuti alur yang diizinkan.
10. Setiap perubahan status menyimpan waktu dan pengguna.
11. Semua waktu disimpan dalam UTC dan ditampilkan sesuai zona waktu Asia/Jakarta.
12. Uang disimpan sebagai bilangan bulat Rupiah, bukan FLOAT.
13. Berat atau jumlah desimal menggunakan tipe data dengan presisi yang sesuai.
14. Diskon, pembatalan, koreksi pembayaran, dan perubahan penting masuk audit log.
15. Omzet dihitung berdasarkan tanggal pesanan.
16. Kas diterima dihitung berdasarkan tanggal pembayaran.
17. Pengeluaran dihitung berdasarkan tanggal pengeluaran.
18. Estimasi laba dashboard merupakan indikator sederhana, bukan laporan laba akuntansi.
19. AI hanya membaca data melalui tools backend yang telah ditentukan.
20. AI tidak memiliki hak untuk mengubah data bisnis.
21. Akses AI hanya diberikan kepada Owner.
22. Data pribadi pelanggan hanya dikirim ke AI sesuai kebutuhan.
23. Semua tindakan penting harus dapat ditelusuri melalui audit log.

---

## 8. Arsitektur dan Teknologi

| Lapisan | Teknologi |
|---|---|
| Frontend | Next.js, React, TypeScript |
| Styling | Tailwind CSS |
| UI Components | shadcn/ui |
| Backend | Go + Gin |
| ORM | GORM |
| Database | MySQL 8 / InnoDB |
| Migrasi Database | golang-migrate atau Goose |
| AI | Google Gemini API |
| AI Integration | Function Calling dengan tools backend |
| Autentikasi | Session cookie atau JWT yang aman |
| Deployment | Docker + Docker Compose |
| Reverse Proxy | Nginx atau Caddy |
| CI | GitHub Actions, opsional |
| Version Control | Git |

### 8.1 Struktur backend

Backend menggunakan arsitektur monolit modular dengan pemisahan tanggung jawab.

Struktur konseptual:

- `handler`: menangani HTTP request dan response.
- `service`: menjalankan aturan bisnis.
- `repository`: mengakses database melalui GORM.
- `model`: mendefinisikan struktur data.
- `middleware`: autentikasi, otorisasi, logging, dan rate limiting.
- `routes`: mendaftarkan endpoint.
- `config`: membaca konfigurasi environment.
- `ai`: integrasi Gemini dan pengelolaan tools.
- `audit`: pencatatan aktivitas penting.
- `database`: koneksi dan konfigurasi database.

Contoh struktur direktori:

```text
backend/
├── cmd/
│   └── api/
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── model/
│   ├── middleware/
│   ├── routes/
│   ├── config/
│   ├── ai/
│   │   ├── provider/
│   │   ├── tools/
│   │   └── prompts/
│   └── audit/
├── migrations/
├── tests/
├── go.mod
├── go.sum
└── .env.example
```

Struktur tersebut merupakan acuan awal dan dapat disesuaikan selama pemisahan tanggung jawab tetap jelas.

### 8.2 Struktur database utama

| Tabel | Fungsi |
|---|---|
| `users` | Akun dan peran pengguna |
| `settings` | Pengaturan bisnis |
| `customers` | Data pelanggan |
| `services` | Layanan dan harga |
| `orders` | Data pesanan |
| `order_items` | Item dan snapshot harga |
| `order_status_history` | Riwayat status |
| `order_sequences` | Penomoran harian |
| `payments` | Catatan pembayaran |
| `expenses` | Data pengeluaran |
| `expense_categories` | Kategori pengeluaran |
| `ai_conversations` | Percakapan AI |
| `ai_messages` | Pesan dan jawaban AI |
| `audit_logs` | Riwayat tindakan penting |
| `idempotency_keys` | Pencegahan request ganda |

Struktur database tidak menggunakan `tenant_id` atau `branch_id` karena aplikasi hanya melayani satu bisnis.

### 8.3 Environment variables

Contoh `.env.example`:

```env
APP_NAME=SatSet
APP_ENV=development
APP_PORT=8080
APP_TIMEZONE=Asia/Jakarta

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=satset
DB_DSN=

SESSION_SECRET=replace_with_secure_random_secret
JWT_ACCESS_SECRET=replace_with_secure_random_secret
JWT_REFRESH_SECRET=replace_with_secure_random_secret

GEMINI_API_KEY=
GEMINI_MODEL=

AI_MAX_TOOL_CALLS=5
AI_DAILY_LIMIT_PER_USER=30
AI_REQUEST_TIMEOUT_SECONDS=30

FRONTEND_URL=http://localhost:3000
```

Catatan:

- Konfigurasi database menyesuaikan lingkungan lokal dan produksi.
- Jika menggunakan XAMPP, database lokal dapat memakai MySQL/MariaDB yang disediakan XAMPP, dengan penyesuaian versi dan kompatibilitas.
- Pilih satu strategi autentikasi yang konsisten; variabel session dan JWT tidak harus digunakan bersamaan.
- `DB_DSN` sebaiknya dibentuk dari konfigurasi database dan tidak perlu berisi kredensial yang di-commit.
- `GEMINI_MODEL` harus menggunakan model yang tersedia pada akun dan API yang dipakai.
- Secret produksi harus dibuat secara acak dan kuat.
- File `.env` tidak boleh di-commit ke Git.
- Kredensial produksi dikelola secara terpisah dan aman.

### 8.4 Prinsip pengembangan

- Tidak menggunakan microservices pada MVP.
- Tidak membutuhkan Redis untuk tahap awal.
- Tidak membutuhkan Kubernetes.
- Tidak membutuhkan message broker.
- Hindari abstraksi berlebihan yang belum diperlukan.
- Semua perhitungan bisnis utama dilakukan di backend.
- Database menjadi sumber data utama.
- Integrasi AI dipisahkan agar provider dapat diganti di masa depan.
- Query laporan harus mempertimbangkan indeks dan konsistensi data.

---

## 9. Kebutuhan Nonfungsional

### 9.1 Keamanan

- HTTPS pada lingkungan produksi.
- Password di-hash.
- Secret disimpan di environment variable.
- Validasi input dilakukan di backend.
- Query database menggunakan parameter yang aman.
- Hak akses diperiksa di setiap endpoint.
- Login dan AI menggunakan rate limiting.
- Cegah akses lintas peran.
- Jangan mencatat password, token, atau API key ke log.
- Gunakan transaksi database untuk operasi yang saling berkaitan.
- Catat perubahan transaksi penting pada audit log.

### 9.2 Performa

- Target awal respons API biasa sekitar 500 ms pada persentil ke-95 untuk beban bisnis yang wajar.
- Gunakan indeks pada nomor pesanan, telepon pelanggan, status, dan tanggal.
- Daftar pesanan menggunakan pagination.
- Laporan memakai query yang efisien.
- AI menggunakan streaming jika didukung integrasi yang dipilih.
- Operasi AI tidak boleh memblokir operasi kasir secara keseluruhan.

### 9.3 Keandalan dan backup

- Pembuatan pesanan dan item dilakukan dalam satu transaksi database.
- Pencatatan pembayaran harus konsisten dengan tagihan.
- Sediakan endpoint kesehatan seperti `/healthz`.
- Log aplikasi diarahkan ke stdout.
- Backup database dilakukan secara rutin.
- Backup harus disimpan terpisah dari database utama.
- Uji pemulihan backup sebelum aplikasi dipakai secara nyata.
- Kegagalan Gemini tidak boleh merusak data atau menghentikan fitur utama.

### 9.4 Usability

- Antarmuka menggunakan Bahasa Indonesia.
- Format uang: `Rp 1.234.567`.
- Halaman POS dioptimalkan untuk input cepat.
- Tombol penting mudah dijangkau di layar sentuh.
- Status pesanan menggunakan warna dan label teks.
- Setiap halaman data memiliki loading, empty state, dan error state.
- Konfirmasi diberikan untuk tindakan penting.
- Formulir meminimalkan input berulang.

### 9.5 Privasi

- Data pelanggan dikumpulkan sesuai kebutuhan operasional.
- Data pribadi yang dikirim ke Gemini diminimalkan.
- Akses data pelanggan dibatasi berdasarkan peran.
- Riwayat AI mengikuti kebijakan retensi yang ditentukan.
- Pengelolaan data mengikuti ketentuan perlindungan data pribadi yang berlaku di Indonesia.

---

## 10. Halaman Aplikasi

| Peran | Halaman |
|---|---|
| Owner | Dashboard |
| Owner | AI Assistant |
| Owner | Buat Pesanan |
| Owner | Daftar Pesanan |
| Owner | Detail Pesanan |
| Owner | Pelanggan |
| Owner | Layanan & Harga |
| Owner | Pengeluaran |
| Owner | Laporan |
| Owner | Manajemen Pengguna |
| Owner | Pengaturan Bisnis |
| Owner | Audit Log |
| Kasir | Dashboard Operasional |
| Kasir | Buat Pesanan |
| Kasir | Daftar Pesanan |
| Kasir | Detail Pesanan |
| Kasir | Pelanggan |
| Operator | Daftar Pekerjaan |
| Operator | Detail Pekerjaan |
| Semua peran | Login |

Navigasi, menu, dan halaman harus disesuaikan dengan hak akses pengguna yang sedang login.

---

## 11. Pengujian

### 11.1 Unit testing

Pengujian untuk:

- Perhitungan subtotal dan total.
- Validasi diskon.
- Perhitungan status pembayaran.
- Validasi pembayaran lunas.
- Transisi status pengerjaan.
- Perhitungan pesanan terlambat.
- Perhitungan omzet dan kas diterima.
- Perhitungan pengeluaran dan estimasi laba.
- Validasi hak akses.

### 11.2 Integration testing

- Login dan logout.
- CRUD pelanggan.
- CRUD layanan.
- Pembuatan pesanan.
- Penomoran pesanan otomatis.
- Pencatatan pembayaran penuh.
- Pembatalan pesanan.
- Perubahan status pengerjaan.
- Pengeluaran dan laporan.
- Pemeriksaan hak akses per endpoint.
- Integrasi AI tools dengan database.

### 11.3 Pengujian konkurensi

- Beberapa pesanan dibuat bersamaan tanpa nomor ganda.
- Request berulang tidak menghasilkan transaksi ganda.
- Kegagalan transaksi tidak menghasilkan data pesanan setengah jadi.
- Perubahan status dan pembayaran tetap konsisten.

### 11.4 Pengujian AI

- Pertanyaan bisnis normal.
- Perbandingan periode.
- Data kosong.
- Pertanyaan yang membutuhkan data pelanggan.
- Pertanyaan di luar kewenangan.
- Upaya prompt injection melalui pertanyaan atau catatan pelanggan.
- Percobaan memanggil tools yang tidak tersedia.
- Gemini mengalami timeout atau kuota habis.
- Jawaban angka harus sesuai hasil tools.
- AI tidak dapat mengubah atau menghapus data.
- Pengguna non-Owner tidak dapat mengakses AI.
- Batas pemanggilan tools dan kuota harian diterapkan.

### 11.5 User Acceptance Testing

Owner, Kasir, dan Operator mencoba skenario operasional nyata untuk memastikan aplikasi sesuai dengan proses laundry sehari-hari.

---

## 12. Rencana Pengerjaan

| Tahap | Fokus | Hasil |
|---|---|---|
| 1 | Setup repository, konfigurasi, database, migrasi, login, dan peran | Aplikasi dapat dijalankan dan pengguna dapat login |
| 2 | Pengaturan bisnis, layanan, pelanggan, dan POS | Kasir dapat membuat pesanan |
| 3 | Pembayaran, nota, status pengerjaan, dan audit log | Transaksi dapat diproses hingga selesai |
| 4 | Pengeluaran, dashboard, dan laporan | Owner dapat memantau bisnis |
| 5 | Gemini API, tools read-only, dan UI AI Assistant | Owner dapat bertanya tentang bisnis |
| 6 | Testing, backup, keamanan, dan penyempurnaan UI | MVP siap diuji |
| 7 | Deployment dan uji coba operasional | Aplikasi siap digunakan secara nyata |

AI dikembangkan setelah fitur transaksi dan laporan stabil agar jawaban AI menggunakan data yang benar.

---

## 13. Kriteria Selesai MVP

MVP dianggap selesai jika:

- Owner, Kasir, dan Operator dapat login.
- Setiap peran hanya dapat mengakses fitur yang diizinkan.
- Owner dapat mengelola pengguna, layanan, dan pengaturan bisnis.
- Kasir dapat membuat pesanan dengan cepat.
- Nomor pesanan otomatis dan unik.
- Harga lama tetap tersimpan pada transaksi.
- Kasir dapat mencatat pembayaran penuh.
- Nota dapat dicetak.
- Status pengerjaan dapat dilacak sampai selesai.
- Pembatalan transaksi meninggalkan riwayat.
- Owner dapat mencatat pengeluaran.
- Dashboard dan laporan konsisten dengan data database.
- AI Assistant hanya dapat digunakan Owner.
- Gemini dapat mengakses data yang relevan melalui tools backend yang diizinkan.
- AI memberikan jawaban berdasarkan hasil data aktual.
- AI tidak dapat mengubah data bisnis.
- Gangguan AI tidak mengganggu fitur utama.
- Backup database berhasil dibuat dan dapat dipulihkan.
- Aplikasi dapat dijalankan pada lingkungan deployment yang ditentukan.

---

## 14. Keputusan dan Hal yang Belum Ditentukan

| No. | Hal | Keputusan |
|---|---|---|
| 1 | Daftar layanan dan harga awal | Belum ditentukan; dapat diatur melalui aplikasi |
| 2 | Akses AI Assistant | Hanya Owner |
| 3 | Jenis printer nota | Belum ditentukan; desain harus fleksibel |
| 4 | Alur status pengerjaan | Disetujui untuk MVP |
| 5 | Metode pembayaran | Wajib lunas; tidak ada DP atau cicilan |
| 6 | Hosting produksi | Belum ditentukan |
| 7 | Framework backend | Go + Gin |
| 8 | ORM | GORM |
| 9 | Database | MySQL 8 |
| 10 | Akses Gemini | Luas terhadap data bisnis yang relevan melalui tools read-only |
| 11 | Pengiriman data pelanggan ke Gemini | Diizinkan jika diperlukan, dengan pembatasan data pribadi |
| 12 | Model Gemini | Ditentukan saat implementasi berdasarkan model yang tersedia |

Hal-hal yang belum ditentukan tidak boleh menghambat pengembangan fitur inti. Nilai yang memang bersifat konfigurasi harus dibuat mudah diubah tanpa perlu memodifikasi banyak bagian kode.

---

## 15. Arahan Implementasi

Pengembangan SatSet harus mengutamakan aplikasi yang benar-benar siap digunakan untuk operasional laundry, bukan sekadar demonstrasi atau prototype.

Prioritas implementasi:

1. Akurasi data dan transaksi.
2. Keamanan autentikasi dan hak akses.
3. Kecepatan proses Kasir.
4. Kemudahan penggunaan oleh Owner dan Operator.
5. Konsistensi laporan.
6. Integrasi Gemini yang aman dan bermanfaat.
7. Tampilan profesional dan responsif.
8. Struktur kode yang mudah dipelihara dan dikembangkan.

**Keputusan akhir:** SatSet adalah aplikasi manajemen laundry untuk satu bisnis, dengan backend Go + Gin, ORM GORM, database MySQL, frontend Next.js, serta AI Assistant Gemini khusus Owner yang memiliki akses baca luas melalui tools backend yang terkontrol.