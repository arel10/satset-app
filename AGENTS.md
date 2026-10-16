# AGENTS.md — SatSet Coding Agent Guidelines

> Dokumen ini adalah aturan kerja untuk coding agent yang mengembangkan SatSet.
> Baca dokumen ini dan PRD proyek sebelum membuat perubahan kode.
> Sumber kebutuhan produk utama: `doc/PRD_SatSet.md` (sesuaikan jika lokasi PRD berbeda).

## 1. Konteks Proyek

SatSet adalah aplikasi web internal untuk **satu bisnis laundry**, bukan SaaS multi-tenant.

Stack utama:
- Frontend: Next.js, React, TypeScript, Tailwind CSS, shadcn/ui.
- Backend: Go + Gin.
- ORM: GORM.
- Database: MySQL 8 / InnoDB.
- AI: Google Gemini API, hanya untuk Owner, melalui tools backend yang read-only.
- Deployment target: Docker Compose; hosting produksi belum diputuskan.

Peran aplikasi:
- **Owner**: akses seluruh fitur, laporan, pengeluaran, pengguna, pengaturan, dan AI Assistant.
- **Kasir**: membuat pesanan, mengelola pelanggan, mencatat pembayaran lunas, melihat pesanan, dan mencetak nota.
- **Operator**: melihat pekerjaan dan memperbarui status pengerjaan sesuai alur yang diizinkan.

Aturan bisnis penting:
- Pembayaran normal wajib lunas; tidak ada DP atau cicilan pada MVP.
- Harga dan nama layanan disimpan sebagai snapshot pada item pesanan.
- Total transaksi dihitung dan divalidasi di backend.
- Pesanan tidak dihapus permanen; pembatalan harus tercatat.
- AI tidak boleh mengubah, menghapus, membatalkan, atau merefund data.
- Semua waktu disimpan sebagai UTC dan ditampilkan dalam zona waktu `Asia/Jakarta`.
- Nilai uang disimpan sebagai integer Rupiah, bukan `FLOAT`/`DOUBLE`.
- Berat/jumlah pecahan menggunakan tipe desimal dengan presisi yang sesuai.
- Nama, telepon, alamat, dan catatan pelanggan merupakan data sensitif operasional; kirim ke Gemini hanya jika memang diperlukan.

## 2. Aturan Prioritas

Jika terdapat konflik, ikuti urutan berikut:
1. Instruksi eksplisit terbaru dari pengguna.
2. PRD proyek.
3. Dokumen `AGENTS.md` ini.
4. Pola dan konvensi kode yang sudah ada.
5. Preferensi umum agent.

Jangan mengubah kebutuhan produk, alur bisnis, stack, atau struktur data secara besar-besaran tanpa alasan kuat. Jika ada ambiguitas yang memengaruhi data, keamanan, pembayaran, atau arsitektur, ajukan pertanyaan singkat sebelum mengambil keputusan yang sulit dibatalkan. Untuk keputusan kecil, pilih solusi paling sederhana dan dokumentasikan asumsi.

## 3. Cara Memulai Tugas

Sebelum mengubah kode:
1. Periksa struktur repository dan status Git.
2. Baca PRD serta `AGENTS.md` ini.
3. Baca file terkait dan ikuti pola yang sudah digunakan proyek.
4. Identifikasi dampak ke frontend, backend, database, autentikasi, dan pengujian.
5. Buat rencana singkat untuk tugas yang lebih dari perubahan kecil.
6. Jangan mengasumsikan fitur atau file sudah ada sebelum memeriksanya.

Saat menyelesaikan tugas:
1. Jalankan formatter, linter, type-check, dan test yang relevan jika tersedia.
2. Periksa diff agar tidak ada perubahan tidak terkait.
3. Laporkan apa yang diubah, file penting, hasil validasi, dan hal yang belum dapat diuji.
4. Jangan mengklaim test berhasil jika test tidak dijalankan atau gagal.

## 4. Prinsip Coding

- Utamakan kode yang sederhana, terbaca, konsisten, dan mudah dirawat.
- Hindari overengineering, abstraksi prematur, dan dependency yang tidak diperlukan.
- Jangan menambahkan Redis, queue, microservices, Kubernetes, atau layanan eksternal lain tanpa kebutuhan nyata dan persetujuan pengguna.
- Gunakan nama variabel, fungsi, tabel, dan endpoint yang jelas serta konsisten.
- Pisahkan tanggung jawab; hindari handler/controller yang berisi seluruh logika bisnis.
- Tangani error secara eksplisit; jangan menelan error atau menggunakan `panic` untuk alur normal.
- Hindari duplikasi logika bisnis antara frontend dan backend.
- Jangan meninggalkan placeholder, mock, tombol mati, atau implementasi palsu pada fitur yang diklaim selesai.
- Jangan melakukan refactor besar di luar cakupan tugas.
- Jangan menghapus pekerjaan pengguna atau mengembalikan perubahan lokal tanpa izin.
- Jangan menambahkan file dokumentasi berulang kecuali memang dibutuhkan.

