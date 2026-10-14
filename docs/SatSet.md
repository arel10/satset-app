# Product Requirements Document (PRD)

# SatSet — All-in-One Laundry Management Platform

| Atribut | Nilai |
|---|---|
| **Versi** | 1.1 (revisi dari 1.0) |
| **Status** | Draft untuk pengembangan |
| **Tanggal** | 14 Oktober 2026 |
| **Nama produk** | SatSet (nama kerja) |
| **Kategori** | Laundry Management System, Business Intelligence, AI Assistant |
| **Platform** | Web responsif (PWA-ready) |
| **Bahasa utama** | Bahasa Indonesia |
| **Frontend** | React / Next.js (TypeScript) |
| **Backend** | Go (Golang) dengan framework HTTP (Gin / Fiber / Echo) |
| **Database** | MySQL 8.x |
| **LLM** | Google Gemini API via Google AI Studio (free tier) |
| **Infrastruktur** | Docker, Docker Compose, CI/CD |

### Ringkasan Perubahan dari Versi 1.0

1. Database diganti dari PostgreSQL ke **MySQL 8**; desain data disesuaikan (tipe uang, UTC, indeks, locking, nomor pesanan).
2. Backend dipastikan memakai **framework Go** (rekomendasi: Gin atau Fiber) dengan struktur clean/modular architecture.
3. Ditambahkan bab **CI/CD dan Docker** (pipeline, environment, strategi rilis).
4. Ditambahkan bab **Peningkatan Performa dan Skalabilitas** (Redis, job queue, reverse proxy, CDN, observability, load testing).
5. Bagian **AI Assistant** disesuaikan dengan keterbatasan **Gemini free tier** (kuota, rate limit, privasi data, caching, fallback).
6. Ditambahkan: state machine status pesanan, idempotensi transaksi, penanganan uang, konvensi API, strategi pengujian, dan matriks keputusan teknologi.
7. Struktur dokumen dirapikan agar lebih mudah dijadikan acuan pengembangan.

---

## 1. Ringkasan Produk

SatSet adalah platform manajemen bisnis laundry yang mengintegrasikan operasional harian dan asisten bisnis berbasis Large Language Model (LLM) dalam satu aplikasi.

Aplikasi membantu pemilik dan karyawan mengelola pesanan, pelanggan, layanan, pembayaran, proses produksi, pengeluaran, karyawan, serta laporan bisnis. AI Business Assistant memungkinkan pengguna bertanya dengan bahasa sehari-hari, memperoleh jawaban berdasarkan data aktual yang diizinkan, dan meminta bantuan menyusun draft tindakan.

Produk dirancang untuk laundry kiloan, satuan, sepatu, dan dry cleaning, dengan satu atau beberapa cabang. Pengembangan dimulai dari MVP untuk satu bisnis, namun arsitektur disiapkan agar dapat berkembang menjadi SaaS multi-tenant.

### 1.1 Visi

Menjadi platform operasional dan asisten bisnis yang membantu pemilik laundry menjalankan usaha secara lebih efisien, terukur, dan berbasis data.

### 1.2 Permasalahan yang Ingin Diselesaikan

1. Pencatatan pesanan dan pembayaran masih manual atau tersebar di beberapa aplikasi.
2. Pemilik sulit mengetahui status pesanan, pekerjaan terlambat, dan beban kerja.
3. Data transaksi tidak mudah diubah menjadi informasi bisnis yang berguna.
4. Pencatatan biaya dan stok tidak terintegrasi dengan transaksi.
5. Komunikasi ke pelanggan (misalnya pesanan selesai) memerlukan pekerjaan berulang.
6. Pemilik butuh cara sederhana menanyakan kondisi bisnis tanpa membuat laporan manual.
7. Sistem perlu fondasi untuk penambahan cabang dan pelanggan bisnis baru.

### 1.3 Solusi

- POS dan manajemen pesanan.
- Alur produksi laundry dengan state machine status.
- Manajemen pelanggan dan layanan.
- Pembayaran, pengeluaran, dan laporan keuangan operasional.
- Inventaris perlengkapan dan pemeliharaan peralatan (pasca-MVP).
- Manajemen karyawan dan hak akses (RBAC).
- AI Business Assistant berbasis tool/function calling melalui backend.
- Fondasi multi-cabang dan multi-tenant.

---

## 2. Tujuan Produk dan Metrik Keberhasilan

### 2.1 Tujuan

1. Memusatkan proses bisnis laundry dalam satu aplikasi.
2. Mengurangi kesalahan pencatatan transaksi dan status pesanan.
3. Mempermudah pemantauan pesanan, pembayaran, dan pekerjaan terlambat.
4. Menyediakan laporan operasional dan keuangan yang konsisten.
5. Memberi insight dari data melalui AI.
6. Menjaga keamanan data, transaksi, dan tindakan yang dijalankan AI.
7. Menyiapkan arsitektur yang dapat diperluas tanpa membangun ulang.
8. Memastikan aplikasi cepat dan stabil bahkan pada perangkat dan koneksi kelas menengah.

### 2.2 Metrik Keberhasilan MVP

Target berikut adalah sasaran awal yang divalidasi lewat pengujian pengguna, bukan klaim hasil.

| Area | Target |
|---|---|
| Usability | Kasir baru dapat membuat transaksi pertama tanpa bantuan teknis |
| Integritas | Setiap pesanan memiliki nomor unik dan riwayat status |
| Pencarian | Pesanan ditemukan lewat nomor nota, nama, atau telepon |
| Konsistensi laporan | Total laporan sama dengan agregasi transaksi sah (selisih 0) |
| AI akurasi | Semua angka pada jawaban AI berasal dari hasil tool; tanpa angka karangan |
| AI keamanan | Tidak ada aksi sensitif tanpa otorisasi dan konfirmasi |
| Isolasi | Tidak ada akses lintas tenant/cabang di luar hak pengguna |
| Performa | Sesuai target di Bab 10 (p95 API, LCP, dsb.) |
| Ketersediaan | Kegagalan AI tidak mengganggu fitur kasir/transaksi |

---

## 3. Target Pengguna

| Peran | Deskripsi |
|---|---|
| **Owner** | Pemilik usaha; kontrol penuh atas operasional, pendapatan, pengeluaran, layanan, pelanggan, karyawan, dan insight |
| **Manager** | Pengelola cabang; memantau pesanan, pekerjaan karyawan, stok, dan laporan sesuai izin |
| **Kasir** | Membuat transaksi, mengelola pelanggan, menerima pembayaran, mencetak/membagikan nota, memperbarui pesanan sesuai izin |
| **Operator Laundry** | Menjalankan pencucian, pengeringan, penyetrikaan, pengemasan, dan kontrol kualitas |
| **Super Admin Platform** | (Tahap SaaS) Mengelola tenant, paket, konfigurasi platform, dukungan. Tidak boleh mengakses data bisnis tenant secara bebas tanpa kebutuhan sah, izin, dan audit |

