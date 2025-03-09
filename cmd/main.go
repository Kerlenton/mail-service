package main

import (
	"context"
	"log"
	"mail-service/internal"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// Завершаем работу сервера с тайм-аутом
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Shutdown(ctxShutdown)
}
