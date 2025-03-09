package internal

import (
	"context"
	"fmt"
	"mail-service/config"
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

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := database.InitDB(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
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
}
