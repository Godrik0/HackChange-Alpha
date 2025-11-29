package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/config"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/container"
)

// @title HackChange-Alpha API
// @version 1.0
// @description API для управления клиентами и расчета ML-скоринга
// @license.name MIT
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app, err := container.New(cfg, nil)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Close()

	app.Logger.Info("Application initialized successfully",
		"version", "1.0.0",
		"log_level", cfg.Log.Level,
		"log_format", cfg.Log.Format,
	)

	go func() {
		if err := app.HTTPServer.Start(); err != nil {
			app.Logger.Error("Server failed to start", "error", err)
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	app.Logger.Info("Server started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Received shutdown signal, gracefully shutting down...")

	shutdownTimeout := time.Duration(cfg.Server.ShutdownTimeout) * time.Second
	if shutdownTimeout == 0 {
		shutdownTimeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.HTTPServer.Shutdown(ctx); err != nil {
		app.Logger.Error("Server forced to shutdown", "error", err)
	}

	app.Logger.Info("Server stopped gracefully")
}
