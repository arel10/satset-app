# 🚀 Panduan Deploy SatSet ke VPS (Production)

Panduan lengkap step-by-step deploy aplikasi SatSet menggunakan Docker di VPS Ubuntu.

> [!IMPORTANT]
> VPS target: `20.41.112.76` (Ubuntu 22.04/24.04), SSH port 22, akses SatSet di port **3000**.
> Port 80 sudah dipakai app lain (cPanel/aaPanel), jadi SatSet pakai port 3000.

---

## Daftar Isi

1. [Persiapan Lokal](#1-persiapan-lokal)
2. [Install Docker di VPS](#2-install-docker-di-vps)
3. [Clone Repository di VPS](#3-clone-repository-di-vps)
4. [Konfigurasi Environment Production](#4-konfigurasi-environment-production)
5. [Build & Jalankan](#5-build--jalankan)
6. [Verifikasi & Akses Aplikasi](#6-verifikasi--akses-aplikasi)
7. [Maintenance & Operasional](#7-maintenance--operasional)
8. [Troubleshooting](#8-troubleshooting)

---

## 1. Persiapan Lokal

### 1.1. Commit & Push Kode ke GitHub/GitLab

Pastikan semua perubahan sudah di-commit dan di-push:

```bash
cd D:\Project\SatSet-App

# Cek status
git status

# Tambahkan file baru (docker-compose.prod.yml, .env.production.example, dll)
git add .

# Commit
git commit -m "feat: add production Docker setup for VPS deployment"

# Push ke remote
git push origin main
```

> [!WARNING]
> Jangan pernah commit file `.env`, `.env.production`, atau file yang berisi secret/API key.
> File `.env.production.example` (template tanpa secret) yang boleh di-commit.

---

## 2. Install Docker di VPS

### 2.1. SSH ke VPS

Dari terminal (PowerShell / Windows Terminal):

```bash
ssh username@20.41.112.76
```

Ganti `username` dengan username SSH kamu. Jika diminta password, masukkan password SSH kamu.

### 2.2. Update Sistem

```bash
sudo apt update && sudo apt upgrade -y
```

### 2.3. Install Dependensi

```bash
sudo apt install -y ca-certificates curl gnupg lsb-release git
```

### 2.4. Tambah Docker GPG Key & Repository

```bash
# Tambah Docker official GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Tambah Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

### 2.5. Install Docker Engine & Docker Compose

```bash
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

### 2.6. Tambah User ke Docker Group (Agar Tidak Perlu `sudo`)

```bash
sudo usermod -aG docker $USER

# Logout lalu login ulang agar berlaku
exit
```

Lalu SSH lagi:
```bash
ssh username@20.41.112.76
```

### 2.7. Verifikasi Instalasi

```bash
docker --version
# Docker version 28.x.x atau lebih baru

docker compose version
# Docker Compose version v2.x.x
```

---

## 3. Clone Repository di VPS

### 3.1. Clone dari GitHub/GitLab

```bash
# Buat folder project
mkdir -p ~/apps
cd ~/apps

# Clone repository (ganti URL sesuai repo kamu)
git clone https://github.com/USERNAME/SatSet-App.git
cd SatSet-App
```

> [!TIP]
> Jika repo private, kamu perlu setup SSH key atau personal access token di VPS.
> ```bash
> # Generate SSH key di VPS (jika belum ada)
> ssh-keygen -t ed25519 -C "email@kamu.com"
> cat ~/.ssh/id_ed25519.pub
> # Copy output, lalu tambahkan di GitHub → Settings → SSH Keys
> ```

---

## 4. Konfigurasi Environment Production

### 4.1. Buat File `.env.production`

```bash
# Salin template
cp .env.production.example .env.production

# Edit dengan nano
nano .env.production
```

### 4.2. Generate Secret yang Kuat

Jalankan perintah ini untuk menghasilkan secret JWT:

```bash
# Generate JWT Access Secret
openssl rand -hex 32
# Output contoh: a3f5b7c9d1e2f4a6b8c0d2e4f6a8b0c2d4e6f8a0b2c4d6e8f0a2b4c6d8e0f2

# Generate JWT Refresh Secret (jalankan lagi untuk secret berbeda)
openssl rand -hex 32

# Generate DB Password
openssl rand -hex 16

# Generate DB Root Password
openssl rand -hex 16
```

### 4.3. Isi `.env.production`

Edit file dan ganti semua nilai `GANTI_INI_...`:

```bash
nano .env.production
```

Yang **WAJIB** diganti:

| Variable | Keterangan |
|---|---|
| `DB_PASSWORD` | Password database user (hasil `openssl rand -hex 16`) |
| `DB_ROOT_PASSWORD` | Password root MySQL (hasil `openssl rand -hex 16`) |
| `JWT_ACCESS_SECRET` | Secret JWT access (hasil `openssl rand -hex 32`) |
| `JWT_REFRESH_SECRET` | Secret JWT refresh (hasil `openssl rand -hex 32`) |
| `FRONTEND_URL` | `http://20.41.112.76` (IP VPS kamu) |
| `GEMINI_API_KEY` | API key Gemini kamu (jika mau pakai AI) |

Contoh `.env.production` yang sudah diisi:

```env
APP_NAME=SatSet
APP_ENV=production
APP_TIMEZONE=Asia/Jakarta

DB_NAME=satset
DB_USER=satset
DB_PASSWORD=a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6
DB_ROOT_PASSWORD=f1e2d3c4b5a6f7e8d9c0b1a2f3e4d5c6

FRONTEND_PORT=3000

JWT_ACCESS_SECRET=a3f5b7c9d1e2f4a6b8c0d2e4f6a8b0c2d4e6f8a0b2c4d6e8f0a2b4c6d8e0f2
JWT_REFRESH_SECRET=b4f6c8d0e2a4b6c8d0e2f4a6b8c0d2e4f6a8b0c2d4e6f8a0b2c4d6e8f0a2b4

GEMINI_API_KEY=your_gemini_api_key_here
GEMINI_MODEL=gemini-3.5-flash-lite

AI_MAX_TOOL_CALLS=5
AI_DAILY_LIMIT_PER_USER=30
AI_REQUEST_TIMEOUT_SECONDS=30
AI_MAX_QUESTION_LENGTH=2000

LOGIN_MAX_ATTEMPTS=5
LOGIN_LOCKOUT_MINUTES=15

FRONTEND_URL=http://20.41.112.76:3000

AUTO_MIGRATE=true

VITE_APP_NAME="SatSet Laundry"
VITE_APP_DESCRIPTION="Sistem Manajemen Laundry & POS"
VITE_APP_TIMEZONE="Asia/Jakarta"
VITE_APP_DEFAULT_LOCALE="id-ID"
VITE_APP_CURRENCY="IDR"
```

Simpan: `Ctrl+O` → `Enter` → `Ctrl+X`

---

## 5. Build & Jalankan

### 5.1. Build dan Start Semua Service

```bash
cd ~/apps/SatSet-App

# Build image dan jalankan dalam background
docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build
```

> [!NOTE]
> Proses build pertama kali bisa memakan waktu **5-15 menit** tergantung kecepatan internet VPS
> karena harus download base image (golang, node, nginx, mysql) dan compile kode.

### 5.2. Pantau Progress Build

```bash
# Lihat log secara real-time
docker compose -f docker-compose.prod.yml --env-file .env.production logs -f

# Tekan Ctrl+C untuk keluar dari log (container tetap berjalan)
```

### 5.3. Cek Status Container

```bash
docker compose -f docker-compose.prod.yml --env-file .env.production ps
```

Output yang diharapkan (semua `healthy`):
```
NAME                   STATUS                    PORTS
satset-db-prod         Up (healthy)
satset-backend-prod    Up (healthy)
satset-frontend-prod   Up (healthy)              0.0.0.0:80->80/tcp
```

---

## 6. Verifikasi & Akses Aplikasi

### 6.1. Tes dari VPS

```bash
# Tes backend health
curl http://localhost:8080/healthz
# Expected: {"status":"ok"} atau response serupa

# Tes frontend
curl -I http://localhost
# Expected: HTTP/1.1 200 OK
```

### 6.2. Akses dari Browser

Buka browser di komputer kamu dan akses:

```
http://20.41.112.76:3000
```

Kamu seharusnya melihat halaman login SatSet! 🎉

> [!WARNING]
> Jika tidak bisa diakses, pastikan:
> - Port 3000 sudah dibuka di firewall VPS (lihat [Troubleshooting](#81-port-3000-tidak-bisa-diakses))
> - Tidak ada service lain yang sudah menggunakan port 3000

---

## 7. Maintenance & Operasional

### 7.1. Melihat Log

```bash
cd ~/apps/SatSet-App

# Log semua service
docker compose -f docker-compose.prod.yml --env-file .env.production logs -f

# Log service tertentu
docker compose -f docker-compose.prod.yml --env-file .env.production logs -f backend
docker compose -f docker-compose.prod.yml --env-file .env.production logs -f frontend
docker compose -f docker-compose.prod.yml --env-file .env.production logs -f db

# Log 100 baris terakhir
docker compose -f docker-compose.prod.yml --env-file .env.production logs --tail 100 backend
```

### 7.2. Update Aplikasi (Deploy Ulang)

Ketika ada perubahan kode baru:

```bash
cd ~/apps/SatSet-App

# 1. Pull kode terbaru
git pull origin main

# 2. Rebuild dan restart (tanpa downtime database)
docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build

# 3. Hapus image lama yang tidak terpakai
docker image prune -f
```

### 7.3. Stop & Start Aplikasi

```bash
cd ~/apps/SatSet-App

# Stop semua
docker compose -f docker-compose.prod.yml --env-file .env.production down

# Start ulang
docker compose -f docker-compose.prod.yml --env-file .env.production up -d

# Restart service tertentu
docker compose -f docker-compose.prod.yml --env-file .env.production restart backend
```

### 7.4. Backup Database

```bash
# Backup database ke file SQL
docker exec satset-db-prod mysqldump -u satset -p'PASSWORD_DB_KAMU' satset > backup_$(date +%Y%m%d_%H%M%S).sql

# Contoh dengan variable dari .env.production
source .env.production
docker exec satset-db-prod mysqldump -u $DB_USER -p"$DB_PASSWORD" $DB_NAME > backup_$(date +%Y%m%d_%H%M%S).sql
```

### 7.5. Restore Database

```bash
# Restore dari backup
docker exec -i satset-db-prod mysql -u satset -p'PASSWORD_DB_KAMU' satset < backup_20261017_120000.sql
```

### 7.6. Akses Database (MySQL CLI)

```bash
# Masuk ke MySQL shell
source .env.production
docker exec -it satset-db-prod mysql -u $DB_USER -p"$DB_PASSWORD" $DB_NAME
```

### 7.7. Melihat Penggunaan Disk

```bash
# Ukuran volume Docker
docker system df -v

# Bersihkan resource yang tidak terpakai
docker system prune -f
```

---

## 8. Troubleshooting

### 8.1. Port 3000 Tidak Bisa Diakses

```bash
# Cek firewall Ubuntu (UFW)
sudo ufw status

# Buka port 3000
sudo ufw allow 3000/tcp
sudo ufw reload

# Jika belum aktif, aktifkan (pastikan port 22 juga dibuka!)
sudo ufw allow 22/tcp
sudo ufw enable
```

Jika VPS kamu di **Azure** (berdasarkan IP 20.x.x.x):
- Buka **Azure Portal** → **Virtual Machine** → **Networking**
- Tambah **Inbound Security Rule**: Port 3000, Protocol TCP, Allow

### 8.2. Container Gagal Start

```bash
# Lihat log error
docker compose -f docker-compose.prod.yml --env-file .env.production logs backend
docker compose -f docker-compose.prod.yml --env-file .env.production logs db

# Cek container yang crash
docker ps -a
```

### 8.3. Database Tidak Bisa Konek

```bash
# Cek apakah container db sudah healthy
docker inspect satset-db-prod --format='{{.State.Health.Status}}'

# Cek log database
docker compose -f docker-compose.prod.yml --env-file .env.production logs db

# Tes koneksi manual
docker exec -it satset-db-prod mysql -u satset -p
```

### 8.4. Build Gagal (Timeout / Out of Memory)

```bash
# Tambah swap jika VPS RAM kecil (1-2GB)
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

# Lalu coba build ulang
docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build
```

### 8.5. Hapus Semua Data dan Mulai Ulang (Reset Total)

> [!CAUTION]
> Perintah ini menghapus SEMUA data termasuk database! Hanya gunakan jika memang ingin reset total.

```bash
# Stop dan hapus container, network, volume
docker compose -f docker-compose.prod.yml --env-file .env.production down -v

# Build ulang dari awal
docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build
```

---

## Ringkasan Perintah Penting

| Aksi | Perintah |
|---|---|
| **Build & Start** | `docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build` |
| **Lihat Status** | `docker compose -f docker-compose.prod.yml --env-file .env.production ps` |
| **Lihat Log** | `docker compose -f docker-compose.prod.yml --env-file .env.production logs -f` |
| **Stop** | `docker compose -f docker-compose.prod.yml --env-file .env.production down` |
| **Restart** | `docker compose -f docker-compose.prod.yml --env-file .env.production restart` |
| **Update / Deploy ulang** | `git pull && docker compose -f docker-compose.prod.yml --env-file .env.production up -d --build` |
| **Backup DB** | `docker exec satset-db-prod mysqldump -u USER -pPASS DB > backup.sql` |
| **Masuk MySQL** | `docker exec -it satset-db-prod mysql -u USER -pPASS DB` |

---

## Perbedaan Development vs Production

| Aspek | Development (`docker-compose.yml`) | Production (`docker-compose.prod.yml`) |
|---|---|---|
| **Port DB** | Exposed (3306) | ❌ Tidak exposed ke publik |
| **Restart** | `unless-stopped` | `always` |
| **Resource Limits** | Tidak ada | Ada (memory limits) |
| **Logging** | Default | JSON file dengan rotation |
| **Network** | Bisa diakses dari host | Hanya internal |
| **Port Frontend** | 3000 | 3000 (HTTP, karena port 80 sudah dipakai app lain) |
| **DB Optimization** | Default | InnoDB buffer pool, max connections |
| **APP_ENV** | `development` | `production` |
