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
	authSvc := service.NewAuthService(userRepo, cfg)
	authHandler := handler.NewAuthHandler(authSvc)

	// Routes
	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 1. Public Routes
		public := api.Group("/public")
		{
			public.GET("/stats", handler.GetPublicStats)
			public.GET("/commodities", handler.GetPublicCommodities)
			public.GET("/testimonials", handler.GetPublicTestimonials)
		}

		// 2. Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.RequireAuth(cfg.JWTSecret), authHandler.Me)
			auth.POST("/logout", middleware.RequireAuth(cfg.JWTSecret), authHandler.Logout)
		}

		// 3. Admin Routes
		admin := api.Group("/admin")
		admin.Use(middleware.RequireAuth(cfg.JWTSecret, "admin"))
		{
			admin.GET("/dashboard", handler.AdminDashboard)
			admin.GET("/users", handler.AdminGetUsers)
			admin.PUT("/users/:id/status", handler.AdminUpdateUserStatus)
			
			admin.GET("/products", handler.AdminGetProducts)
			admin.PUT("/products/:id/approve", handler.AdminApproveProduct)
			admin.PUT("/products/:id/reject", handler.AdminRejectProduct)
			
			admin.GET("/commodities", handler.AdminGetCommodities)
			admin.POST("/commodities", handler.AdminCreateCommodity)
			admin.PUT("/commodities/:id", handler.AdminUpdateCommodity)
			admin.DELETE("/commodities/:id", handler.AdminDeleteCommodity)
			
			admin.GET("/market-prices", handler.AdminGetMarketPrices)
			admin.POST("/market-prices", handler.AdminCreateMarketPrice)
			admin.GET("/price-trends", handler.AdminGetPriceTrends)
			
			admin.GET("/settings", handler.AdminGetSettings)
			admin.PUT("/settings", handler.AdminUpdateSettings)
		}

		// 4. Petani Routes
		petani := api.Group("/petani")
		petani.Use(middleware.RequireAuth(cfg.JWTSecret, "petani"))
		{
			petani.GET("/dashboard", handler.PetaniDashboard)
			
			petani.GET("/products", handler.PetaniGetProducts)
			petani.POST("/products", handler.PetaniCreateProduct)
			petani.PUT("/products/:id", handler.PetaniUpdateProduct)
			petani.DELETE("/products/:id", handler.PetaniDeleteProduct)
			
			petani.GET("/orders", handler.PetaniGetOrders)
			petani.PUT("/orders/:id/status", handler.PetaniUpdateOrderStatus)
			petani.GET("/history", handler.PetaniGetHistory)
			
			petani.GET("/recommendations", handler.PetaniGetRecommendations)
			
			petani.GET("/profile", handler.PetaniGetProfile)
			petani.PUT("/profile", handler.PetaniUpdateProfile)
		}
	}

	return r
}
