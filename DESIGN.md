# DESIGN.md — SatSet Design System & UI/UX Guidelines

> Dokumen ini adalah pedoman visual dan pengalaman pengguna untuk seluruh antarmuka SatSet. Baca bersama `AGENTS.md` dan PRD sebelum mengerjakan UI. Jika ada konflik, aturan bisnis dan keamanan di PRD tetap menjadi acuan; dokumen ini mengatur cara fitur ditampilkan dan digunakan.

## 1. Product Design Direction

SatSet adalah aplikasi operasional laundry yang dipakai sehari-hari oleh Owner, Kasir, dan Operator. Antarmuka harus terasa seperti produk SaaS yang matang: ramah, rapi, cepat dipahami, dan dapat diandalkan.

### Kata kunci desain
- **Friendly:** bahasa sederhana, istilah operasional yang akrab bagi pengguna Indonesia.
- **Clean:** visual bersih, hierarki informasi jelas, dekorasi seperlunya.
- **Efficient:** tugas yang sering dilakukan dapat diselesaikan dengan sedikit langkah.
- **Professional:** konsisten, presisi, tidak terlihat seperti template dashboard generik.
- **Accessible:** kontras cukup, ukuran teks nyaman, status tidak hanya dibedakan dengan warna.

Jangan mengorbankan kejelasan dan kecepatan demi dekorasi. Antarmuka harus tetap terlihat bagus saat berisi data nyata, tabel panjang, nama pelanggan panjang, dan status yang beragam.

## 2. Visual Identity

### 2.1 Font

Gunakan **Poppins** sebagai font utama di seluruh aplikasi.

- Muat Poppins dari sumber font yang disetujui proyek (misalnya Google Fonts atau paket font lokal jika diperlukan).
- Gunakan bobot 400, 500, 600, dan 700 saja kecuali ada alasan kuat.
- Body text: 14–16 px.
- Label/helper text: minimal 12 px; jangan gunakan teks kecil untuk informasi penting.
- Judul halaman: 24–30 px desktop, 22–26 px layar kecil.
- Angka metrik utama: 24–32 px sesuai konteks.
- Jangan menggunakan banyak bobot font sekaligus dalam satu area.
- Pastikan ada fallback: `Poppins, ui-sans-serif, system-ui, sans-serif`.

Contoh CSS:
```css
:root {
  font-family: "Poppins", ui-sans-serif, system-ui, sans-serif;
}
```

### 2.2 Color palette

Gunakan palet warna yang konsisten dan tenang. Warna utama adalah teal/biru kehijauan untuk identitas SatSet; hindari menjadikan seluruh halaman penuh warna.

| Token | Nilai | Penggunaan |
|---|---|---|
| Primary | `#0F766E` | Tombol utama, link penting, elemen aktif |
| Primary hover | `#115E59` | Hover tombol utama |
| Primary soft | `#CCFBF1` | Highlight ringan, icon background terpilih |
| App background | `#F8FAFC` | Latar utama halaman |
| Surface | `#FFFFFF` | Panel, form, tabel, dialog |
| Text primary | `#0F172A` | Heading dan teks utama |
| Text secondary | `#475569` | Teks sekunder |
| Text muted | `#64748B` | Helper text, metadata |
| Border | `#E2E8F0` | Border input, divider, tabel |
| Success | `#15803D` | Berhasil, selesai, pembayaran lunas |
| Success soft | `#DCFCE7` | Latar status sukses |
| Warning | `#B45309` | Perlu perhatian, proses tertunda |
| Warning soft | `#FEF3C7` | Latar status peringatan |
| Danger | `#B91C1C` | Error, pembatalan, tindakan destruktif |
| Danger soft | `#FEE2E2` | Latar status error |
| Info | `#1D4ED8` | Informasi dan status informatif |
| Info soft | `#DBEAFE` | Latar status informasi |

Ketentuan:
- Warna utama digunakan secara selektif untuk aksi utama dan navigasi aktif.
- Hindari gradient dekoratif besar, neon, warna pastel berlebihan, dan latar berwarna pada setiap card.
- Jangan menggunakan warna status sebagai satu-satunya indikator. Sertakan teks status dan, jika membantu, ikon.
- Pastikan rasio kontras teks memadai; targetkan WCAG AA.

