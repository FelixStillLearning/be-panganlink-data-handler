package app

import (
	"github.com/example/be-panganlink-data-handler/internal/config"
	"github.com/example/be-panganlink-data-handler/internal/handler"
	"github.com/example/be-panganlink-data-handler/internal/middleware"
	"github.com/example/be-panganlink-data-handler/internal/repository"
	"github.com/example/be-panganlink-data-handler/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	komoditasRepo := repository.NewKomoditasRepository(db)
	productRepo := repository.NewProductRepository(db)

	authSvc := service.NewAuthService(userRepo, cfg)
	komoditasSvc := service.NewKomoditasService(komoditasRepo)
	productSvc := service.NewProductService(productRepo)
	aiSvc := service.NewAIService(cfg.AIServiceURL)

	authHandler := handler.NewAuthHandler(authSvc)
	publicHandler := handler.NewPublicHandler(komoditasSvc)
	adminHandler := handler.NewAdminHandler(komoditasSvc, aiSvc)
	petaniHandler := handler.NewPetaniHandler(productSvc, aiSvc)

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

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.RequireAuth(cfg.JWTSecret), authHandler.Me)
			auth.POST("/logout", middleware.RequireAuth(cfg.JWTSecret), authHandler.Logout)
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
		}
	}

	return r
}