---

## 4. Ruang Lingkup Produk

### 4.1 MVP

1. Autentikasi dan otorisasi.
2. Profil bisnis dan cabang awal.
3. Manajemen layanan dan harga.
4. Manajemen pelanggan.
5. POS dan pembuatan pesanan.
6. Pelacakan status pesanan.
7. Pencatatan pembayaran.
8. Dashboard dan laporan dasar.
9. Pencatatan pengeluaran sederhana.
10. AI Business Assistant (read-only, tool terbatas).
11. Audit log tindakan penting.
12. Pengaturan dasar, keamanan, dan fondasi DevOps (Docker + CI/CD).

### 4.2 Pasca-MVP

- Inventaris lengkap dan pemasok.
- Loyalitas, voucher, retensi pelanggan.
- Integrasi WhatsApp Business Platform (API resmi/penyedia sah).
- QRIS/payment gateway dan rekonsiliasi.
- Pemeliharaan mesin.
- Penjadwalan karyawan dan analisis produktivitas.
- Analisis pelanggan berisiko churn.
- Multi-cabang lanjutan.
- SaaS multi-tenant: langganan, billing, kuota AI.
- Ekspor lanjutan dan integrasi akuntansi.
- Portal/aplikasi pelanggan.
- PWA offline-first untuk kasir (antrean transaksi saat koneksi buruk).

### 4.3 Di Luar Ruang Lingkup MVP

- Melatih foundation model sendiri.
- AI mengubah saldo, menghapus transaksi, refund, atau mengubah harga otomatis tanpa kontrol.
- Integrasi semua penyedia pembayaran dan pesan sekaligus.
- Aplikasi native Android/iOS.
- Prediksi bisnis kompleks sebelum data historis memadai.
- Pengganti penuh sistem akuntansi profesional.

---

## 5. Modul dan Persyaratan Fungsional

Prioritas: **P0** wajib MVP · **P1** penting setelah MVP stabil · **P2** pengembangan lanjutan.

### 5.1 Autentikasi dan Hak Akses (P0)

- Login/logout dengan JWT access token berumur pendek + refresh token (rotasi, disimpan aman; refresh token dalam httpOnly cookie).
- Password di-hash dengan **argon2id** atau **bcrypt** (cost sesuai rekomendasi terkini).
- RBAC: role dan permission granular per modul/aksi.
- Backend memvalidasi identitas dan izin pada setiap endpoint; frontend hanya membantu UX.
- Pengguna hanya melihat data tenant dan cabang yang menjadi kewenangannya.
- Owner mengelola akun karyawan dan peran.
- Reset/ubah password yang aman (token sekali pakai, kedaluwarsa).
- Rate limiting dan penguncian sementara pada percobaan login gagal berulang.

**Kriteria penerimaan**
- Tanpa autentikasi, endpoint privat mengembalikan 401.
- Kasir tanpa izin tidak dapat membuka fitur Owner (403).
- Mengganti ID pada request tidak memberi akses ke data tenant lain (diuji otomatis).

### 5.2 Profil Bisnis dan Cabang

**P0**
- Owner mengatur nama bisnis, alamat, telepon, logo, mata uang (default IDR), zona waktu (default Asia/Jakarta), dan footer nota.
- MVP mendukung satu cabang; skema siap multi-cabang.
- Nomor pesanan dan laporan dapat ditelusuri ke tenant dan cabang.

**P1**
- Jam operasional, hari libur, konfigurasi per cabang.
- Transfer pesanan/stok antarcabang dengan pencatatan.

### 5.3 Manajemen Layanan dan Harga

**P0**
- CRUD layanan, aktif/nonaktif.
- Jenis: kiloan, satuan, sepatu, dry cleaning, kategori lain.
- Atribut: nama, satuan, harga, estimasi waktu pengerjaan, status aktif.
- Validasi harga dan kuantitas tidak negatif.
- Harga pada transaksi disimpan sebagai **snapshot**.

**P1**
- Harga berbeda per cabang/jenis pelanggan/tingkat layanan.
- Biaya tambahan (express, noda khusus, layanan tambahan).

### 5.4 Manajemen Pelanggan

**P0**
- Tambah, lihat, ubah, cari pelanggan (nama, telepon).
- Data: nama, nomor telepon (dinormalisasi ke format E.164, mis. +628…), alamat opsional, catatan.
- Riwayat transaksi pelanggan sesuai izin.
- Pencegahan duplikasi: unique constraint `(tenant_id, phone_normalized)`.

**P1**
- Poin loyalitas, voucher, segmentasi, catatan komplain, riwayat komunikasi.

### 5.5 POS dan Manajemen Pesanan

**P0**
- Pesanan dengan satu atau banyak item layanan.
- Nomor pesanan unik dan mudah dicari (format usulan: `LDR-{KODECABANG}-{YYMMDD}-{URUT}`, urutan dihasilkan backend secara atomik per cabang per hari).
- Data pesanan: pelanggan, cabang, waktu masuk, estimasi selesai, item (kuantitas/berat, harga snapshot), diskon, biaya tambahan, total, catatan, pembuat.
- Kalkulasi subtotal, diskon, biaya tambahan, dan total dilakukan **di backend**.
- Pembayaran penuh dan sebagian.
- Status pembayaran: `UNPAID`, `PARTIAL`, `PAID`, `REFUNDED`.
- Nota dapat dicetak (CSS print / thermal 58–80 mm) atau disimpan sebagai PDF.
- Pencarian berdasarkan nomor, nama, atau telepon.
- Pembatalan/koreksi sesuai izin dan masuk audit log.
- Transaksi berstatus dibayar tidak dapat dihapus permanen lewat operasi biasa (soft delete / status `CANCELLED`).
- **Idempotensi**: request pembuatan pesanan/pembayaran membawa `Idempotency-Key` agar klik ganda atau retry jaringan tidak menggandakan data.

**Kriteria penerimaan**
- Total di frontend = hasil backend.
- Nomor pesanan tidak pernah duplikat (diuji pada pembuatan paralel).
- Perubahan harga layanan tidak mengubah transaksi lama.
- Permintaan gagal/di-retry tidak menghasilkan pesanan atau pembayaran ganda.

### 5.6 Alur Produksi Laundry

**P0** — Status awal (dapat dikonfigurasi):

```text
DITERIMA → MENUNGGU_PROSES → DICUCI → DIKERINGKAN → DISETRIKA → DIKEMAS → SIAP_DIAMBIL → SELESAI
                                                                                   ↘ DIBATALKAN (dari status tertentu, perlu izin)
```

- Transisi status diatur oleh **state machine** di backend; transisi tidak sah ditolak (409/422).
- Layanan tertentu dapat melewati langkah (mis. cuci kering tanpa setrika) melalui konfigurasi alur per layanan.
- Pesanan terlambat ditandai berdasarkan estimasi selesai dan status.
- Setiap perubahan status menyimpan waktu, pengguna, dan catatan opsional.
- Catatan khusus per pesanan (instruksi perawatan).