### 2.3 Shape, borders, shadows

- Gunakan radius secara konsisten: input/button 8–10 px; panel/card 12–16 px; dialog 16 px.
- Border tipis `1px solid #E2E8F0` lebih diutamakan daripada shadow tebal.
- Shadow hanya untuk elemen yang benar-benar mengambang seperti dropdown, popover, dan dialog.
- Hindari setiap elemen dibungkus card. Gunakan spacing, heading, dan divider untuk mengelompokkan konten.
- Jangan menggunakan efek glassmorphism, glow, gradient berlebihan, atau animasi dekoratif yang tidak membantu tugas.

### 2.4 Spacing and sizing

Gunakan skala spacing yang konsisten berbasis kelipatan 4 px:
- 4 px: jarak mikro.
- 8 px: jarak elemen yang berkaitan.
- 12–16 px: padding input dan jarak komponen.
- 20–24 px: jarak antarbagian.
- 24–32 px: padding halaman atau panel besar.
- 32–40 px: pemisah kelompok konten utama.

Usahakan halaman terasa lega, tetapi jangan membuang ruang terlalu banyak sehingga data penting harus sering di-scroll.

## 3. Application Layout

### 3.1 Main shell

Gunakan pola aplikasi dashboard yang konsisten:
- Sidebar kiri pada desktop, dengan logo/nama SatSet, navigasi berikon dan label, serta penanda halaman aktif yang jelas.
- Area konten utama memiliki header halaman, breadcrumb bila berguna, judul, deskripsi singkat bila perlu, dan aksi utama.
- Informasi pengguna/akun ditempatkan konsisten di area header atau sidebar.
- Pada tablet/mobile, sidebar berubah menjadi drawer atau navigasi yang sesuai ukuran layar.
- Jangan membuat sidebar terlalu lebar; target sekitar 224–256 px pada desktop.
- Sidebar tidak boleh mengambil terlalu banyak perhatian dari konten.

### 3.2 Role-aware navigation

Tampilkan hanya menu yang relevan untuk role pengguna yang sedang login. Ini membantu mengurangi kebingungan, tetapi **bukan pengganti authorization di backend**.

- **Owner:** Dashboard, Pesanan, Pelanggan, Layanan, Pengeluaran, Laporan, AI Assistant, Pengaturan.
- **Kasir:** Dashboard ringkas, Pesanan/POS, Pelanggan, pembayaran dan cetak struk sesuai izin.
- **Operator:** Daftar pekerjaan/pesanan, detail pesanan yang diperlukan, dan pembaruan status sesuai izin.
- Jangan tampilkan menu yang tidak diizinkan hanya untuk kemudian mengandalkan error ketika diklik.

### 3.3 Page header

Setiap halaman sebaiknya memiliki:
1. Judul yang jelas.
2. Deskripsi satu kalimat hanya jika membantu.
3. Aksi utama yang paling sering digunakan, ditempatkan konsisten di kanan atas desktop dan mudah ditemukan di mobile.
4. Filter atau rentang tanggal di dekat data yang difilter.

Jangan menampilkan terlalu banyak tombol primer sekaligus. Gunakan satu aksi utama, lalu aksi sekunder sebagai tombol outline, ghost, atau menu.

## 4. Component Guidelines

### 4.1 Buttons

- Primary: latar Primary, teks putih, digunakan untuk aksi utama.
- Secondary: latar putih/neutral dengan border halus.
- Ghost: untuk aksi ringan seperti navigasi atau aksi tabel yang tidak dominan.
- Destructive: warna merah hanya untuk tindakan berisiko seperti pembatalan atau penghapusan.
- Tinggi tombol umum 40–44 px; tombol kecil tidak boleh terlalu sempit untuk disentuh.
- Label harus berupa kata kerja yang jelas: `Buat Pesanan`, `Simpan Perubahan`, `Proses Pembayaran`.
- Tampilkan loading state dan nonaktifkan tombol selama permintaan sedang diproses agar tidak terjadi submit ganda.
- Untuk tindakan destruktif atau sulit dibatalkan, minta konfirmasi dan jelaskan dampaknya.

