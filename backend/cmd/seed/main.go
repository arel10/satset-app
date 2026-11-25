package main

import (
	"context"
	"log"
	"time"

	"satset-app/internal/audit"
	"satset-app/internal/config"
	"satset-app/internal/database"
	"satset-app/internal/repository"
	"satset-app/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Pastikan data bawaan (settings, kategori pengeluaran, dan akun default) tersedia
	settingRepo := repository.NewSettingRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	userRepo := repository.NewUserRepository(db)
	auditLogger := audit.NewLogger(db)
	userService := service.NewUserService(userRepo, auditLogger)

	if err := settingRepo.InitDefaults(ctx); err != nil {
		log.Fatalf("Gagal inisialisasi pengaturan default: %v", err)
	}
	if err := expenseRepo.InitDefaultCategories(ctx); err != nil {
		log.Fatalf("Gagal inisialisasi kategori pengeluaran default: %v", err)
	}
	if err := userService.SeedDefaultUsers(ctx); err != nil {
		log.Fatalf("Gagal inisialisasi akun default: %v", err)
	}

	log.Println("Memulai seeding data sampel SatSet...")
	if err := database.SeedSampleData(ctx, db); err != nil {
		log.Fatalf("Gagal melakukan seeding: %v", err)
	}

	log.Println("Seeding data sampel SatSet selesai!")
}
