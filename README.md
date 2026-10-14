# SatSet — All-in-One Laundry Management Platform

Platform manajemen operasional bisnis laundry terintegrasi dengan **AI Business Assistant**. Dirancang untuk laundry kiloan, satuan, sepatu, dan dry cleaning, dari single store hingga multi-cabang (multi-tenant ready).

Rujukan spesifikasi dan arsitektur lengkap: [SatSet.md](docs/SatSet.md).

---

## 🏗️ Tech Stack

- **Frontend (`apps/web`)**: Next.js 15 (App Router), React 19, TypeScript, Tailwind CSS, shadcn/ui, pnpm
- **Backend (`apps/api`)**: Go 1.24+, Gin Framework, sqlc, zerolog, validator
- **Database & Cache**: MySQL 8 (InnoDB, utf8mb4), Redis 7
- **Object Storage**: MinIO (S3 compatible)
- **AI Engine**: Google Gemini API via adapter backend (read-only pada MVP)
- **Container & CI/CD**: Docker multi-stage, Docker Compose, GitHub Actions

---

## 📁 Struktur Monorepo

```text
SatSet-App/
├── apps/
│   ├── api/                    # Backend modular monolith (Go + Gin)
│   │   ├── cmd/server/         # Entrypoint aplikasi server
│   │   ├── internal/
│   │   │   ├── http/           # HTTP Handlers, Middleware, Router (/healthz, /readyz)
│   │   │   └── platform/       # Config, Logger (zerolog), DB (MySQL), Redis
│   │   ├── Dockerfile          # Multi-stage minimal runner (non-root)
│   │   └── go.mod
│   └── web/                    # Frontend (Next.js 15 + Tailwind CSS + shadcn/ui)
│       ├── src/
│       │   ├── app/            # Next.js App Router (layout, page, styles)
│       │   ├── components/ui/  # Reusable UI components
│       │   └── lib/            # Utility helpers (clsx, tailwind-merge)
│       ├── Dockerfile          # Multi-stage standalone output
│       └── package.json
├── deploy/
│   ├── docker/mysql/           # MySQL init scripts (utf8mb4, UTC timezone)
│   └── nginx/                  # Nginx reverse proxy configuration
├── docs/                       # Dokumentasi arsitektur & PRD
├── .github/workflows/          # GitHub Actions CI pipeline
├── docker-compose.yml          # Local development stack (MySQL, Redis, MinIO, API, Web)
├── .env.example                # Template environment variables
├── Makefile                    # Perintah otomatisasi pengembangan
├── .golangci.yml               # Konfigurasi linter Go
├── .pre-commit-config.yaml     # Pre-commit hook configuration
└── AGENTS.md                   # Aturan & konvensi engineering proyek
```

---

## 🚀 Panduan Menjalankan Lokal

### 1. Prasyarat Sistem
- **Docker & Docker Compose**
- **Go 1.24+** (untuk pengembangan lokal API)
- **Node.js 22+** dan **pnpm** (untuk pengembangan lokal Web)

---

### 2. Konfigurasi Lingkungan
Salin berkas `.env.example` ke `.env`:
```bash
cp .env.example .env
```

---

### 3. Menjalankan Seluruh Stack (Docker Compose)
Cukup jalankan satu perintah:
```bash
make up
# Atau tanpa make:
docker compose up -d --build
```

Setelah semua container berjalan:
- **Web App**: [http://localhost:3000](http://localhost:3000)
- **API Server**: [http://localhost:8080](http://localhost:8080)
- **API Liveness Probe**: [http://localhost:8080/healthz](http://localhost:8080/healthz)
- **API Readiness Probe**: [http://localhost:8080/readyz](http://localhost:8080/readyz)
- **MinIO Console**: [http://localhost:9001](http://localhost:9001) (`minioadmin` / `minioadminpassword`)

Untuk menghentikan seluruh layanan:
```bash
make down
# Atau tanpa make:
docker compose down -v
```

---

### 4. Menjalankan Komponen Secara Terpisah (Local Development)

#### Backend (Go API)
```bash
cd apps/api
go run cmd/server/main.go
```
*Pastikan MySQL dan Redis sudah berjalan (dapat menggunakan `docker compose up -d mysql redis`).*

#### Frontend (Next.js)
```bash
cd apps/web
pnpm install
pnpm run dev
```
Akses di [http://localhost:3000](http://localhost:3000).

---

## 🧪 Testing & Linting

### Menjalankan Unit Test
```bash
make test
# Atau secara manual:
# Backend:
cd apps/api && go test -v -race ./...
# Frontend:
cd apps/web && pnpm run build
```

### Menjalankan Linter
```bash
make lint
# Atau secara manual:
# Backend:
cd apps/api && go vet ./...
# Frontend:
cd apps/web && pnpm run lint
```

---

## 🔒 Konvensi & Aturan Kunci
- Logika bisnis hanya berada di lapisan service backend (`handler -> service -> repository`).
- Setiap transaksi multi-tabel wajib atomik dalam transaksi database.
- Tipe data uang disimpan dalam satuan Rupiah integer `BIGINT` (tanpa float).
- Waktu disimpan dalam format UTC dan ditampilkan sesuai timezone bisnis (default Asia/Jakarta).
- API key Google Gemini hanya berada di backend, frontend tidak pernah mengakses LLM secara langsung.
