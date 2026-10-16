# DOKUMENTASI STRUKTUR DAN ARSITEKTUR BACKEND SATSET

Dokumen ini menjelaskan struktur kode, arsitektur, pola dependensi, dan aliran data aktual dari backend aplikasi **SatSet**.

---

## 1. Arsitektur Umum

Backend SatSet menerapkan **Layered Monolith Architecture (N-Tier)** dengan pemisahan tanggung jawab yang tegas antar-lapisan. Pola ini dipilih secara pragmatis untuk bisnis laundry tunggal guna memastikan kode mudah dibaca, mudah diuji, performa tinggi, dan bebas dari abstraksi berlebihan (mengikuti prinsip PRD Section 8.1 dan AGENTS.md Section 5).

### Diagram Lapisan (Layer Hierarchy)
```text
┌────────────────────────────────────────────────────────┐
│                   HTTP Requests (Client)               │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│                  Cross-Cutting Middlewares             │
│   (CORS, Recovery, Logger, JWT Auth, Role RBAC, Rate)  │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│                  Presentation Layer                    │
│            internal/handler (HTTP Controllers)         │
└───────────────────────────┬────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────┐
│                   Business Logic Layer                 │
│         internal/service (Use Cases & Orchestration)   │
└──────────────┬─────────────────────────┬───────────────┘
               │                         │
               ▼                         ▼
┌────────────────────────────┐ ┌─────────────────────────┐
│     Data Access Layer      │ │   External Subsystem    │
│  internal/repository (ORM) │ │ internal/ai (Gemini AI) │
└──────────────┬─────────────┘ └─────────────────────────┘
               │
               ▼
┌────────────────────────────────────────────────────────┐
│                     Database Layer                     │
│                  MySQL 8 / InnoDB Engine               │
└────────────────────────────────────────────────────────┘
```

---

## 2. Struktur Direktori dan Komponen

```text
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point, bootstrap, wiring, graceful shutdown
├── internal/
│   ├── ai/
│   │   ├── prompts/
│   │   │   └── system_prompt.go    # System prompt SatSet, aturan read-only, anti-injection
│   │   ├── provider/
│   │   │   └── gemini.go           # Gemini REST API provider client (HTTP standard library)
│   │   └── tools/
│   │       ├── tools.go            # 13 definisi & eksekutor read-only tools AI
│   │       └── tools_test.go       # Unit test verifikasi 13 tools read-only
│   ├── audit/
│   │   └── audit.go                # Logger audit trail (Log dan LogWithTx ke database)
│   ├── config/
│   │   └── config.go               # Loader & validator environment variables
│   ├── database/
│   │   └── database.go             # Koneksi GORM MySQL & connection pooling
│   ├── handler/                    # HTTP Controllers (Gin)
│   │   ├── auth_handler.go         # Login, refresh, logout, me
│   │   ├── customer_handler.go     # Pelanggan CRUD & order history
│   │   ├── dashboard_handler.go    # Dashboard Owner & Kasir
│   │   ├── expense_handler.go      # Pengeluaran & kategori pengeluaran
│   │   ├── order_handler.go        # POS order creation, list, lookup, status, cancel
│   │   ├── payment_handler.go      # Pencatatan bayar lunas & pembatalan
│   │   ├── report_handler.go       # Laporan bisnis periodik
│   │   ├── response.go             # Standard response builder & pagination parser
│   │   ├── service_handler.go      # Layanan & tarif laundry
│   │   ├── setting_handler.go      # Pengaturan toko & teks nota
│   │   ├── user_handler.go         # Manajemen akun pengguna
│   │   └── ai_handler.go           # AI Assistant chat & conversation history
│   ├── middleware/                 # HTTP Middlewares
│   │   ├── auth_middleware.go      # Validasi Bearer JWT token
│   │   ├── role_middleware.go      # Otorisasi RBAC (Owner, Kasir, Operator)
│   │   ├── rate_limit_middleware.go# In-memory limiter proteksi brute-force login
│   │   └── cors_middleware.go      # Header CORS (termasuk Idempotency-Key)
│   ├── model/                      # Database Models (GORM)
│   │   ├── user.go / setting.go / customer.go / service.go
│   │   ├── order.go / order_item.go / order_status_history.go / order_sequence.go
│   │   ├── payment.go / expense.go / idempotency_key.go / audit_log.go
│   │   ├── ai_conversation.go
│   │   └── order_test.go           # Unit test transisi status & status terminal
│   ├── repository/                 # Data Access Layer
│   │   ├── user_repository.go / setting_repository.go / customer_repository.go
│   │   ├── service_repository.go / order_repository.go / payment_repository.go
│   │   ├── expense_repository.go / ai_repository.go / audit_repository.go
│   ├── routes/
│   │   └── routes.go               # Pendaftaran seluruh rute & wiring middleware
│   └── service/                    # Business Logic Layer
│       ├── auth_service.go / auth_service_test.go
│       ├── user_service.go / setting_service.go
│       ├── customer_service.go / customer_service_test.go
│       ├── service_service.go / order_service.go / payment_service.go
│       ├── expense_service.go / dashboard_service.go / report_service.go
│       └── ai_service.go
├── migrations/
│   ├── 000001_initial_schema.up.sql   # DDL MySQL 8 InnoDB (16 tabel)
│   └── 000001_initial_schema.down.sql # DDL Rollback
├── .env.example                    # Template konfigurasi environment
├── go.mod / go.sum                 # Dependency Go modules
```

