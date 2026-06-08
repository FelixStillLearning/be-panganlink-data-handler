package main

import (
	"log"

	"github.com/example/be-panganlink-data-handler/internal/app"
	"github.com/example/be-panganlink-data-handler/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	db := config.InitDatabase(cfg)

	router := app.SetupRouter(db, cfg)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