### 4.2 Forms and inputs

- Label selalu terlihat; jangan hanya mengandalkan placeholder.
- Gunakan pesan validasi dekat dengan field terkait, dengan bahasa yang spesifik dan membantu.
- Bedakan state default, hover, focus, error, disabled, dan read-only.
- Fokus keyboard harus terlihat jelas.
- Gunakan input yang sesuai: angka untuk nominal/berat, tanggal untuk tanggal, dan pencarian untuk pilihan dengan banyak data.
- Tampilkan format rupiah dan satuan secara konsisten.
- Hindari form panjang satu kolom jika dapat dikelompokkan secara logis. Untuk form panjang, gunakan section dengan heading.
- Jangan menghapus input pengguna saat validasi gagal.
- Konfirmasi sebelum meninggalkan form yang memiliki perubahan belum disimpan bila risikonya relevan.

### 4.3 Cards and metrics

- Gunakan card untuk metrik penting atau konten yang memang perlu dikelompokkan.
- Dashboard cukup menampilkan metrik prioritas, misalnya omzet pada periode terpilih, jumlah pesanan, pesanan dalam proses, dan pesanan siap diambil.
- Setiap metrik memiliki label, angka, dan konteks periode yang jelas.
- Hindari terlalu banyak card, ikon besar, dekorasi, dan grafik yang tidak membantu keputusan.
- Jangan membuat angka dummy atau tren palsu pada production UI.

### 4.4 Tables and lists

- Tabel harus mudah dipindai: header jelas, alignment angka konsisten, row height nyaman, dan status mudah ditemukan.
- Nominal uang rata kanan dan menggunakan format Rupiah.
- Gunakan pagination atau mekanisme pemuatan data yang sesuai; jangan memuat data tak terbatas.
- Search, filter, dan sort ditempatkan konsisten di atas tabel.
- Kolom tindakan jangan terlalu ramai. Gunakan menu tindakan jika diperlukan.
- Di layar kecil, tabel dapat berubah menjadi list/card ringkas atau memakai horizontal scrolling yang disengaja. Jangan membuat seluruh halaman overflow horizontal.
- Sediakan empty state yang menjelaskan apa yang belum ada dan tindakan berikutnya.

### 4.5 Status badges

Gunakan badge dengan warna lembut dan label eksplisit. Contoh status pesanan:
- Baru diterima: neutral/info.
- Diproses: info.
- Siap diambil: warning/teal sesuai sistem.
- Selesai: success.
- Dibatalkan: danger.

Pastikan pemetaan warna sama di seluruh halaman. Jangan mengubah arti warna antarfitur.

### 4.6 Dialogs, drawers, and toasts

- Dialog dipakai untuk konfirmasi singkat atau tugas kecil, bukan sebagai wadah form yang terlalu panjang.
- Drawer cocok untuk melihat detail atau mengedit informasi pendukung tanpa kehilangan konteks halaman.
- Toast untuk feedback singkat seperti berhasil menyimpan. Error penting harus tetap dapat ditemukan dan dipahami, bukan hanya hilang sebagai toast.
- Setelah aksi berhasil, perbarui data yang terlihat dan tampilkan feedback yang relevan.
- Modal harus bisa ditutup dengan jelas, menjaga fokus keyboard, dan tidak menutup secara tidak sengaja ketika pengguna sedang mengisi form penting.

### 4.7 Icons

- Gunakan satu library ikon yang konsisten, misalnya Lucide.
- Ikon harus membantu mengenali fungsi, bukan menjadi dekorasi.
- Aksi penting sebaiknya memiliki teks; ikon tanpa label perlu tooltip atau accessible label.
- Jangan mencampur banyak gaya ikon.

## 5. Core User Flows

### 5.1 Kasir — buat pesanan