**P1**
- Penugasan operator, checklist QC.
- Pencatatan pakaian rusak/hilang/komplain.
- Foto pendukung (batas ukuran, akses terkontrol, penyimpanan object storage).

### 5.7 Pembayaran dan Pengeluaran

**P0**
- Catat pembayaran tunai, transfer, QRIS manual, atau metode lain sebagai catatan.
- Simpan jumlah, metode, waktu, pencatat, referensi opsional, status.
- Pembayaran sebagian dihitung terhadap sisa tagihan.
- Pengeluaran: kategori, nominal, tanggal, deskripsi, bukti opsional.
- Rekap pemasukan dan pengeluaran tercatat.
- Koreksi pembayaran, pembatalan, refund memerlukan izin khusus + audit log.

**Catatan**
- Pencatatan transfer/QRIS di MVP bukan verifikasi otomatis.
- Laporan bukan pengganti akuntansi/pajak profesional.
- Laba bersih hanya ditampilkan final jika cakupan biaya didefinisikan; jika belum, tampilkan **"Estimasi Laba"** dengan label jelas.

### 5.8 Dashboard dan Laporan

**P0**
- Ringkasan: jumlah pesanan, omzet tercatat, pesanan aktif, terlambat, belum lunas, pengeluaran.
- Pilih rentang tanggal; zona waktu bisnis yang benar.
- Drill-down ke data pendukung sesuai izin.
- Membedakan: **omzet transaksi**, **kas diterima**, **pengeluaran tercatat**, **estimasi laba**.
- Aturan kalkulasi metrik terdokumentasi (lihat Bab 8, aturan 10–11).

**P1**
- Perbandingan antarperiode dan antarcabang.
- Layanan terlaris, rata-rata nilai pesanan, pelanggan aktif, tren keterlambatan.
- Ekspor CSV/PDF/XLSX (diproses di background job untuk data besar).

### 5.9 Inventaris dan Peralatan (P1)

- Barang, satuan, stok, batas minimum, harga beli, pemasok.
- Pergerakan stok: masuk, pemakaian, penyesuaian, penghapusan (wajib alasan).
- Peringatan stok minimum.
- Mesin/peralatan: jadwal perawatan, kerusakan, biaya, riwayat servis.
- Penyesuaian stok perlu izin + audit log.

### 5.10 CRM dan Notifikasi

**P1**
- Daftar pelanggan dengan pesanan siap diambil.
- Template pesan (nomor pesanan, nama, status).
- Pesan ditinjau manusia sebelum dikirim pada tahap awal.
- WhatsApp via API resmi/penyedia sah; tidak memakai otomasi yang melanggar ketentuan platform.
- Status pengiriman/kegagalan dicatat jika didukung provider.
- Pengiriman dilakukan lewat **job queue** (retry, backoff, dead-letter).

**P2**
- Pengingat otomatis terjadwal, kampanye promo, loyalitas, voucher.

---

## 6. AI Business Assistant

### 6.1 Tujuan

AI membantu pengguna memahami data bisnis dan menyelesaikan pekerjaan yang didukung sistem. **Database dan backend adalah sumber kebenaran**, bukan AI.

### 6.2 Provider dan Model

- Provider awal: **Google Gemini API via Google AI Studio (free tier)**.
- Pemilihan model berdasarkan: dukungan function calling, kualitas Bahasa Indonesia, batas kuota/rate limit, latensi, dan ketentuan terkini.
- Nama model **tidak di-hardcode**; gunakan environment variable.
- API key hanya di backend (secret), tidak pernah dikirim ke browser.
- Backend menyediakan **abstraksi provider** (`LLMProvider` interface) agar mudah berpindah ke model/provider lain.
- Tangani: kuota habis (429), timeout, respons tidak valid, gangguan provider, dan safety block.

### 6.3 Strategi Menghadapi Keterbatasan Free Tier

| Risiko free tier | Strategi |
|---|---|
| Rate limit (RPM/RPD) dan token terbatas | Rate limit per pengguna & global di Redis; antrean permintaan; batas input/output token |
| Kuota dan syarat dapat berubah | Konfigurasi via env; jangan menjanjikan ketersediaan; siapkan provider cadangan |
| Konten pada free tier dapat digunakan Google untuk peningkatan produk (verifikasi ketentuan terbaru) | **Minimalkan data pribadi** ke model: kirim agregat/ringkasan, samarkan nama/telepon pelanggan (mis. "Pelanggan A"); sebelum produksi komersial, pindah ke tier berbayar dan tinjau kebijakan data |
| Latensi tinggi | Streaming respons (SSE), indikator proses, timeout tegas |
| Pertanyaan berulang boros kuota | **Cache hasil tool** (TTL pendek) dan cache jawaban untuk pertanyaan identik dalam periode data yang sama |
| Provider down | Fallback: pesan informatif + tautan ke dashboard/laporan manual; fitur kasir tetap jalan |

### 6.4 Cara Kerja AI

1. Pengguna mengirim pertanyaan melalui UI chat.
2. Backend memvalidasi identitas, tenant, cabang, role, dan kuota pengguna.
3. Backend mengirim pertanyaan + definisi tool yang diizinkan ke Gemini.
4. Jika model meminta tool, backend memvalidasi nama fungsi dan argumen (schema validation).
5. Backend menjalankan fungsi internal dengan scope data pengguna.
6. Hasil terstruktur dikirim kembali ke model.
7. Backend memvalidasi respons dan mengalirkannya ke frontend (SSE).
8. Jika ada perubahan data, terapkan otorisasi dan konfirmasi eksplisit.

### 6.5 Tool AI untuk MVP (read-only)

| Tool | Fungsi | Izin minimum |
|---|---|---|
| `get_revenue_summary` | Ringkasan omzet tercatat per rentang waktu | `report.view` |
| `get_order_summary` | Jumlah pesanan per status dan periode | `order.view` |
| `search_orders` | Pencarian pesanan sesuai kriteria | `order.view` |
| `get_overdue_orders` | Pesanan melewati estimasi selesai | `order.view` |
| `get_unpaid_orders` | Pesanan belum lunas | `payment.view` |
| `get_service_performance` | Performa layanan berdasarkan transaksi | `report.view` |
| `get_expense_summary` | Ringkasan pengeluaran tercatat | `expense.view` |
| `calculate_service_price` | Simulasi harga dari biaya, satuan, target margin | `service.view` |
| `draft_customer_message` | Menyusun draft pesan (tanpa mengirim) | `customer.view` |
| `get_inventory_alerts` | (Pasca-MVP) peringatan stok | `inventory.view` |

Aturan tool:
- Skema input/output jelas (JSON Schema), validasi argumen, pemeriksaan izin, batas hasil (mis. maks 50 baris), parameterized query.
- **Dilarang menerima SQL mentah** dari model.
- Hasil tool menyertakan `period`, `source`, dan `generated_at` agar AI dapat menyebut dasar angkanya.

