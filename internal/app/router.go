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

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Example Protected Route
		protected := api.Group("/protected")
		protected.Use(middleware.RequireAuth(cfg.JWTSecret))
		{
			protected.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{"user_id": c.GetString("user_id"), "role": c.GetString("role")})
			})
		}
	}

	return r
}