Buat alur transaksi secepat dan sejelas mungkin:
1. Cari pelanggan atau buat pelanggan baru tanpa kehilangan data pesanan yang sudah diisi.
2. Pilih layanan dan masukkan berat/jumlah sesuai jenis layanan.
3. Tampilkan rincian harga, subtotal, dan total secara langsung.
4. Tampilkan ringkasan sebelum pesanan disimpan.
5. Simpan pesanan dan tampilkan nomor pesanan yang mudah ditemukan.
6. Lanjutkan ke pembayaran sesuai aturan bisnis, kemudian tawarkan cetak struk.

Ketentuan:
- Total dihitung oleh backend sebagai sumber kebenaran; frontend hanya memberi preview.
- Semua label, satuan, dan nominal harus jelas.
- Hindari memaksa kasir berpindah banyak halaman untuk menyelesaikan transaksi.
- Jika terjadi error, pertahankan data yang telah dimasukkan bila aman.
- Jangan menampilkan opsi uang muka atau cicilan karena aturan bisnis SatSet mewajibkan pembayaran penuh.

### 5.2 Pembayaran

- Tampilkan nomor pesanan, nama pelanggan, total tagihan, jumlah yang dibayar, dan status pembayaran dalam satu ringkasan.
- Jelaskan secara eksplisit bahwa pembayaran dilakukan penuh.
- Cegah submit ganda dan tampilkan status proses.
- Setelah berhasil, tampilkan konfirmasi serta tindakan lanjutan seperti cetak struk.
- Jangan menganggap pembayaran berhasil hanya karena tombol diklik; gunakan respons backend.

### 5.3 Operator — pembaruan status

- Prioritaskan daftar pesanan yang perlu dikerjakan sekarang.
- Buat status dan langkah berikutnya mudah dipahami.
- Sediakan tombol aksi yang jelas sesuai status yang diizinkan.
- Konfirmasi bila aksi berpotensi sulit dibatalkan.
- Jangan tampilkan data finansial atau fungsi administratif yang tidak diperlukan role Operator.

### 5.4 Owner — dashboard dan laporan

- Mulai dengan ringkasan yang menjawab: bagaimana kondisi usaha, pesanan mana yang perlu perhatian, dan bagaimana kinerja pada periode yang dipilih.
- Filter tanggal harus terlihat dan mudah diubah.
- Grafik wajib memiliki judul, label, dan konteks periode.
- Sediakan tampilan tabel atau rincian ketika pengguna perlu memeriksa angka di balik grafik.
- Hindari chart hanya untuk mempercantik tampilan.

### 5.5 AI Assistant

- Bedakan area percakapan dari sumber data/insight yang dikembalikan.
- Tampilkan loading, error, dan status kosong secara jelas.
- Jawaban AI harus mudah dibaca: ringkasan dahulu, detail setelahnya, gunakan tabel/list bila membantu.
- Jika insight didasarkan pada rentang tanggal atau data tertentu, tampilkan konteks tersebut.
- AI hanya untuk membaca dan menganalisis data sesuai izin. Jangan sediakan kontrol UI yang menyiratkan AI dapat mengubah pesanan, pembayaran, layanan, atau database.
- Beri penjelasan ketika data tidak tersedia atau hasil tidak pasti; jangan membuat angka atau transaksi fiktif.

## 6. Responsive Behavior

Rancang dari kebutuhan desktop operasional, lalu pastikan pengalaman tablet dan mobile tetap baik.

### Desktop (≥ 1024 px)
- Sidebar tetap terlihat.
- Dashboard dapat memakai grid yang rapi.
- Tabel menampilkan kolom penting tanpa membuat teks terlalu rapat.

### Tablet (768–1023 px)
- Sidebar dapat diciutkan atau menjadi drawer.
- Kurangi kolom yang tidak esensial.
- Form dua kolom hanya jika tiap field masih cukup lebar.

### Mobile (< 768 px)
- Navigasi menjadi drawer/menu yang jelas.
- Form umumnya satu kolom.
- Aksi utama mudah dijangkau dan tidak tertutup elemen lain.
- Tabel diubah menjadi list/card ringkas atau scroll horizontal pada area tabel saja.
- Hindari tombol terlalu kecil, dialog yang lebih lebar daripada layar, dan teks yang terpotong tanpa cara melihatnya.

