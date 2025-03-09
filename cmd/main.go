package main

import (
	"log"
	"mail-service/config"
	"mail-service/internal/database"
	"mail-service/internal/handlers"
	"mail-service/internal/repository"
	"mail-service/internal/router"
	"mail-service/internal/services"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database.InitDB(cfg)
	repo := repository.NewUserRepository(database.DB)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	r := router.SetupRouter(handler)
	r.Run(":" + cfg.Server.Port)
}