### 6.6 Contoh Pertanyaan

- "Berapa omzet saya minggu ini dibanding minggu lalu?"
- "Tampilkan pesanan yang terlambat hari ini."
- "Pesanan mana yang belum lunas?"
- "Layanan apa yang paling banyak dipesan bulan ini?"
- "Berapa rata-rata nilai transaksi saya?"
- "Bantu hitung harga per kilogram jika biaya operasionalnya sekian."
- "Buatkan pesan WhatsApp untuk pelanggan yang belum mengambil laundry."
- "Apa saja hal yang perlu saya perhatikan dari laporan bulan ini?"

### 6.7 Guardrails

- AI tidak boleh mengarang angka, transaksi, status pembayaran, atau data pelanggan.
- Jika data tidak ada atau tool gagal, AI menyatakan keterbatasan.
- Semua angka bisnis bersumber dari hasil tool terverifikasi.
- Akses data sesuai tenant, cabang, dan izin pengguna.
- Isi pesan pelanggan atau data database **bukan** instruksi sistem (mitigasi prompt injection); data tool diberi delimiter dan diperlakukan sebagai data.
- Output model dianggap tidak tepercaya dan divalidasi.
- Tidak ada SQL bebas, perubahan role, ekspor seluruh database, atau penembusan akses.
- Pengiriman pesan, perubahan harga massal, refund, pembatalan, penghapusan, dan perubahan finansial memerlukan konfirmasi eksplisit + izin tepat.
- **MVP: AI read-only + draft**; aksi tulis ditambahkan setelah otorisasi dan audit matang.
- Batas panjang input/output, jumlah tool call (`AI_MAX_TOOL_CALLS`), waktu eksekusi, dan kuota harian per pengguna.
- Jangan menjanjikan layanan AI gratis selalu tersedia.

### 6.8 Riwayat dan Evaluasi

**P0**
- UI chat menampilkan riwayat percakapan pengguna.
- Catat metadata: waktu, status, penggunaan token (jika tersedia), tool yang dijalankan — tanpa menyimpan rahasia.
- Minimalkan data pribadi di log dan terapkan kebijakan retensi (mis. 90 hari, dapat dikonfigurasi).
- Kegagalan AI tidak boleh menggagalkan proses kasir/transaksi.

**P1**
- Feedback 👍/👎 pada jawaban.
- Dataset uji: pertanyaan normal, data kosong, permintaan tidak sah, prompt injection, kegagalan provider.
- Pantau akurasi angka, keberhasilan tool, latensi, kuota, dan jawaban tanpa dukungan data.

---

## 7. User Stories Utama

### Owner
- Melihat ringkasan performa bisnis untuk memantau kondisi usaha.
- Mencatat pengeluaran untuk memahami biaya operasional.
- Bertanya kepada AI tentang laporan agar tidak menyaring data manual.
- Mengelola layanan dan harga sesuai strategi.
- Mengelola karyawan dan hak akses agar data aman.

### Kasir
- Membuat pesanan cepat agar pelanggan tidak menunggu.
- Mencari pelanggan dan pesanan lama.
- Mencatat pembayaran agar tagihan akurat.
- Melihat status dan estimasi selesai untuk informasi ke pelanggan.

### Operator
- Melihat daftar pekerjaan yang harus diproses.
- Memperbarui status pengerjaan.
- Mencatat masalah pakaian agar dapat ditelusuri.

### AI Assistant
- Bertanya informasi bisnis dengan bahasa alami.
- Melihat sumber angka dan periode data yang digunakan AI.
- Meminta AI membuat draft pesan pelanggan.
- Dimintai konfirmasi sebelum tindakan penting.

---

## 8. Aturan Bisnis

1. Setiap transaksi memiliki ID internal dan nomor pesanan unik.
2. Harga layanan pada transaksi disimpan sebagai snapshot.
3. Total dihitung di backend; frontend bukan sumber kebenaran.
4. Pembayaran dapat terdiri dari beberapa catatan pembayaran.
5. Status pembayaran dihitung dari total tagihan, pembayaran valid, dan refund.
6. Pesanan dibatalkan tetap memiliki jejak audit dan tidak dihapus permanen lewat operasi biasa.
7. Perubahan status mengikuti transisi yang diizinkan (state machine).
8. Waktu disimpan dalam **UTC** dan ditampilkan sesuai zona waktu bisnis.
9. Setiap data bisnis terkait `tenant_id` dan, bila relevan, `branch_id`.
10. Laporan harus mendefinisikan basis tanggal: tanggal pesanan, tanggal pembayaran, atau tanggal penyelesaian.
11. Omzet, kas diterima, pengeluaran, dan estimasi laba adalah metrik berbeda.
12. Diskon, refund, koreksi pembayaran, dan penyesuaian stok memerlukan otorisasi dan audit.
13. AI tidak boleh melewati aturan bisnis yang berlaku pada UI normal.
14. Aksi via AI memakai service layer backend yang sama dengan aplikasi utama.
15. **Uang disimpan sebagai bilangan bulat (satuan rupiah, `BIGINT`) atau `DECIMAL(15,2)`; tidak memakai `FLOAT/DOUBLE`.**
16. Berat kiloan disimpan dengan presisi tetap (mis. `DECIMAL(10,3)`).

---

## 9. Arsitektur Teknis

### 9.1 Tech Stack

| Lapisan | Teknologi | Catatan |
|---|---|---|
| Frontend | **Next.js (App Router) + React + TypeScript** | SSR/CSR sesuai kebutuhan, PWA-ready |
| UI | Tailwind CSS + shadcn/ui (atau Radix UI) | Komponen aksesibel, konsisten |
| State/data fetching | TanStack Query, React Hook Form + Zod | Cache klien, validasi form |
| Backend | **Go + framework HTTP (Gin / Fiber / Echo)** | Rekomendasi: Gin (ekosistem luas) atau Fiber (throughput tinggi) |
| Akses data | **sqlc** (type-safe SQL) atau GORM | Rekomendasi sqlc untuk performa dan kontrol query |
| Migration | golang-migrate atau goose | Versi skema di repository |
| Database | **MySQL 8.x (InnoDB, utf8mb4)** | Primary + opsional read replica |
| Cache & rate limit | **Redis** | Cache laporan, rate limiting, session/blacklist token |
| Job queue | **Asynq** (Redis-based) | Notifikasi, ekspor, laporan berat, retry |
| Object storage | S3-compatible (MinIO / Cloudflare R2 / S3) | Foto, bukti pengeluaran, PDF |
| LLM | **Gemini API (Google AI Studio, free tier)** | Via adapter provider |
| Reverse proxy | **Nginx / Caddy / Traefik** | TLS, kompresi, rate limit dasar, HTTP/2 |
| Kontainer | **Docker + Docker Compose** | Multi-stage build; Kubernetes opsional di fase SaaS |
| CI/CD | **GitHub Actions** (atau GitLab CI) | Lihat Bab 11 |
| Dokumentasi API | OpenAPI 3 (swaggo/swag atau oapi-codegen) | Swagger UI di non-produksi |
| Observability | Prometheus, Grafana, Loki, OpenTelemetry, Sentry | Lihat Bab 10 |
| Version control | Git (trunk-based / GitHub Flow) | Conventional Commits |

