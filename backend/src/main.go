package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/m-milek/leszmonitor/api"
	"github.com/m-milek/leszmonitor/appconfig"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/workers"
)

//go:embed all:static
var staticFiles embed.FS

func runComponents(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		workers.StartUptimeWorker(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		workers.StartDataCleanupWorker(ctx)
	}()
}

func main() {
	logger := log.New()

	appCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	appCtx = log.WithContext(appCtx, &logger)

	var wg sync.WaitGroup

	err := config.Validate()
	if err != nil {
		logger.Fatal().Err(err).Msg("Environment variable validation failed")
	}
	logger.Info().Msg("Environment variable validation OK")

	err = db.InitFromEnv(appCtx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize SQLite connection")
	}

	// Start the server
	serverConfig := api.DefaultServerConfig()
	logger.Info().Msg("Starting API server...")
	server, done, err := api.StartServer(appCtx, serverConfig, staticFiles)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to start API server")
		os.Exit(1)
	}
	logger.Info().Msg("API server started successfully")

	runComponents(appCtx, &wg)

	<-appCtx.Done()
	logger.Info().Msg("Shutdown signal received")

	// Create a timeout context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown API server
	logger.Info().Msg("Shutting down API server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("API server shutdown error")
	} else {
		logger.Info().Msg("API server stopped gracefully")
	}
	close(done)

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Info().Msg("All processes terminated successfully")
}
