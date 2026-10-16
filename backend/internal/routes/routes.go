package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"satset-app/internal/config"
	"satset-app/internal/handler"
	"satset-app/internal/middleware"
	"satset-app/internal/model"
	"satset-app/internal/service"
)

// RouterConfig contains all dependencies needed to configure application routes.
type RouterConfig struct {
	Config           *config.Config
	AuthService      *service.AuthService
	AuthHandler      *handler.AuthHandler
	UserHandler      *handler.UserHandler
	SettingHandler   *handler.SettingHandler
	CustomerHandler  *handler.CustomerHandler
	ServiceHandler   *handler.ServiceHandler
	OrderHandler     *handler.OrderHandler
	PaymentHandler   *handler.PaymentHandler
	ExpenseHandler   *handler.ExpenseHandler
	DashboardHandler *handler.DashboardHandler
	ReportHandler    *handler.ReportHandler
	AIHandler        *handler.AIHandler
	AuditHandler     *handler.AuditHandler
}

// SetupRouter initializes the Gin router with all middlewares and route definitions.
func SetupRouter(rc RouterConfig) *gin.Engine {
	if rc.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(rc.Config.CORS.FrontendURL))

	// Static file serving for uploads (logos, etc.)
	r.Static("/uploads", "./uploads")

	// Health check endpoint (GET and HEAD supported for Docker healthcheck)
	healthHandler := func(c *gin.Context) {
		handler.RespondOK(c, "SatSet Backend is healthy", gin.H{
			"app":       rc.Config.App.Name,
			"status":    "UP",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
	r.GET("/healthz", healthHandler)
	r.HEAD("/healthz", healthHandler)

	// Rate limiter for login: 5 attempts per 15 minutes by default
	maxAttempts := rc.Config.RateLimit.LoginMaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	lockoutMins := rc.Config.RateLimit.LoginLockoutMinutes
	if lockoutMins <= 0 {
		lockoutMins = 15
	}
	loginLimiter := middleware.NewMemoryRateLimiter(maxAttempts, time.Duration(lockoutMins)*time.Minute, 15*time.Minute)

	// API Group
	api := r.Group("/api")
	{
		// Public Authentication Routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", middleware.LoginRateLimitMiddleware(loginLimiter), rc.AuthHandler.Login)
			authGroup.POST("/refresh", rc.AuthHandler.RefreshToken)
		}

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(rc.AuthService))
		{
			// Auth context
			protected.GET("/auth/me", rc.AuthHandler.Me)
			protected.POST("/auth/logout", rc.AuthHandler.Logout)

			// 1. User Management (Owner only)
			userGroup := protected.Group("/users")
			userGroup.Use(middleware.RequireRole(model.RoleOwner))
			{
				userGroup.POST("", rc.UserHandler.Create)
				userGroup.GET("", rc.UserHandler.List)
				userGroup.GET("/:id", rc.UserHandler.GetByID)
				userGroup.PUT("/:id", rc.UserHandler.Update)
				userGroup.DELETE("/:id", rc.UserHandler.Delete)
			}

			// 2. Business Settings
			settingsGroup := protected.Group("/settings")
			{
				settingsGroup.GET("", rc.SettingHandler.GetAll)
				settingsGroup.PUT("", middleware.RequireRole(model.RoleOwner), rc.SettingHandler.Update)
				settingsGroup.POST("/logo", middleware.RequireRole(model.RoleOwner), rc.SettingHandler.UploadLogo)
			}

			// 3. Customers (Owner & Kasir)
			customerGroup := protected.Group("/customers")
			customerGroup.Use(middleware.RequireRole(model.RoleOwner, model.RoleKasir))
			{
				customerGroup.POST("", rc.CustomerHandler.Create)
				customerGroup.GET("", rc.CustomerHandler.Search)
				customerGroup.GET("/:id", rc.CustomerHandler.GetByID)
				customerGroup.PUT("/:id", rc.CustomerHandler.Update)
				customerGroup.GET("/:id/orders", rc.CustomerHandler.GetOrderHistory)
			}

			// 4. Laundry Services
			serviceGroup := protected.Group("/services")
			{
				// Active services for POS: Owner, Kasir, Operator
				serviceGroup.GET("/active", middleware.RequireRole(model.RoleOwner, model.RoleKasir, model.RoleOperator), rc.ServiceHandler.GetActive)
				// Full listing & CRUD: Owner only
				serviceGroup.GET("", middleware.RequireRole(model.RoleOwner), rc.ServiceHandler.List)
				serviceGroup.POST("", middleware.RequireRole(model.RoleOwner), rc.ServiceHandler.Create)
				serviceGroup.GET("/:id", middleware.RequireRole(model.RoleOwner), rc.ServiceHandler.GetByID)
				serviceGroup.PUT("/:id", middleware.RequireRole(model.RoleOwner), rc.ServiceHandler.Update)
			}

			// 5. Orders & POS
			orderGroup := protected.Group("/orders")
			{
				// Create order (POS): Owner & Kasir
				orderGroup.POST("", middleware.RequireRole(model.RoleOwner, model.RoleKasir), rc.OrderHandler.Create)
				// List orders: Owner, Kasir, Operator
				orderGroup.GET("", middleware.RequireRole(model.RoleOwner, model.RoleKasir, model.RoleOperator), rc.OrderHandler.List)
				orderGroup.GET("/number/:orderNumber", middleware.RequireRole(model.RoleOwner, model.RoleKasir, model.RoleOperator), rc.OrderHandler.GetByNumber)
				orderGroup.GET("/:id", middleware.RequireRole(model.RoleOwner, model.RoleKasir, model.RoleOperator), rc.OrderHandler.GetByID)
				orderGroup.GET("/:id/payments", middleware.RequireRole(model.RoleOwner, model.RoleKasir, model.RoleOperator), rc.PaymentHandler.GetByOrder)
				// Update work status: Owner, Kasir, Operator
				orderGroup.PATCH("/:id/status", middleware.RequireRole(model.RoleOwner, model.RoleKasir, model.RoleOperator), rc.OrderHandler.UpdateStatus)
				// Cancel order: Owner & Kasir
				orderGroup.POST("/:id/cancel", middleware.RequireRole(model.RoleOwner, model.RoleKasir), rc.OrderHandler.Cancel)
			}

			// 6. Payments
			paymentGroup := protected.Group("/payments")
			{
				// Record payment (full payment): Owner & Kasir
				paymentGroup.POST("", middleware.RequireRole(model.RoleOwner, model.RoleKasir), rc.PaymentHandler.RecordPayment)
				// Cancel payment: Owner only
				paymentGroup.PUT("/:id/cancel", middleware.RequireRole(model.RoleOwner), rc.PaymentHandler.CancelPayment)
			}

			// 7. Expenses (Owner only)
			expenseGroup := protected.Group("/expenses")
			expenseGroup.Use(middleware.RequireRole(model.RoleOwner))
			{
				expenseGroup.POST("", rc.ExpenseHandler.Create)
				expenseGroup.GET("", rc.ExpenseHandler.List)
				expenseGroup.GET("/:id", rc.ExpenseHandler.GetByID)
				expenseGroup.PUT("/:id", rc.ExpenseHandler.Update)
				expenseGroup.POST("/:id/cancel", rc.ExpenseHandler.Cancel)

				// Expense Categories
				expenseGroup.GET("/categories", rc.ExpenseHandler.ListCategories)
				expenseGroup.POST("/categories", rc.ExpenseHandler.CreateCategory)
				expenseGroup.PUT("/categories/:id", rc.ExpenseHandler.UpdateCategory)
			}

			// 8. Dashboards
			dashboardGroup := protected.Group("/dashboard")
			{
				dashboardGroup.GET("/owner", middleware.RequireRole(model.RoleOwner), rc.DashboardHandler.GetOwnerDashboard)
				dashboardGroup.GET("/kasir", middleware.RequireRole(model.RoleOwner, model.RoleKasir), rc.DashboardHandler.GetKasirDashboard)
			}

			// 9. Reports (Owner only)
			reportGroup := protected.Group("/reports")
			reportGroup.Use(middleware.RequireRole(model.RoleOwner))
			{
				reportGroup.GET("/business", rc.ReportHandler.GetBusinessReport)
			}

			// 10. AI Assistant (Owner only)
			aiGroup := protected.Group("/ai")
			aiGroup.Use(middleware.RequireRole(model.RoleOwner))
			{
				aiGroup.POST("/chat", rc.AIHandler.Chat)
				aiGroup.GET("/conversations", rc.AIHandler.GetConversations)
				aiGroup.GET("/conversations/:id", rc.AIHandler.GetConversationDetail)
				aiGroup.DELETE("/conversations/:id", rc.AIHandler.DeleteConversation)
				aiGroup.GET("/tools", rc.AIHandler.ListTools)
			}

			// 11. Audit Logs (Owner only)
			auditGroup := protected.Group("/audit-logs")
			auditGroup.Use(middleware.RequireRole(model.RoleOwner))
			{
				auditGroup.GET("", rc.AuditHandler.List)
			}
		}
	}

	return r
}
