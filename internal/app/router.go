package app

import (
	"time"
	"github.com/example/be-panganlink-data-handler/internal/config"
	"github.com/example/be-panganlink-data-handler/internal/handler"
	"github.com/example/be-panganlink-data-handler/internal/middleware"
	"github.com/example/be-panganlink-data-handler/internal/repository"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/example/be-panganlink-data-handler/pkg/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust for production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Security Headers (Helmet-equivalent)
	r.Use(middleware.SecurityHeaders())

	// Rate Limiting (DDoS & Brute Force protection)
	r.Use(middleware.RateLimiter())

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	komoditasRepo := repository.NewKomoditasRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	authSvc := service.NewAuthService(userRepo, cfg)
	komoditasSvc := service.NewKomoditasService(komoditasRepo)
	productSvc := service.NewProductService(productRepo)
	aiSvc := service.NewAIService(cfg.AIServiceURL)
	paymentSvc := service.NewPaymentService(cfg.MidtransServerKey, false)
	orderSvc := service.NewOrderService(orderRepo, paymentSvc, notifRepo)

	// Cloud Storage
	var azureHelper *storage.AzureHelper
	if cfg.AzureAccountName != "" {
		// NewAzureHelper now expects connectionString and containerName
		// We'll pass AccountKey as the connection string or modify NewAzureHelper.
		// Usually ConnectionString looks like: DefaultEndpointsProtocol=https;AccountName=...;AccountKey=...;EndpointSuffix=core.windows.net
		connStr := "DefaultEndpointsProtocol=https;AccountName=" + cfg.AzureAccountName + ";AccountKey=" + cfg.AzureAccountKey + ";EndpointSuffix=core.windows.net"
		azureHelper, _ = storage.NewAzureHelper(connStr, cfg.AzureContainerName)
	}

	authHandler := handler.NewAuthHandler(authSvc)
	publicHandler := handler.NewPublicHandler(komoditasSvc)
	adminHandler := handler.NewAdminHandler(komoditasSvc, aiSvc, userRepo, productSvc)
	petaniHandler := handler.NewPetaniHandler(productSvc, aiSvc, orderSvc, userRepo)
	pembeliHandler := handler.NewPembeliHandler(orderSvc, userRepo)
	paymentHandler := handler.NewPaymentHandler(orderSvc)
	uploadHandler := handler.NewUploadHandler(azureHelper)
	notifHandler := handler.NewNotificationHandler(notifRepo)

	// Background Job: Cancel expired orders every hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			orderSvc.CancelExpiredOrders()
		}
	}()

	// Routes
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

		public := api.Group("/public")
		{
			public.GET("/stats", publicHandler.GetStats)
			public.GET("/commodities", publicHandler.GetCommodities)
			public.GET("/testimonials", publicHandler.GetTestimonials)
		}

		// Let's add upload as public for now, or authenticated. We will add it under /upload
		api.POST("/upload", uploadHandler.UploadImage)

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.RequireAuth(cfg.JWTSecret), authHandler.Me)
			auth.POST("/logout", middleware.RequireAuth(cfg.JWTSecret), authHandler.Logout)
		}

		payments := api.Group("/payments")
		{
			// Webhook Midtrans (No Auth Required)
			payments.POST("/webhook", paymentHandler.Webhook)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.RequireAuth(cfg.JWTSecret, "admin"))
		{
			admin.GET("/dashboard", adminHandler.Dashboard)
			admin.GET("/users", adminHandler.GetUsers)
			admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
			
			admin.GET("/products", adminHandler.GetProducts)
			admin.PUT("/products/:id/approve", adminHandler.ApproveProduct)
			admin.PUT("/products/:id/reject", adminHandler.RejectProduct)
			
			admin.GET("/commodities", adminHandler.GetCommodities)
			admin.POST("/commodities", adminHandler.CreateCommodity)
			admin.PUT("/commodities/:id", adminHandler.UpdateCommodity)
			admin.DELETE("/commodities/:id", adminHandler.DeleteCommodity)
			
			admin.GET("/market-prices", adminHandler.GetMarketPrices)
			admin.POST("/market-prices", adminHandler.CreateMarketPrice)
			admin.GET("/price-trends", adminHandler.GetPriceTrends)
			
			admin.GET("/settings", adminHandler.GetSettings)
			admin.PUT("/settings", adminHandler.UpdateSettings)
		}

		petani := api.Group("/petani")
		petani.Use(middleware.RequireAuth(cfg.JWTSecret, "petani"))
		{
			petani.GET("/dashboard", petaniHandler.Dashboard)
			
			petani.GET("/products", petaniHandler.GetProducts)
			petani.POST("/products", petaniHandler.CreateProduct)
			petani.PUT("/products/:id", petaniHandler.UpdateProduct)
			petani.DELETE("/products/:id", petaniHandler.DeleteProduct)
			
			petani.GET("/orders", petaniHandler.GetOrders)
			petani.PUT("/orders/:id/status", petaniHandler.UpdateOrderStatus)
			petani.GET("/history", petaniHandler.GetHistory)
			
			petani.GET("/recommendations", petaniHandler.GetRecommendations)
			
			petani.GET("/profile", petaniHandler.GetProfile)
			petani.PUT("/profile", petaniHandler.UpdateProfile)

			petani.GET("/notifications", notifHandler.GetNotifications)
			petani.PUT("/notifications/read-all", notifHandler.MarkAllAsRead)
			petani.PUT("/notifications/:id/read", notifHandler.MarkAsRead)
		}

		pembeli := api.Group("/pembeli")
		pembeli.Use(middleware.RequireAuth(cfg.JWTSecret, "pembeli"))
		{
			pembeli.GET("/dashboard", pembeliHandler.Dashboard)
			pembeli.GET("/products", pembeliHandler.GetProducts)
			pembeli.GET("/orders", pembeliHandler.GetOrders)
			pembeli.POST("/orders/checkout", pembeliHandler.Checkout)
			pembeli.PUT("/orders/:id/status", pembeliHandler.UpdateOrderStatus)
			pembeli.GET("/history", pembeliHandler.GetHistory)
			pembeli.GET("/profile", pembeliHandler.GetProfile)
			pembeli.PUT("/profile", pembeliHandler.UpdateProfile)

			pembeli.GET("/notifications", notifHandler.GetNotifications)
			pembeli.PUT("/notifications/read-all", notifHandler.MarkAllAsRead)
			pembeli.PUT("/notifications/:id/read", notifHandler.MarkAsRead)
		}
	}

	return r
}
