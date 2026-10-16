package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"satset-app/internal/ai/provider"
	"satset-app/internal/ai/tools"
	"satset-app/internal/audit"
	"satset-app/internal/config"
	"satset-app/internal/database"
	"satset-app/internal/handler"
	"satset-app/internal/model"
	"satset-app/internal/repository"
	"satset-app/internal/routes"
	"satset-app/internal/service"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	log.Printf("Starting %s in %s mode on port %s", cfg.App.Name, cfg.App.Env, cfg.App.Port)

	// 2. Connect to Database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// In development, ensure tables are present
	if cfg.IsDevelopment() {
		log.Println("Development mode: running AutoMigrate...")
		if err := db.AutoMigrate(
			&model.User{},
			&model.RefreshToken{},
			&model.Setting{},
			&model.Customer{},
			&model.Service{},
			&model.Order{},
			&model.OrderItem{},
			&model.OrderStatusHistory{},
			&model.OrderSequence{},
			&model.Payment{},
			&model.ExpenseCategory{},
			&model.Expense{},
			&model.AIConversation{},
			&model.AIMessage{},
			&model.AuditLog{},
			&model.IdempotencyKey{},
		); err != nil {
			log.Printf("Peringatan AutoMigrate: %v", err)
		}
	}

	// 3. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	settingRepo := repository.NewSettingRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	aiRepo := repository.NewAIRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// 4. Initialize Audit Logger
	auditLogger := audit.NewLogger(db)

	// 5. Seed default data (Owner, default categories, default settings)
	initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer initCancel()

	if err := settingRepo.InitDefaults(initCtx); err != nil {
		log.Printf("Peringatan init settings: %v", err)
	}
	if err := expenseRepo.InitDefaultCategories(initCtx); err != nil {
		log.Printf("Peringatan init expense categories: %v", err)
	}

	// 6. Initialize AI Provider & Tools
	geminiProvider := provider.NewGeminiProvider(cfg.Gemini.APIKey, cfg.Gemini.Model, cfg.AI.RequestTimeoutSecs)
	toolExecutor := tools.NewToolExecutor(db, orderRepo, paymentRepo, expenseRepo, customerRepo, serviceRepo)

	// 7. Initialize Services
	authService := service.NewAuthService(userRepo, cfg, auditLogger)
	userService := service.NewUserService(userRepo, auditLogger)
	settingService := service.NewSettingService(settingRepo, auditLogger)
	customerService := service.NewCustomerService(customerRepo)
	serviceService := service.NewServiceService(serviceRepo, auditLogger)
	orderService := service.NewOrderService(orderRepo, serviceRepo, customerRepo, paymentRepo, auditLogger)
	paymentService := service.NewPaymentService(paymentRepo, orderRepo, auditLogger)
	expenseService := service.NewExpenseService(expenseRepo, auditLogger)
	dashboardService := service.NewDashboardService(orderRepo, paymentRepo, expenseRepo, customerRepo, serviceRepo, cfg)
	reportService := service.NewReportService(orderRepo, paymentRepo, expenseRepo, serviceRepo, customerRepo, cfg)
	aiService := service.NewAIService(aiRepo, geminiProvider, toolExecutor, cfg, auditLogger)
	auditService := service.NewAuditService(auditRepo)

	// Seed default Owner account if no owner exists
	if err := userService.SeedDefaultOwner(initCtx); err != nil {
		log.Printf("Peringatan seeding default owner: %v", err)
	}

	// 8. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	settingHandler := handler.NewSettingHandler(settingService)
	customerHandler := handler.NewCustomerHandler(customerService, orderService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	reportHandler := handler.NewReportHandler(reportService)
	aiHandler := handler.NewAIHandler(aiService)
	auditHandler := handler.NewAuditHandler(auditService)

	// 9. Setup Router
	router := routes.SetupRouter(routes.RouterConfig{
		Config:           cfg,
		AuthService:      authService,
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		SettingHandler:   settingHandler,
		CustomerHandler:  customerHandler,
		ServiceHandler:   serviceHandler,
		OrderHandler:     orderHandler,
		PaymentHandler:   paymentHandler,
		ExpenseHandler:   expenseHandler,
		DashboardHandler: dashboardHandler,
		ReportHandler:    reportHandler,
		AIHandler:        aiHandler,
		AuditHandler:     auditHandler,
	})

	// 10. Background periodic cleaner for expired tokens and idempotency keys (every 1 hour)
	cleanerStop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cleanCtx, cleanCancel := context.WithTimeout(context.Background(), 30*time.Second)
				_ = userRepo.DeleteExpiredRefreshTokens(cleanCtx)
				_ = orderRepo.DeleteExpiredIdempotencyKeys(cleanCtx)
				cleanCancel()
			case <-cleanerStop:
				return
			}
		}
	}()

	// 11. Start HTTP Server with Graceful Shutdown
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("SatSet backend running at http://localhost:%s", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down SatSet backend server...")
	close(cleanerStop)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("SatSet backend server exited cleanly")
}