Uji setidaknya pada lebar 360 px, 390 px, 768 px, 1024 px, dan 1440 px.

## 7. UX States: Required on Every Data-Driven Screen

Setiap halaman yang mengambil atau mengubah data harus memiliki state yang dirancang, bukan hanya happy path:

- **Loading:** skeleton atau indikator yang sesuai; hindari spinner besar untuk seluruh halaman jika hanya satu bagian yang memuat.
- **Empty:** jelaskan mengapa daftar kosong dan apa yang bisa dilakukan berikutnya.
- **Error:** pesan manusiawi, tindakan retry bila sesuai, jangan bocorkan stack trace atau detail internal.
- **Success:** konfirmasi singkat dan pembaruan data.
- **Disabled:** jelaskan atau buat alasan status disabled dapat dipahami.
- **Permission denied:** pesan yang sopan tanpa membocorkan informasi yang tidak boleh diakses.
- **Unsaved changes:** peringatan ketika perubahan penting akan hilang.

Jangan menggunakan placeholder data seolah-olah itu data nyata. Data contoh hanya boleh digunakan pada seed/demo yang memang ditujukan untuk pengembangan.

## 8. Accessibility and Interaction

- Gunakan elemen HTML semantik dan struktur heading berurutan.
- Semua fungsi utama dapat digunakan dengan keyboard.
- Focus ring terlihat; jangan menghilangkan outline tanpa pengganti.
- Setiap input memiliki label yang terhubung.
- Ikon dekoratif disembunyikan dari screen reader; ikon fungsional memiliki accessible name.
- Jangan mengandalkan warna saja untuk menunjukkan status/error.
- Hormati `prefers-reduced-motion`.
- Animasi transisi singkat dan fungsional, umumnya 150–200 ms; tidak ada animasi yang menghambat pekerjaan.
- Area klik target minimal sekitar 40 × 40 px jika memungkinkan.

## 9. Content and Microcopy

Gunakan Bahasa Indonesia yang natural, singkat, konsisten, dan berorientasi tindakan.

Contoh:
- `Buat Pesanan`, bukan `Create New Order`.
- `Belum ada pesanan`, diikuti arahan tindakan yang sesuai.
- `Periksa koneksi Anda lalu coba lagi.`, bukan `Something went wrong`.
- `Pembayaran berhasil`, hanya setelah backend mengonfirmasi.
- `Hapus layanan ini? Tindakan ini tidak dapat dibatalkan.`, untuk konfirmasi penghapusan.

Ketentuan:
- Hindari istilah teknis untuk pengguna operasional.
- Hindari sapaan berlebihan, copywriting yang terlalu ramai, dan teks placeholder seperti `Lorem ipsum`.
- Format mata uang konsisten dengan `Rp` dan pemisah ribuan Indonesia.
- Tanggal/waktu ditampilkan dalam format Indonesia yang mudah dibaca dan zona waktu aplikasi yang konsisten.
- Gunakan istilah bisnis yang sama untuk satu konsep di seluruh aplikasi.

## 10. Charts and Data Visualization

- Gunakan chart hanya jika membantu membandingkan tren atau distribusi.
- Hindari efek 3D, gradient, animasi berlebihan, dan terlalu banyak warna seri.
- Selalu tampilkan judul, rentang waktu, label sumbu yang relevan, tooltip yang berguna, serta empty state.
- Pastikan data dapat dipahami tanpa mengandalkan warna saja.
- Sertakan angka ringkasan atau tabel jika nilai tepat lebih penting daripada tren visual.
- Jangan menyajikan angka buatan, prediksi, atau persentase tanpa sumber/perhitungan yang jelas.

## 11. Technical Implementation Guidance

Ikuti teknologi yang telah dipilih di PRD dan `AGENTS.md`. Jangan mengganti stack atau menambah dependency besar hanya untuk meniru desain.