> **Catatan MySQL**: Pastikan MySQL 8.0+ (window function, CTE, `JSON`, `SKIP LOCKED`). Gunakan `utf8mb4_0900_ai_ci` (atau `utf8mb4_unicode_ci`), engine InnoDB, dan simpan waktu sebagai `DATETIME(3)`/`TIMESTAMP` dalam UTC.

### 9.2 Prinsip Arsitektur

- Frontend **tidak** memanggil Gemini langsung; backend satu-satunya penyimpan API key dan penghubung ke provider AI.
- **Modular monolith** pada MVP (satu deployable, modul terpisah secara logis) agar sederhana dan cepat; dapat dipecah menjadi service di masa depan.
- Modul domain: `auth`, `tenant`, `branch`, `service`, `customer`, `order`, `production`, `payment`, `expense`, `report`, `ai`, `audit`, `notification`.
- Pola lapisan: `handler → service → repository`; logika bisnis hanya di service.
- Semua perubahan transaksi lewat service layer dan **transaksi database**.
- AI hanya memanggil tool terdefinisi.
- Provider AI dibungkus adapter.
- Semua endpoint: autentikasi, otorisasi, validasi, pembatasan data per tenant.
- Proses bisnis inti berjalan tanpa AI.
- Aplikasi mengikuti **12-factor** (konfigurasi via environment, log ke stdout, stateless).

### 9.3 Struktur Repository (Usulan)

```text
laundryos/
├── apps/
│   ├── web/                 # Next.js
│   └── api/                 # Go backend
│       ├── cmd/server/
│       ├── internal/
│       │   ├── auth/ order/ payment/ report/ ai/ audit/ ...
│       │   ├── platform/    # db, redis, logger, config, llm adapter
│       │   └── http/        # router, middleware
│       ├── migrations/
│       └── openapi/
├── deploy/
│   ├── docker/              # Dockerfile, compose
│   ├── nginx/
│   └── monitoring/
├── .github/workflows/
└── docs/
```

### 9.4 Alur Integrasi AI

```text
Web Frontend (Next.js)
    │  (SSE/stream)
    ▼
Go Backend API ── Auth & RBAC ── Rate Limit (Redis)
    │
    ▼
AI Orchestrator ──► Cache (Redis)
    │
    ▼
Gemini API (via adapter)
    │  tool call request
    ▼
Tool Validator (schema + permission)
    │
    ▼
Domain Services ──► MySQL (read-only query terparameter)
    │  tool result
    ▼
Gemini API → jawaban → validasi → Frontend
```

### 9.5 Environment Variables (Contoh)

```env
APP_ENV=development
APP_PORT=8080
APP_BASE_URL=http://localhost:3000

# MySQL
DB_DSN=laundry:password@tcp(mysql:3306)/laundryos?parseTime=true&loc=UTC&charset=utf8mb4
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME_MINUTES=30

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=

# Auth
JWT_ACCESS_SECRET=change_me
JWT_REFRESH_SECRET=change_me
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_DAYS=14

# AI
GEMINI_API_KEY=your_secret_key
GEMINI_MODEL=your_selected_model
AI_REQUEST_TIMEOUT_SECONDS=30
AI_MAX_TOOL_CALLS=5
AI_MAX_INPUT_CHARS=2000
AI_DAILY_LIMIT_PER_USER=20
AI_GLOBAL_RPM_LIMIT=10
AI_CACHE_TTL_SECONDS=300

# Storage
S3_ENDPOINT=
S3_BUCKET=
S3_ACCESS_KEY=
S3_SECRET_KEY=
```

Nilai di atas hanyalah contoh, bukan kredensial nyata. Jangan commit `.env` berisi secret; gunakan GitHub Secrets / secret manager untuk CI/CD dan produksi.

---

## 10. Peningkatan Performa, Skalabilitas, dan Observability

### 10.1 Target Performa (Awal, Diuji Ulang Setelah Load Test)

| Metrik | Target |
|---|---|
| API operasional normal (tanpa provider eksternal) | p95 < 300–500 ms |
| Pencarian pesanan | p95 < 1 detik pada volume MVP |
| Dashboard (dengan cache) | p95 < 800 ms |
| Next.js LCP (halaman utama kasir) | < 2,5 detik pada 4G |
| Pembuatan pesanan (POS) | < 1 detik end-to-end |
| AI | Token pertama via streaming ≤ beberapa detik; ada indikator proses |
| UI saat AI down | Tetap berfungsi penuh |

### 10.2 Teknologi dan Teknik Peningkat Performa

| Teknologi / Teknik | Manfaat | Prioritas |
|---|---|---|
| **Redis cache** (hasil laporan, daftar layanan, konfigurasi bisnis, hasil tool AI) | Mengurangi beban MySQL dan kuota Gemini | P0 |
| **Redis rate limiting** (login, AI, endpoint sensitif) | Anti-abuse dan melindungi kuota free tier | P0 |
| **Indeks MySQL terencana** (komposit `tenant_id, branch_id, created_at`, `order_number`, `phone_normalized`, `status`) + analisis `EXPLAIN` | Query cepat | P0 |
| **Connection pooling** (`database/sql` tuning) | Stabil saat konkurensi | P0 |
| **Pagination berbasis cursor/keyset** untuk daftar besar | Menghindari `OFFSET` lambat | P0 |
| **Kompresi (gzip/brotli) dan HTTP/2** di reverse proxy | Payload lebih kecil | P0 |
| **Next.js optimasi**: Server Components, code splitting, `next/image`, ISR/caching, font optimization | Muat halaman cepat | P0 |
| **TanStack Query** (cache, dedupe, stale-while-revalidate) | UI responsif, request lebih sedikit | P0 |
| **Asynq job queue** (ekspor, notifikasi, laporan berat) | Menjauhkan pekerjaan berat dari request | P1 |
| **Tabel ringkasan/materialized aggregate** (mis. `daily_sales_summary` diperbarui lewat job) | Laporan cepat tanpa scan transaksi penuh | P1 |
| **SSE streaming** untuk jawaban AI | Persepsi cepat | P1 |
| **CDN** (Cloudflare) untuk aset statis & gambar | Latensi rendah, perlindungan DDoS dasar | P1 |
| **Read replica MySQL** untuk laporan | Memisahkan beban baca | P2 |
| **ProxySQL** | Pooling & routing baca/tulis | P2 |
| **Partisi tabel** (mis. audit/log per bulan) | Pengelolaan data besar | P2 |
| **PWA + service worker** (offline queue kasir) | Tahan koneksi buruk | P2 |
| **Horizontal scaling** (beberapa replika API di belakang load balancer) | Kapasitas | P2 |
| **Go `pprof` + benchmark** | Optimasi bottleneck berbasis data | P1 |

