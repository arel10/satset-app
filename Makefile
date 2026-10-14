.PHONY: help up down test lint migrate-up migrate-down dev-api dev-web

help: ## Menampilkan bantuan perintah
	@echo "Perintah yang tersedia:"
	@echo "  make up           - Menjalankan seluruh stack via Docker Compose"
	@echo "  make down         - Menghentikan seluruh container dan menghapus volume"
	@echo "  make test         - Menjalankan unit test backend dan frontend"
	@echo "  make lint         - Menjalankan linter untuk backend dan frontend"
	@echo "  make migrate-up   - Menjalankan migrasi database ke versi terbaru"
	@echo "  make migrate-down - Membatalkan satu langkah migrasi database"
	@echo "  make dev-api      - Menjalankan backend Go secara lokal"
	@echo "  make dev-web      - Menjalankan frontend Next.js secara lokal"

up: ## Menjalankan docker compose (MySQL, Redis, MinIO, API, Web)
	docker compose up -d --build

down: ## Menghentikan stack
	docker compose down -v

test: test-backend test-frontend ## Menjalankan seluruh test

test-backend: ## Menjalankan test Go API
	cd apps/api && go test -v -race ./...

test-frontend: ## Menjalankan validasi build frontend
	cd apps/web && pnpm run build

lint: lint-backend lint-frontend ## Menjalankan seluruh linter

lint-backend: ## Menjalankan go vet & linter backend
	cd apps/api && go vet ./...

lint-frontend: ## Menjalankan Next.js ESLint
	cd apps/web && pnpm run lint

migrate-up: ## Menjalankan migrasi DB
	@echo "Migrasi database siap diintegrasikan pada Tahap 2"

migrate-down: ## Rollback migrasi DB
	@echo "Rollback database siap diintegrasikan pada Tahap 2"

dev-api: ## Menjalankan API lokal
	cd apps/api && go run cmd/server/main.go

dev-web: ## Menjalankan Web lokal
	cd apps/web && pnpm run dev