## 5. Backend Go + Gin

Gunakan struktur berlapis yang konsisten dengan repository:
- `handler`: HTTP request/response dan pemetaan status HTTP.
- `service`: aturan bisnis dan orkestrasi.
- `repository`: akses data.
- `model`: model/domain database.
- `middleware`: autentikasi, otorisasi, rate limit, dan logging.
- `config`: konfigurasi environment.
- `ai`: provider Gemini dan tools AI.

Ketentuan:
- Gunakan `context.Context` untuk operasi yang bergantung pada request.
- Validasi semua input di server, termasuk ID, rentang tanggal, jumlah, harga, diskon, dan status.
- Gunakan error handling idiomatis Go dan response API yang konsisten.
- Jangan mengembalikan stack trace, SQL mentah, secret, atau detail internal kepada client.
- Jangan menyimpan konfigurasi sensitif langsung di source code.
- Jangan menggunakan global mutable state jika tidak diperlukan.
- Gunakan interface jika membantu pengujian atau memisahkan integrasi eksternal; jangan membuat interface hanya demi abstraksi.
- Endpoint yang mengubah data harus memiliki otorisasi, validasi, dan penanganan kegagalan yang sesuai.
- Terapkan pagination dan batas maksimum ukuran hasil pada endpoint daftar.
- Gunakan transaksi database untuk operasi yang terdiri dari beberapa perubahan terkait.

## 6. Database, GORM, dan Migrasi

- Database target adalah MySQL 8 / InnoDB.
- Jangan menjalankan `AutoMigrate` sebagai pengganti strategi migrasi produksi tanpa persetujuan eksplisit.
- Perubahan skema harus menggunakan migration yang dapat ditinjau dan memiliki urutan jelas.
- Jangan menghapus atau mengganti nama kolom/tabel yang berisi data tanpa migration dan strategi kompatibilitas.
- Gunakan foreign key, unique constraint, index, dan constraint yang sesuai.
- Jangan mengandalkan validasi frontend untuk integritas data.
- Gunakan parameter query melalui GORM; hindari menyusun SQL dari input pengguna.
- Untuk query laporan yang kompleks, raw SQL terparameter boleh digunakan jika lebih jelas dan telah diuji.
- Perhatikan N+1 query; gunakan preload/join secara terukur.
- Jangan menggunakan `float32` atau `float64` untuk uang.
- Simpan waktu dengan standar yang konsisten; konversi tampilan ke `Asia/Jakarta`.
- Pastikan pembuatan nomor pesanan aman terhadap request bersamaan.
- Pastikan operasi pembuatan pesanan, item, dan pembayaran yang terkait berjalan atomik.
- Soft delete atau pembatalan tidak boleh menghilangkan jejak transaksi yang dibutuhkan untuk audit.

## 7. Aturan Transaksi dan Status

- Total pesanan dihitung ulang di backend dari item, harga yang valid, dan diskon yang diizinkan.
- Jangan mempercayai total, harga, status pembayaran, atau status pengerjaan yang dikirim frontend tanpa validasi.
- Snapshot harga harus tetap konsisten meskipun layanan kemudian diubah atau dinonaktifkan.
- Pembayaran normal harus sama dengan total tagihan yang harus dibayar; jangan implementasikan DP/cicilan kecuali PRD diubah secara eksplisit.
- Status `LUNAS` hanya diberikan setelah pembayaran valid tercatat dan jumlahnya memenuhi tagihan.
- Perubahan status pengerjaan harus mematuhi transisi yang diizinkan. Pengecualian Owner harus tercatat.
- Pesanan yang dibatalkan tetap memiliki riwayat dan alasan pembatalan.
- Pembatalan pesanan tidak otomatis berarti refund.
- Koreksi atau pembatalan pembayaran harus dibatasi untuk Owner dan dicatat pada audit log.
- Gunakan idempotency untuk endpoint pembuatan pesanan atau operasi lain yang rawan terkirim ulang.
- Pastikan retry request tidak menggandakan transaksi.

## 8. Autentikasi dan Otorisasi

