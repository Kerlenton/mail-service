// @title Mail Service API
// @version 1.0.0
// @description API for handling mail operations.
// @contact.name Semyon Usachev
// @contact.email semyonschv@yandex.ru
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mail-service/internal"
)

func main() {
	app := internal.NewApp()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	// Shutdown server with timeout
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.Shutdown(ctxShutdown)
}