### 10.3 Konsistensi Data pada MySQL

- Gunakan transaksi (`BEGIN … COMMIT`) untuk pembuatan pesanan + item + pembayaran + audit.
- Level isolasi default InnoDB (REPEATABLE READ) cukup bila disertai `SELECT … FOR UPDATE` pada operasi kritis (pembayaran terhadap pesanan yang sama, pembuatan nomor urut).
- Nomor urut pesanan: tabel `order_sequences (tenant_id, branch_id, date, last_number)` dengan `INSERT … ON DUPLICATE KEY UPDATE` atau locking baris.
- Constraint: `FOREIGN KEY`, `CHECK` (MySQL 8.0.16+), `UNIQUE (tenant_id, order_number)`.
- Hindari N+1 query; gunakan join/batch.

### 10.4 Observability

| Kebutuhan | Alat |
|---|---|
| Metrik (latensi, error rate, throughput, pool DB, kuota AI) | Prometheus + Grafana |
| Log terstruktur (JSON, correlation/request ID) | zerolog/zap → Loki (atau ELK) |
| Tracing | OpenTelemetry (opsional di MVP) |
| Error tracking | Sentry (frontend & backend) |
| Health check | `/healthz` (liveness), `/readyz` (cek DB, Redis) |
| Uptime monitoring | Uptime Kuma / UptimeRobot |
| Alert | Alertmanager / notifikasi Telegram/Slack |

### 10.5 Pengujian Performa

- **k6** atau **vegeta** untuk load test skenario: login, pembuatan pesanan, pencarian, dashboard.
- Lighthouse CI untuk frontend.
- Data uji realistis (mis. 100 ribu pesanan) sebelum pilot.

---

## 11. Docker dan CI/CD

### 11.1 Docker

- **Backend**: multi-stage build (`golang` → `distroless`/`alpine`), binary statis, user non-root.
- **Frontend**: multi-stage build Next.js (`output: "standalone"`).
- **Docker Compose** untuk pengembangan lokal: `web`, `api`, `mysql`, `redis`, `minio`, `nginx` (opsional), `prometheus/grafana` (opsional).
- Healthcheck pada setiap service; volume terpisah untuk data MySQL dan MinIO.
- Image diberi tag `git-sha` dan `semver`; simpan di registry (GitHub Container Registry / Docker Hub).

Contoh layanan Compose (ringkas):

```yaml
services:
  api:
    build: ./apps/api
    env_file: .env
    depends_on: [mysql, redis]
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/healthz"]
  web:
    build: ./apps/web
    depends_on: [api]
  mysql:
    image: mysql:8
    volumes: [mysql_data:/var/lib/mysql]
  redis:
    image: redis:7-alpine
volumes:
  mysql_data:
```

### 11.2 Pipeline CI (setiap Pull Request)

1. Checkout dan cache dependensi (Go modules, pnpm).
2. **Lint & format**: `golangci-lint`, `gofmt`; `eslint`, `prettier`, `tsc --noEmit`.
3. **Unit test** Go (`go test -race -cover`) dan frontend (Vitest/Jest).
4. **Integration test** dengan service container MySQL + Redis (termasuk uji isolasi tenant & RBAC).
5. **Migration check**: migrasi naik/turun berhasil pada DB kosong.
6. **Security scan**: `govulncheck`, `gosec`, `npm audit`/`pnpm audit`, **Trivy** (image), **gitleaks** (secret scanning).
7. **Build** image Docker (tanpa push untuk PR).
8. Validasi OpenAPI (kontrak API tidak rusak).
9. (Opsional) E2E Playwright terhadap stack Compose.

### 11.3 Pipeline CD

| Lingkungan | Pemicu | Tindakan |
|---|---|---|
| **Development/Preview** | Merge ke `develop` atau PR | Deploy otomatis ke server dev |
| **Staging** | Merge ke `main` | Build, push image, jalankan migrasi, deploy otomatis, smoke test |
| **Production** | Tag rilis (`vX.Y.Z`) + persetujuan manual | Backup DB → migrasi → deploy → smoke test → rollback otomatis bila gagal |

- Strategi rilis: **rolling/blue-green** menggunakan Docker Compose + Traefik/Nginx (atau Kubernetes di fase lanjut).
- Deploy via SSH/GitHub Actions runner atau Watchtower/Argo (sesuai infrastruktur).
- Migrasi database **backward-compatible** (expand → migrate → contract) agar rollback aman.
- Secret dikelola lewat GitHub Secrets / secret manager; tidak ada secret di image.
- Rollback: redeploy image tag sebelumnya.
- Semua deploy dicatat (versi, waktu, pelaku).

### 11.4 Kebijakan Branch

- `main` selalu dapat dirilis; perlindungan branch (review wajib, CI hijau).
- Fitur lewat branch pendek + Pull Request; Conventional Commits untuk changelog otomatis.

---

## 12. Entitas Data Utama (MySQL)

Rancangan awal; dapat disesuaikan saat desain database.

| Tabel | Keterangan |
|---|---|
| `tenants` | Bisnis pengguna platform |
| `branches` | Cabang milik tenant |
| `users` | Akun pengguna |
| `roles`, `permissions`, `role_permissions` | RBAC |
| `customers` | Pelanggan (`phone_normalized` + unique per tenant) |
| `services` | Jenis layanan |
| `orders` | Data utama pesanan |
| `order_items` | Rincian layanan (harga snapshot) |
| `order_status_history` | Riwayat status |
| `order_sequences` | Penghasil nomor urut pesanan |
| `payments` | Catatan pembayaran |
| `refunds` | Refund (jika diaktifkan) |
| `expenses`, `expense_categories` | Pengeluaran |
| `employees` | Profil karyawan (jika terpisah) |
| `inventory_items`, `inventory_movements` | Inventaris (pasca-MVP) |
| `equipment`, `maintenance_records` | Peralatan (pasca-MVP) |
| `ai_conversations`, `ai_messages` | Riwayat percakapan AI (dengan retensi) |
| `ai_tool_logs` | Metadata tool call |
| `audit_logs` | Catatan tindakan penting |
| `idempotency_keys` | Mencegah request ganda |
| `daily_sales_summary` | Agregat harian untuk laporan cepat (P1) |
| `subscriptions`, `plans` | Tahap SaaS |
| `notifications` | Status notifikasi |

