package internal

import (
	"context"
	"fmt"
	"mail-service/config"
	"mail-service/internal/auth"
	"mail-service/internal/database"
	"mail-service/internal/handlers"
	"mail-service/internal/mail"
	"mail-service/internal/middleware"
	"mail-service/internal/repository"
	"mail-service/internal/router"
	"mail-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	server      *http.Server
	db          *database.Database
	logger      *zap.Logger
	mailService *mail.MailService
}

func NewApp() *App {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Set the JWT key from config.
	auth.SetJWTKey(cfg.Auth.JwtSecret)

	logger, _ := zap.NewProduction()

	db, err := database.InitDB(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}

	// Initialize the mail service
	mailSvc, err := mail.NewMailService()
	if err != nil {
		logger.Error("Failed to initialize mail service", zap.Error(err))
		panic(err)
	}

	repo := repository.NewUserRepository(db.DB)
	service := services.NewUserService(repo, logger)
	handler := handlers.NewUserHandler(service, logger)

	r := gin.Default()
	r.Use(middleware.LoggerMiddleware(logger), middleware.ErrorHandler())

	router.SetupRouter(r, handler)

	return &App{
		server: &http.Server{
			Addr:    ":" + fmt.Sprintf("%d", cfg.Server.Port),
			Handler: r,
		},
		db:          db,
		logger:      logger,
		mailService: mailSvc,
	}
}

func (a *App) Run() error {
	a.logger.Info("Starting server", zap.String("port", a.server.Addr))
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) {
	a.logger.Info("Closing database connection")
	a.db.Close()
	a.mailService.Close()
	a.logger.Sync()
}