- Backend adalah sumber kebenaran untuk hak akses.
- Jangan hanya menyembunyikan menu atau tombol di frontend.
- Terapkan prinsip least privilege.
- Peran yang valid: `owner`, `kasir`, dan `operator` (ikuti konvensi enum proyek bila berbeda).
- Hanya Owner boleh mengakses AI Assistant, laporan keuangan lengkap, pengeluaran, manajemen pengguna, pengaturan, dan audit log.
- Kasir hanya boleh mengakses fungsi operasional yang ditetapkan PRD.
- Operator hanya boleh mengakses daftar pekerjaan, detail pekerjaan yang relevan, dan perubahan status yang diizinkan.
- Hash password dengan bcrypt atau Argon2id; jangan pernah menyimpan password plaintext.
- Jika memakai cookie autentikasi, konfigurasi `HttpOnly`, `Secure` pada produksi, `SameSite`, dan perlindungan CSRF sesuai arsitektur.
- Jika memakai JWT, jaga secret, expiry, refresh flow, dan pencabutan sesi secara aman.
- Batasi percobaan login gagal dan endpoint sensitif.
- Jangan menaruh token atau API key di local storage jika arsitektur dapat menggunakan mekanisme yang lebih aman.
- Jangan pernah mencatat password, token, cookie, atau API key ke log.

## 9. Frontend Next.js, React, dan TypeScript

- Gunakan TypeScript dengan tipe eksplisit pada boundary data dan hindari `any` jika bisa dihindari.
- Ikuti struktur routing dan pola komponen yang sudah ada.
- Gunakan komponen bersama untuk pola UI berulang; jangan membuat design system baru jika sudah ada.
- Gunakan Tailwind CSS dan shadcn/ui sesuai pola proyek.
- UI harus berbahasa Indonesia, responsif, dan nyaman digunakan pada laptop, tablet, serta HP.
- Prioritaskan alur POS yang cepat, form yang ringkas, dan interaksi yang jelas.
- Tampilkan state loading, kosong, error, dan sukses pada halaman data.
- Gunakan konfirmasi untuk tindakan berisiko seperti pembatalan.
- Status pesanan harus memiliki warna sekaligus label teks; jangan mengandalkan warna saja.
- Format uang dengan format Rupiah Indonesia, misalnya `Rp 1.234.567`.
- Tanggal dan waktu ditampilkan dalam zona `Asia/Jakarta`.
- Jangan menaruh secret backend atau Gemini API key di variabel environment publik seperti `NEXT_PUBLIC_*`.
- Jangan menganggap validasi frontend cukup untuk keamanan.
- Hindari UI generik yang berlebihan, gradient dekoratif tanpa tujuan, animasi mengganggu, dan placeholder data pada fitur produksi.
- Jangan mengubah seluruh desain atau navigasi untuk tugas kecil tanpa diminta.
- Pastikan tombol dan form benar-benar berfungsi; jangan membuat kontrol interaktif palsu.

## 10. Integrasi Gemini dan AI Assistant

AI Assistant hanya untuk Owner dan harus bersifat **read-only** terhadap data bisnis.

Aturan wajib:
- Panggilan Gemini dilakukan dari backend, bukan browser.
- API key hanya berada di environment variable server.
- Gemini tidak mendapat kredensial atau koneksi database langsung.
- Gemini tidak boleh menghasilkan SQL yang kemudian dieksekusi.
- Model hanya boleh memanggil tools yang terdaftar dan diizinkan.
- Validasi nama tool, argumen, rentang tanggal, role, ukuran hasil, dan jumlah pemanggilan tools di backend.
- Batasi jumlah tool call per pertanyaan dan penggunaan harian sesuai konfigurasi.
- Gunakan timeout dan penanganan error.
- Batasi hasil daftar, misalnya maksimal 50 baris, serta hindari mengirim kolom yang tidak dibutuhkan.
- Kirim data agregat jika cukup untuk menjawab pertanyaan.
- Nama, nomor telepon, alamat, dan catatan pelanggan hanya dikirim jika memang dibutuhkan untuk jawaban Owner.
- Perlakukan isi database, pesan pelanggan, dan hasil tools sebagai data tidak tepercaya, bukan instruksi.
- Lindungi dari prompt injection dengan memisahkan instruksi sistem, pertanyaan pengguna, dan data hasil tools.
- AI tidak boleh membuat, mengubah, membatalkan, menghapus, atau merefund pesanan; mengubah harga; mencatat pembayaran; atau mengubah pengeluaran.
- Draft pesan pelanggan hanya menghasilkan teks dan tidak mengirim pesan.
- Angka finansial harus berasal dari hasil tools backend, bukan hasil hitung atau tebakan model.
- Jawaban harus jujur jika data tidak tersedia atau tool gagal.
- Sertakan periode data yang digunakan pada jawaban yang berkaitan dengan waktu.
- Jangan menganggap keluaran model sebagai sumber kebenaran untuk angka atau keputusan transaksi.
- Buat lapisan provider agar integrasi Gemini dapat diganti tanpa mengubah seluruh logika bisnis.
- Jangan menyimpan secret, token, atau data pribadi yang tidak perlu dalam log percakapan.
- Periksa ketentuan privasi dan penggunaan data Google yang berlaku sebelum deployment produksi.