**Prinsip data**
- Primary key konsisten (disarankan `BIGINT UNSIGNED AUTO_INCREMENT` atau UUID v7/ULID; bila publik gunakan ID non-sekuensial untuk menghindari enumerasi).
- `tenant_id` pada semua entitas bisnis; `branch_id` pada data operasional.
- Indeks komposit dengan `tenant_id` di depan.
- Constraint `CHECK`, `NOT NULL`, `FOREIGN KEY`, `UNIQUE`.
- Kolom `created_at`, `updated_at` (UTC), `deleted_at` untuk soft delete bila sesuai.
- Strategi retensi per jenis data (log AI, audit, bukti).
- Backend wajib menerapkan filter tenant; tidak boleh bergantung pada frontend.

---

## 13. Konvensi API

- REST, JSON, versi di path: `/api/v1/...`.
- Autentikasi: `Authorization: Bearer <access_token>`.
- Respons error standar:
  ```json
  { "error": { "code": "ORDER_INVALID_TRANSITION", "message": "Status tidak dapat diubah dari X ke Y", "details": {} } }
  ```
- Pagination: cursor (`?cursor=&limit=`) untuk daftar besar.
- Filter tanggal memakai ISO 8601; server menafsirkan sesuai zona waktu bisnis.
- Header `Idempotency-Key` untuk `POST` pesanan/pembayaran.
- Header `X-Request-ID` untuk korelasi log.
- Dokumentasi OpenAPI dipublikasikan otomatis dari kode/spesifikasi.
- Rate limit dikembalikan dengan `429` + `Retry-After`.

---

## 14. Halaman dan Navigasi Utama

### Owner/Manager
Dashboard · POS/Buat Pesanan · Daftar Pesanan · Detail Pesanan · Produksi Laundry · Pelanggan · Layanan dan Harga · Pembayaran · Pengeluaran · Laporan · AI Business Assistant · Karyawan dan Hak Akses · Pengaturan Bisnis · Inventaris dan Peralatan (fase berikutnya)

### Kasir
Dashboard operasional sederhana · Buat Pesanan · Daftar Pesanan · Detail Pesanan · Pelanggan · Pembayaran · AI Assistant (akses terbatas, jika diaktifkan)

### Operator
Daftar pekerjaan · Detail pekerjaan · Perbarui status · Catatan masalah dan QC

### Super Admin (fase SaaS)
Dashboard platform · Daftar tenant · Paket dan langganan · Status layanan dan kuota · Audit dan dukungan platform

---

## 15. Desain UI/UX

- Gaya profesional, bersih, modern, berorientasi operasional.
- Hindari dashboard generik yang ramai (gradien berlebihan, kartu dekoratif, grafik tanpa tujuan).
- Prioritaskan hierarki informasi, keterbacaan, kecepatan input, dan navigasi konsisten.
- POS meminimalkan langkah untuk transaksi umum (pencarian pelanggan cepat, pintasan keyboard, tombol layanan favorit).
- Label status pesanan dan pembayaran konsisten; warna **selalu disertai teks**.
- Dashboard hanya menampilkan metrik yang membantu keputusan.
- AI Assistant: konteks bisnis jelas, indikator proses, sumber/periode data, pesan error informatif, tombol "salin" dan "tanya lanjutan".
- Responsif: laptop, tablet, ponsel; target sentuh memadai untuk kasir.
- Aksesibilitas: kontras, fokus keyboard, label form, standar WCAG 2.1 AA sebagai target.
- State loading (skeleton), kosong, dan error pada setiap halaman data.

---

## 16. Kebutuhan Nonfungsional

### 16.1 Keamanan
- HTTPS (TLS) di seluruh deployment; HSTS.
- Hashing password argon2id/bcrypt.
- Secret di environment/secret manager.
- RBAC dan isolasi tenant pada setiap endpoint.
- Parameterized query/ORM; validasi input di backend.
- Validasi upload: ukuran, tipe (MIME sniffing), nama file, hak akses; penyimpanan di luar webroot.
- Rate limiting: autentikasi, AI, endpoint sensitif.
- Header keamanan (CSP, X-Content-Type-Options, X-Frame-Options), CORS ketat, proteksi CSRF untuk cookie.
- Audit log perubahan penting.
- Jangan mencatat password, token, API key, atau data pembayaran sensitif di log.
- Backup MySQL terjadwal (mysqldump/xtrabackup), terenkripsi, dan **uji restore berkala**.
- Pemindaian dependensi dan image di CI.
- Ikuti OWASP Top 10 dan OWASP API Security Top 10 sebagai checklist.

### 16.2 Ketersediaan dan Pemulihan
- Kesalahan AI tidak merusak/menggandakan transaksi.
- Transaksi database untuk operasi kritis.
- Graceful shutdown pada backend.
- Backup terjadwal + prosedur restore terdokumentasi; target RPO/RTO ditentukan sebelum produksi.
- Error tercatat dan dapat ditelusuri tanpa membocorkan rahasia.

### 16.3 Usability dan Aksesibilitas
- UI Bahasa Indonesia; format mata uang `Rp 1.234.567` dan tanggal lokal.
- Form transaksi dioptimalkan untuk kasir.
- Validasi dan pesan error jelas.
- AI menampilkan periode data dan konteks jawaban.

### 16.4 Privasi dan Kepatuhan
- Kumpulkan data seperlunya (minimisasi).
- Kebijakan retensi dan penghapusan data.
- Jelaskan penggunaan AI dan jenis data yang dikirim ke provider eksternal.
- Anonimkan/samarkan identitas pelanggan pada data yang dikirim ke model.
- Selaraskan dengan **UU Pelindungan Data Pribadi (UU PDP No. 27 Tahun 2022)**: dasar pemrosesan, hak subjek data, dan prosedur insiden.
- Tinjau kebijakan data Gemini API sebelum produksi (khususnya perbedaan free vs berbayar).

### 16.5 Maintainability
- Cakupan test pada logika kritis (kalkulasi, pembayaran, state machine, RBAC, isolasi tenant).
- Dokumentasi: README setup lokal, ADR (Architecture Decision Record), OpenAPI, ERD.
- Linting dan formatting otomatis (pre-commit hook).

---

## 17. Strategi Pengujian

| Jenis | Cakupan | Alat |
|---|---|---|
| Unit | Kalkulasi harga/total, state machine, validasi, status pembayaran | `go test`, Vitest |
| Integrasi | Repository + MySQL nyata, endpoint API, isolasi tenant, RBAC | testcontainers-go / service container CI |
| Kontrak API | Sesuai OpenAPI | schemathesis / dredd (opsional) |
| E2E | Alur login → buat pesanan → bayar → ubah status → laporan | Playwright |
| Keamanan | Otorisasi (IDOR), injection, rate limit | gosec, ZAP baseline, uji manual |
| AI | Pertanyaan normal, data kosong, akses tidak sah, prompt injection, provider gagal | Dataset evaluasi + mock provider |
| Performa | Load test skenario utama | k6 |
| Konkurensi | Pembuatan pesanan paralel, pembayaran ganda, retry idempotensi | Test khusus |
| UAT | Pilot dengan pemilik/karyawan laundry nyata | Checklist skenario |