Jika proyek menggunakan Tailwind dan shadcn/ui:
- Gunakan design token/CSS variables untuk warna, radius, border, dan focus ring.
- Gunakan komponen dasar yang konsisten dan komposisikan ulang untuk kebutuhan khusus.
- Jangan menyalin komponen yang sama ke banyak tempat jika dapat dibuat reusable.
- Hindari inline style acak untuk nilai yang seharusnya menjadi token.
- Pastikan font Poppins dimuat dengan benar dan tidak menyebabkan flash/layout shift yang mengganggu.
- Gunakan ikon dari satu library yang telah dipilih proyek.
- Jangan menginstal library UI tambahan tanpa kebutuhan nyata dan persetujuan sesuai `AGENTS.md`.

Sebelum membuat komponen baru, periksa apakah proyek sudah memiliki komponen serupa. Pertahankan pola, struktur folder, dan conventions repo yang sudah ada.

## 12. Anti-Patterns — Do Not Do These

Jangan:
- Membuat UI seperti template dashboard AI yang penuh gradient ungu/biru.
- Menggunakan terlalu banyak card dengan shadow besar.
- Menaruh emoji sebagai ikon navigasi utama.
- Membuat semua teks terlalu kecil atau abu-abu pucat.
- Menggunakan terlalu banyak warna aksen dalam satu halaman.
- Menampilkan chart dummy, angka palsu, atau aktivitas palsu.
- Menggunakan modal untuk setiap interaksi kecil.
- Menyembunyikan label field hanya dengan placeholder.
- Membuat semua tombol memiliki warna dan bobot visual yang sama.
- Menyisakan halaman kosong tanpa empty state yang informatif.
- Mengandalkan hover sebagai satu-satunya cara mengakses aksi.
- Mengabaikan loading, error, disabled, dan permission state.
- Membuat halaman hanya bagus di desktop dan rusak di mobile.
- Menggunakan font selain Poppins secara tidak konsisten.
- Menambahkan animasi, efek blur, glassmorphism, atau dekorasi tanpa manfaat UX.

## 13. Definition of Done for UI Work

Sebelum menyatakan pekerjaan UI selesai, pastikan:
- [ ] Mengikuti `DESIGN.md`, `AGENTS.md`, dan kebutuhan terkait di PRD.
- [ ] Font Poppins digunakan konsisten.
- [ ] Warna, spacing, radius, tombol, input, tabel, dan status konsisten.
- [ ] Layout terlihat rapi dan proporsional pada ukuran layar yang relevan.
- [ ] Alur utama dapat diselesaikan tanpa langkah yang tidak perlu.
- [ ] Loading, empty, error, success, dan permission state yang relevan tersedia.
- [ ] Form memiliki label, validasi, dan focus state yang jelas.
- [ ] Tidak ada data dummy yang tampil seolah data production.
- [ ] Aksi dan data mengikuti role serta izin backend.
- [ ] Tidak ada overflow horizontal tak disengaja atau konten terpotong.
- [ ] Kontras dan keyboard navigation layak.
- [ ] Tes, lint, dan type-check yang tersedia dijalankan; masalah yang tersisa dilaporkan.
- [ ] Tidak ada dependency atau perubahan stack yang tidak diperlukan.

## 14. Guidance for the Coding Agent

Sebelum mengimplementasikan halaman:
1. Baca `AGENTS.md` dan PRD.
2. Periksa komponen, token, font, dan pola layout yang sudah ada.
3. Tentukan tujuan utama halaman dan aksi terpenting pengguna.
4. Rancang loading, empty, error, success, dan permission state sebelum menyelesaikan UI.
5. Implementasikan dengan komponen reusable dan token desain.
6. Uji responsive behavior dan interaksi nyata.
7. Jalankan pemeriksaan yang tersedia dan laporkan hasilnya.

Jika desain yang diminta pengguna bertentangan dengan pedoman ini, ikuti arahan eksplisit terbaru dari pengguna untuk aspek visual, tetapi jangan melanggar aturan bisnis, authorization, privasi, atau keamanan pada `AGENTS.md` dan PRD.
