package main

import (
	"context"
	"github.com/m-milek/leszmonitor/api"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/log-capture"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/uptime"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func runComponents(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		uptime.StartUptimeWorker(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log_capture.StartLogWorker(ctx)
	}()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup

	err := env.Validate()
	if err != nil {
		logging.Init.Fatal().Err(err).Msg("Environment variable validation failed")
	}
	logging.Init.Info().Msg("Environment variable validation OK")

	// Initialize logger
	logConfig := logging.GetLoggerConfig()

	err = logging.InitLogging(logConfig)
	if err != nil {
		logging.Main.Fatal().Err(err).Msg("Failed to initialize logger")
	}
	logging.Main.Info().Msg("Logger initialized successfully")

	err = db.InitDBClient(ctx)
	if err != nil {
		logging.Init.Fatal().Err(err).Msg("Failed to initialize MongoDB connection")
	}

	// Start the server
	serverConfig := api.DefaultServerConfig()
	logging.Main.Info().Msg("Starting API server...")
	server, done, err := api.StartServer(serverConfig)
	if err != nil {
		logging.Main.Error().Err(err).Msg("Failed to start API server")
		os.Exit(1)
	}
	logging.Main.Info().Msg("API server started successfully")

	runComponents(ctx, &wg)

	<-ctx.Done()
	logging.Main.Info().Msg("Shutdown signal received")

	// Create a timeout context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown API server
	logging.Main.Info().Msg("Shutting down API server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logging.Main.Error().Err(err).Msg("API server shutdown error")
	} else {
		logging.Main.Info().Msg("API server stopped gracefully")
	}
	close(done)

	// Wait for all goroutines to finish
	wg.Wait()
	logging.Main.Info().Msg("All processes terminated successfully")
}