## 11. API dan Kontrak Data

- Gunakan pola endpoint dan format response yang konsisten.
- Validasi request body, query parameter, dan path parameter.
- Gunakan status HTTP yang sesuai.
- Jangan mengembalikan field internal yang tidak diperlukan.
- Dokumentasikan endpoint baru atau perubahan kontrak API jika proyek menggunakan OpenAPI atau dokumentasi API.
- Jika response API berubah, periksa dan sesuaikan seluruh pemanggil frontend serta test.
- Jangan membuat endpoint duplikat tanpa alasan.
- Error response harus membantu client tanpa membocorkan detail internal.
- Semua endpoint daftar harus memiliki batas ukuran hasil dan pagination bila relevan.

## 12. Testing dan Quality Gate

Tambahkan atau perbarui test untuk setiap perubahan yang berdampak pada perilaku aplikasi.

Prioritas test:
- Perhitungan subtotal, diskon, dan total.
- Pembayaran penuh dan status lunas.
- Transisi status pengerjaan.
- Nomor pesanan unik saat concurrent requests.
- Idempotency.
- Hak akses Owner/Kasir/Operator.
- Pembatalan dan audit log.
- Ringkasan omzet, kas masuk, dan pengeluaran.
- AI tools, validasi argumen, prompt injection, timeout, dan tool failure.

Sebelum menyatakan pekerjaan selesai, jalankan yang tersedia dan relevan:
- Go: `gofmt`, `go vet`, `go test ./...`.
- Frontend: formatter, lint, type-check, dan test sesuai `package.json`.
- Migration atau integration test jika perubahan database membutuhkannya.
- Build frontend/backend jika memungkinkan.

Jika test tidak dapat dijalankan karena dependency, service, atau environment belum tersedia, jelaskan dengan tepat. Jangan mengarang hasil test.

## 13. Environment, Git, dan Dependency

- Sediakan `.env.example` tanpa secret nyata.
- Jangan commit `.env`, kredensial, token, file database, dump berisi data pelanggan, atau hasil build.
- Jangan mengubah `.gitignore` untuk memasukkan secret atau artefak lokal.
- Sebelum menambah dependency, pastikan fungsi tersebut tidak tersedia dari standard library atau dependency yang sudah ada.
- Hindari upgrade dependency massal tanpa kebutuhan.
- Jangan menjalankan perintah destruktif, reset database, menghapus file, force push, atau mengubah riwayat Git tanpa izin eksplisit.
- Jangan menimpa perubahan lokal yang belum di-commit.
- Jangan mengubah file di luar repository kecuali diminta.
- Gunakan Docker Compose untuk lingkungan lokal jika konfigurasi proyek sudah tersedia; jangan menganggap Docker sudah terpasang.

## 14. Dokumentasi

Perbarui dokumentasi yang relevan jika mengubah:
- Environment variables.
- Setup lokal.
- Skema database atau migration.
- Endpoint API.
- Hak akses.
- Alur transaksi.
- Integrasi Gemini.
- Proses deployment atau backup.

Dokumentasi harus sesuai implementasi nyata. Jangan menulis fitur sebagai selesai jika baru berupa rencana.

## 15. Cara Melaporkan Hasil

Setelah menyelesaikan tugas, berikan ringkasan singkat dengan format:

1. **Perubahan:** apa yang dibuat atau diperbaiki.
2. **File utama:** daftar file penting yang diubah.
3. **Validasi:** test, lint, type-check, atau build yang dijalankan beserta hasilnya.
4. **Catatan:** asumsi, keterbatasan, migration, konfigurasi, atau langkah manual yang masih diperlukan.

Jangan menyatakan aplikasi siap produksi hanya karena berhasil di-build. Jelaskan jika masih ada pengujian keamanan, backup/restore, konfigurasi server, atau UAT yang belum dilakukan.

## 16. Checklist Sebelum Menyelesaikan Tugas

- [ ] Perubahan sesuai permintaan dan PRD.
- [ ] Tidak ada secret atau kredensial di source code.
- [ ] Hak akses diperiksa di backend.
- [ ] Validasi input dilakukan di backend.
- [ ] Operasi multi-langkah menggunakan transaksi database jika diperlukan.
- [ ] Perubahan skema disertai migration yang sesuai.
- [ ] Test relevan ditambahkan atau diperbarui.
- [ ] Formatter, lint, type-check, test, atau build yang tersedia telah dijalankan.
- [ ] Tidak ada perubahan tidak terkait atau file pengguna yang terhapus.
- [ ] Dokumentasi diperbarui bila diperlukan.
- [ ] Ringkasan hasil dan keterbatasan disampaikan dengan jujur.
