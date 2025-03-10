package internal

import (
	"context"
	"fmt"
	"mail-service/config"
	"mail-service/internal/auth"
	"mail-service/internal/database"
	"mail-service/internal/handlers"
	"mail-service/internal/middleware"
	"mail-service/internal/repository"
	"mail-service/internal/router"
	"mail-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	server *http.Server
	db     *database.Database
	logger *zap.Logger
}

func NewApp() *App {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	auth.SetJWTKey(cfg.Auth.JwtSecret)
	logger, _ := zap.NewProduction()

	db, err := database.InitDB(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}
	if db.DB == nil {
		logger.Fatal("Database.DB is nil after initialization")
	}

	// Repositories
	userRepo := repository.NewUserRepository(db.DB)
	msgRepo := repository.NewMessageRepository(db.DB)

	// Services
	userSvc := services.NewUserService(userRepo, logger)
	msgSvc := services.NewMessageService(msgRepo, userRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userSvc, logger)
	mailHandler := handlers.NewMailHandler(msgSvc)
	adminHandler := handlers.NewAdminHandler(userRepo, logger)
	authHandler := handlers.NewAuthHandler(userRepo, logger)

	// Router setup
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware(logger), middleware.ErrorHandler())

	router.SetupRouter(r, userHandler)
	router.SetupExpandedRoutes(r, mailHandler, adminHandler)
	router.SetupAuthRoutes(r, authHandler)

	return &App{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: r,
		},
		db:     db,
		logger: logger,
	}
}

func (a *App) Run() error {
	a.logger.Info("Starting server", zap.String("port", a.server.Addr))
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) {
	a.logger.Info("Closing database connection")
	a.db.Close()
	a.logger.Sync()
}
