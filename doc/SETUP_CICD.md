# 🤖 Panduan Setup CI/CD Otomatis (GitHub Actions ➔ VPS)

Dokumen ini menjelaskan cara menghubungkan **GitHub Actions** dengan **VPS kamu** agar setiap kali kamu menjalankan `git push origin main`, aplikasi di VPS otomatis ter-update tanpa perlu SSH manual.

---

## 🛠️ Ringkasan 4 Rahasia (GitHub Secrets) yang Dibutuhkan

| Nama Secret | Nilai | Keterangan |
|---|---|---|
| `VPS_HOST` | `20.41.112.76` | Alamat IP publik VPS kamu |
| `VPS_USER` | `afchaarel` | Username login SSH VPS kamu |
| `VPS_PORT` | `22` | Port SSH VPS |
| `VPS_SSH_KEY` | *(Private Key dari Langkah 1)* | Kunci otentikasi SSH rahasia |

---

## Langkah 1: Buat SSH Key Khusus GitHub Actions di VPS

Masuk ke terminal VPS kamu:
```bash
ssh afchaarel@20.41.112.76
```

Jalankan perintah ini untuk membuat pasangan kunci SSH baru khusus bot GitHub:
```bash
ssh-keygen -t ed25519 -C "github-actions-satset" -f ~/.ssh/github_actions_key -N ""
```

Lalu daftarkan kunci publiknya ke daftar kunci yang diizinkan masuk:
```bash
cat ~/.ssh/github_actions_key.pub >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
chmod 700 ~/.ssh
```

Sekarang, tampilkan **Private Key** yang akan dimasukkan ke GitHub:
```bash
cat ~/.ssh/github_actions_key
```

Salin **seluruh teks** yang muncul di layar, mulai dari:
```text
-----BEGIN OPENSSH PRIVATE KEY-----
...kode panjang...
-----END OPENSSH PRIVATE KEY-----
```

---

## Langkah 2: Daftarkan Secrets di Repository GitHub

1. Buka browser, buka repository kamu di GitHub:
   `https://github.com/arel10/SatSet-App`
2. Klik tab **Settings** (di sebelah kanan tab Pull Requests & Actions).
3. Di sidebar sebelah kiri, klik **Secrets and variables** ➔ pilih **Actions**.
4. Klik tombol hijau **New repository secret**.
5. Tambahkan 4 rahasia berikut satu per satu:

### Secret 1: `VPS_HOST`
* **Name**: `VPS_HOST`
* **Secret**: `20.41.112.76`
* Klik **Add secret**.

### Secret 2: `VPS_USER`
* **Name**: `VPS_USER`
* **Secret**: `afchaarel`
* Klik **Add secret**.

### Secret 3: `VPS_PORT`
* **Name**: `VPS_PORT`
* **Secret**: `22`
* Klik **Add secret**.

### Secret 4: `VPS_SSH_KEY`
* **Name**: `VPS_SSH_KEY`
* **Secret**: *(Tempel seluruh isi teks private key dari Langkah 1 tadi)*
* Klik **Add secret**.

---

## Langkah 3: Uji Coba CI/CD Otomatis

Setelah 4 Secrets terpasang:
1. Di laptop kamu, buat perubahan kecil atau langsung push file workflow ini:
   ```bash
   git add .
   git commit -m "ci: add GitHub Actions automated CI/CD pipeline"
   git push origin main
   ```
2. Buka tab **Actions** di GitHub (`https://github.com/arel10/SatSet-App/actions`).
3. Kamu akan melihat alur kerja **CI/CD Pipeline - SatSet Production** sedang berjalan:
   * 🟡 **Test & Build Verification (CI)**: Mengecek backend Go test & frontend Vite build.
   * 🚀 **Deploy to Production VPS (CD)**: Otomatis SSH ke VPS, jalankan `git pull`, dan `docker compose up -d --build`.
4. Setelah statusnya centang hijau (✅), buka web kamu di browser:
   `http://20.41.112.76:3000`
   Perubahan langsung live!