---

## 18. Risiko dan Mitigasi

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Kuota/aturan Gemini free tier habis atau berubah | AI tidak tersedia | Rate limit, cache, fallback UI, adapter provider, rencana upgrade ke tier berbayar |
| Data pelanggan dikirim ke provider free tier | Risiko privasi | Minimisasi/penyamaran data, kebijakan privasi, evaluasi sebelum produksi |
| AI memberi jawaban tidak akurat | Keputusan bisnis keliru | Angka dari tool, validasi, evaluasi berkala, tampilkan sumber/periode |
| Kebocoran data antar bisnis | Pelanggaran privasi | Filter tenant di backend, uji IDOR otomatis, audit |
| AI melakukan tindakan tak diinginkan | Perubahan data salah | Read-only default, konfirmasi, izin khusus, audit |
| Data transaksi tidak konsisten | Laporan salah | Kalkulasi backend, constraint, transaksi DB, idempotensi, test |
| Performa MySQL menurun saat data besar | UI lambat | Indeks, keyset pagination, cache, tabel agregat, read replica |
| Deploy gagal / migrasi merusak data | Downtime | CI/CD dengan backup pra-deploy, migrasi backward-compatible, rollback |
| Integrasi pesan/pembayaran kompleks | Jadwal molor | Mulai manual, integrasi bertahap |
| Scope terlalu besar | MVP terlambat | Fokus P0, rilis bertahap |
| Backup tidak dapat dipulihkan | Kehilangan data | Uji restore terjadwal |
| Ketergantungan single server | Downtime | Monitoring, backup, rencana skala horizontal |

---

## 19. Definition of Done

Sebuah fitur selesai jika:

1. Alur utama berjalan sesuai requirement.
2. Validasi input diterapkan di backend.
3. Hak akses dan isolasi tenant diuji.
4. Kondisi error dan data kosong ditangani.
5. UI memiliki state loading, sukses, dan error.
6. Unit/integration test untuk logika penting lulus di CI.
7. Dokumentasi endpoint (OpenAPI) dan perubahan skema/migrasi diperbarui.
8. Perubahan penting masuk audit log.
9. Tidak ada secret di source code (lolos gitleaks).
10. Tidak merusak alur transaksi yang ada.
11. Tidak ada regresi performa signifikan pada endpoint terkait.
12. Telah ter-deploy ke staging melalui pipeline CI/CD dan lolos smoke test.

### Kriteria Selesai MVP

- Owner dan Kasir login sesuai peran.
- Owner mengatur layanan dan harga.
- Kasir membuat pesanan dan mencatat pembayaran.
- Status pesanan dapat dilacak hingga diambil.
- Owner melihat dashboard dan laporan dasar.
- Pengeluaran dapat dicatat.
- AI menjawab pertanyaan bisnis yang didukung melalui tool backend dan menyatakan keterbatasan saat data tidak ada.
- Akses dibatasi berdasarkan peran dan ruang lingkup bisnis.
- Aplikasi berjalan lewat Docker Compose dan ter-deploy melalui CI/CD.
- Dokumentasi setup lokal, environment, migrasi, dan pengujian tersedia.
- Backup dan restore telah diuji.

---

## 20. Rencana Implementasi

| Sprint | Fokus | Output Utama |
|---|---|---|
| **0 — Fondasi** | Repo, struktur, Docker Compose (MySQL, Redis), migrasi, config, linting, testing, **pipeline CI dasar**, OpenAPI | Skeleton berjalan lokal + CI hijau |
| **1 — Auth & Master Data** | Auth, RBAC, tenant/cabang, profil bisnis, layanan, pelanggan | Login & master data |
| **2 — POS & Pesanan** | Pembuatan transaksi, kalkulasi backend, idempotensi, nomor pesanan, nota, daftar/pencarian/detail | Transaksi end-to-end |
| **3 — Produksi & Pembayaran** | State machine status, riwayat, pembayaran penuh/sebagian, pengeluaran, audit log | Alur operasional lengkap |
| **4 — Dashboard & Laporan** | Ringkasan, laporan periode, cache Redis, verifikasi konsistensi metrik | Laporan akurat & cepat |
| **5 — AI Assistant** | Integrasi Gemini via adapter, chat UI + SSE, tool read-only, rate limit & cache, audit metadata, evaluasi | AI berbasis data |
| **6 — CD, Stabilitas & Pilot** | **Pipeline CD (staging/produksi)**, observability, load test (k6), uji keamanan, backup/restore, pilot pengguna nyata | Siap rilis |

Jadwal disesuaikan dengan kapasitas pengembang dan hasil validasi kebutuhan.

---

## 21. Asumsi dan Keputusan yang Perlu Dikonfirmasi

1. Nama SatSet adalah nama kerja.
2. MVP untuk satu bisnis dan satu cabang; skema siap ekspansi.
3. UI utama Bahasa Indonesia.
4. Aplikasi web responsif (PWA opsional).
5. Backend Go dengan framework (pilih: Gin / Fiber / Echo) dan database MySQL 8.
6. Akses data: sqlc atau GORM (perlu diputuskan tim).
7. Gemini API via Google AI Studio (free tier) sebagai LLM awal; evaluasi migrasi ke tier berbayar sebelum produksi komersial.
8. WhatsApp dan payment gateway belum wajib untuk MVP.
9. AI MVP: pertanyaan berbasis data, analisis, dan draft; aksi yang mengubah data dibatasi.
10. Definisi omzet, kas diterima, biaya, refund, dan laba disepakati sebelum laporan dianggap final.
11. Aturan harga, satuan, status, pembayaran sebagian, dan pembatalan divalidasi terhadap proses nyata laundry.
12. Target hosting (VPS/cloud) dan penyedia CI/CD (GitHub Actions/GitLab CI) ditentukan.
13. Redis dan job queue disetujui sebagai bagian dari stack MVP (disarankan).
14. Paket SaaS, harga, dan kebijakan kuota AI ditentukan setelah MVP tervalidasi.

---

## 22. Ringkasan Akhir

SatSet bertujuan menjadi sistem all-in-one bagi bisnis laundry dengan fondasi operasional yang kuat dan AI yang benar-benar berguna. Prioritas awal adalah transaksi, pesanan, pembayaran, dan laporan yang akurat serta cepat; setelah itu Gemini AI diintegrasikan melalui backend dengan tool/function calling yang aman, hemat kuota, dan menjaga privasi data.

Stack Next.js + Go + MySQL + Redis yang dikemas dengan Docker dan dikirim melalui CI/CD memberi dasar yang ringan, cepat, dan mudah dikembangkan menjadi SaaS.

**Prinsip utama produk:**

> **Operasional tetap berjalan tanpa AI. AI membantu memahami data dan mempercepat pekerjaan, tetapi tidak menggantikan validasi bisnis, kontrol akses, maupun keputusan penting pengguna.**