---

## 3. Penjelasan Rinci Komponen

### 3.1 Entry Point & Bootstrap (`cmd/api/main.go`)
- **Tanggung Jawab:** Titik awal eksekusi backend.
- **Tugas:**
  1. Membaca konfigurasi dari `.env` dan sistem menggunakan `config.Load()`.
  2. Menginisialisasi koneksi MySQL pool menggunakan `database.Connect()`.
  3. Menjalankan auto-migrasi skema dan inisialisasi data bawaan (seeding default Owner, default setting, default kategori pengeluaran).
  4. Melakukan instansiasi repositori, audit logger, provider AI, layanan bisnis (services), dan pengontrol HTTP (handlers).
  5. Mendaftarkan seluruh rute dan middleware pada Gin engine melalui `routes.SetupRouter()`.
  6. Membuka server HTTP pada port konfigurasi (default: 8080) dan menangani *graceful shutdown* menggunakan sinyal OS (`SIGINT`, `SIGTERM`).

### 3.2 Configuration & Environment (`internal/config/config.go`)
- **Tanggung Jawab:** Memuat variabel lingkungan (`.env`), menyediakan fallback nilai default, dan memvalidasi keberadaan secret penting (`JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `DB_NAME`).
- **Parsing:** Menggunakan parser khusus untuk integer, boolean, dan lokasi zona waktu `Asia/Jakarta`.

### 3.3 Database Access & Connection Pool (`internal/database/database.go`)
- **Tanggung Jawab:** Konfigurasi koneksi GORM ke MySQL 8 InnoDB dengan parameter `loc=UTC`, `parseTime=True`, dan `charset=utf8mb4`.
- **Connection Pool:** Mengatur `SetMaxIdleConns(10)`, `SetMaxOpenConns(100)`, dan `SetConnMaxLifetime(time.Hour)`.

### 3.4 Middlewares (`internal/middleware/`)
- `auth_middleware.go`: Membaca header `Authorization: Bearer <token>`, memvalidasi signature HMAC-SHA256, dan menyematkan `user_id`, `username`, serta `role` ke dalam Gin Context.
- `role_middleware.go`: Memeriksa apakah `role` pengguna sesuai dengan daftar peran yang diizinkan pada rute tersebut (Role-Based Access Control).
- `rate_limit_middleware.go`: Limiter in-memory thread-safe (`sync.Mutex`) berbasis sliding window untuk membatasi percobaan login gagal per IP client (default 5 percobaan per 15 menit).
- `cors_middleware.go`: Mengizinkan komunikasi lintas domain dengan Next.js, mendukung header `Authorization` dan `Idempotency-Key`.

### 3.5 Handlers (`internal/handler/`)
- Bertanggung jawab menerima request HTTP, memvalidasi format JSON (`binding:"required"`), mengekstrak parameter/query, memanggil service terkait, dan mengembalikan response terstandarisasi.
- Seluruh response dibentuk melalui utilitas `response.go`:
  ```json
  {
    "success": true,
    "message": "Pesan status",
    "data": { ... },
    "error": null,
    "meta": { "page": 1, "per_page": 20, "total": 100, "total_pages": 5 }
  }
  ```

### 3.6 Services (`internal/service/`)
- Merupakan inti dari aturan bisnis (business logic):
  - `OrderService`: Kalkulasi harga snapshot, subtotal, diskon, urutan status, dan idempotensi.
  - `PaymentService`: Validasi pembayaran lunas penuh, validasi sisa tagihan, pembatalan pembayaran Owner.
  - `CustomerService`: Normalisasi nomor HP ke format E.164 Indonesia (`+628xx`), pencegahan duplikasi.
  - `ServiceService`: Validasi satuan layanan (`kg`, `pcs`, `pasang`), pengaturan harga non-negatif.
  - `ExpenseService`: Pencatatan beban operasional khusus Owner, filter tanggal, soft-cancel.
  - `DashboardService`: Agregasi metrik finansial dan operasional laundry.
  - `ReportService`: Perhitungan 10 jenis laporan bisnis terperinci.
  - `AIService`: Orkestrasi sesi percakapan Gemini AI, loop tool execution, token usage tracking.
  - `UserService`: Hashing password bcrypt, manajemen status akun, revokasi sesi refresh token.

### 3.7 Repositories (`internal/repository/`)
- Melakukan query langsung ke database menggunakan GORM.
- Mengimplementasikan penguncian baris (`clause.Locking{Strength: "UPDATE"}`) untuk nomor pesanan berurutan.
- Menggunakan query terparameter untuk mencegah SQL Injection.

### 3.8 AI Subsystem (`internal/ai/`)
- `provider/gemini.go`: Klien HTTP murni ke Google Gemini API v1beta (`gemini-2.0-flash` / `gemini-1.5-flash`), menangani format message, system instruction, function declarations, dan function call handling.
- `tools/tools.go`: 13 tool deklarasi read-only dan eksekutor kueri database lokal.
- `prompts/system_prompt.go`: Instruksi sistem yang menegakkan batasan read-only, bahasa Indonesia yang ramah, dan pencegahan prompt injection.

### 3.9 Audit Trail (`internal/audit/audit.go`)
- Mencatat aksi audit secara transaksional (`LogWithTx`) atau langsung (`Log`) untuk setiap perubahan kritis: pembuatan pesanan, pembatalan, pembayaran, perubahan harga, login/logout, dan chat AI.

---

## 4. Aliran Data End-to-End: Pembuatan Pesanan (POS)

```text
[Kasir di Halaman POS Frontend]
              │
              │ HTTP POST /api/orders
              │ Headers: Authorization: Bearer <jwt>, Idempotency-Key: <uuid>
              │ Body: { customer_id: 1, items: [...], discount: 5000 }
              ▼
    ┌───────────────────┐
    │  CORSMiddleware   │ -> Validasi Origin, Methods, Headers
    └─────────┬─────────┘
              ▼
    ┌───────────────────┐
    │  AuthMiddleware   │ -> Validasi JWT Token, set user_id=2, role="kasir"
    └─────────┬─────────┘
              ▼
    ┌───────────────────┐
    │    RequireRole    │ -> Cek apakah "kasir" diizinkan (Allowed: owner, kasir)
    └─────────┬─────────┘
              ▼
    ┌───────────────────┐
    │OrderHandler.Create│ -> c.ShouldBindJSON(&req)
    └─────────┬─────────┘
              ▼
   ┌─────────────────────┐
   │ OrderService.Create ├─► Cek IdempotencyKey di DB (jika ada, return response cache)
   │                     ├─► Cek Customer di CustomerRepository
   │                     ├─► Loop items: ambil tarif dari ServiceRepository, hitung Subtotal
   │                     ├─► Validasi: discount <= subtotal
   │                     ├─► Hitung Estimasi Selesai (durasi jam terpanjang)
   │                     │
   │                     ▼
   │            [DB Transaction (tx)]
   │             ├── 1. GenerateOrderNumber (Lock order_sequences: "LDR-261014-001")
   │             ├── 2. Insert orders (total, work_status='DITERIMA', payment_status='BELUM_LUNAS')
   │             ├── 3. Insert order_items (snapshot harga & nama)
   │             ├── 4. Insert order_status_history (initial status)
   │             ├── 5. Insert idempotency_keys (TTL 24 jam)
   │             └── 6. Insert audit_logs (Action: CREATE, Entity: order)
   └──────────┬──────────┘
              ▼
    ┌───────────────────┐
    │RespondCreated 201 │ -> Kirim Order JSON lengkap ke Client (Frontend)
    └───────────────────┘
```

---

## 5. Ringkasan Endpoint API

| Kelompok | Method | Endpoint | Hak Akses | Deskripsi |
|---|---|---|---|---|
| **Health** | `GET` | `/healthz` | Publik | Health check status server |
| **Auth** | `POST` | `/api/auth/login` | Publik (Rate limited) | Login username/email + password |
| **Auth** | `POST` | `/api/auth/refresh` | Publik | Memperbarui access token |
| **Auth** | `GET` | `/api/auth/me` | Authenticated | Info pengguna saat ini |
| **Auth** | `POST` | `/api/auth/logout` | Authenticated | Logout & cabut sesi |
| **Users** | `GET, POST` | `/api/users` | Owner | Daftar & buat pengguna baru |
| **Users** | `GET, PUT, DEL` | `/api/users/:id` | Owner | Detail, edit, soft-delete pengguna |
| **Settings**| `GET` | `/api/settings` | Authenticated | Ambil konfigurasi bisnis |
| **Settings**| `PUT` | `/api/settings` | Owner | Perbarui konfigurasi bisnis |
| **Customers**| `POST, GET` | `/api/customers` | Owner, Kasir | Tambah & cari pelanggan |
| **Customers**| `GET, PUT` | `/api/customers/:id` | Owner, Kasir | Detail & perbarui data pelanggan |
| **Customers**| `GET` | `/api/customers/:id/orders`| Owner, Kasir | Riwayat pesanan pelanggan |
| **Services**| `GET` | `/api/services/active` | Owner, Kasir, Operator | Daftar layanan aktif untuk POS |
| **Services**| `GET, POST` | `/api/services` | Owner | Daftar lengkap & tambah layanan |
| **Services**| `GET, PUT` | `/api/services/:id` | Owner | Detail & perbarui tarif/layanan |
| **Orders** | `POST` | `/api/orders` | Owner, Kasir | Buat pesanan baru (POS) |
| **Orders** | `GET` | `/api/orders` | Owner, Kasir, Operator | Daftar & filter pesanan |
| **Orders** | `GET` | `/api/orders/:id` | Owner, Kasir, Operator | Detail lengkap pesanan |
| **Orders** | `GET` | `/api/orders/number/:num` | Owner, Kasir, Operator | Cari nomor nota pesanan |
| **Orders** | `GET` | `/api/orders/:id/payments` | Owner, Kasir, Operator | Riwayat pembayaran pesanan |
| **Orders** | `PATCH` | `/api/orders/:id/status` | Owner, Kasir, Operator | Perbarui status pengerjaan |
| **Orders** | `POST` | `/api/orders/:id/cancel` | Owner, Kasir | Batalkan pesanan dengan alasan |
| **Payments**| `POST` | `/api/payments` | Owner, Kasir | Catat pembayaran lunas penuh |
| **Payments**| `PUT` | `/api/payments/:id/cancel` | Owner | Batalkan pembayaran |
| **Expenses**| `POST, GET` | `/api/expenses` | Owner | Catat & daftar pengeluaran |
| **Expenses**| `GET, PUT` | `/api/expenses/:id` | Owner | Detail & perbarui pengeluaran |
| **Expenses**| `POST` | `/api/expenses/:id/cancel` | Owner | Batalkan pengeluaran |
| **Expenses**| `GET, POST` | `/api/expenses/categories` | Owner | Kategori pengeluaran |
| **Dashboard**| `GET` | `/api/dashboard/owner` | Owner | Dashboard finansial & operasional |
| **Dashboard**| `GET` | `/api/dashboard/kasir` | Owner, Kasir | Dashboard operasional kasir |
| **Reports** | `GET` | `/api/reports/business` | Owner | 10 jenis laporan bisnis periodik |
| **AI** | `POST` | `/api/ai/chat` | Owner | Tanya jawab Asisten AI (Gemini) |
| **AI** | `GET` | `/api/ai/conversations` | Owner | Riwayat percakapan AI |
| **AI** | `GET, DEL` | `/api/ai/conversations/:id`| Owner | Detail & hapus percakapan AI |
| **AI** | `GET` | `/api/ai/tools` | Owner | Daftar tools AI yang terdaftar |
