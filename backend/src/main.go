package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/m-milek/leszmonitor/api"
	"github.com/m-milek/leszmonitor/config"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/workers"
)

func runComponents(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		workers.StartUptimeWorker(ctx)
	}()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup

	err := config.Validate()
	if err != nil {
		log.Init.Fatal().Err(err).Msg("Environment variable validation failed")
	}
	log.Init.Info().Msg("Environment variable validation OK")

	logConfig := log.GetLoggerConfig()

	err = log.InitLogging(logConfig)
	if err != nil {
		log.Main.Fatal().Err(err).Msg("Failed to initialize logger")
	}
	log.Main.Info().Msg("Logger initialized successfully")

	err = db.InitFromEnv(ctx)
	if err != nil {
		log.Init.Fatal().Err(err).Msg("Failed to initialize PostgreSQL connection")
	}

	// Start the server
	serverConfig := api.DefaultServerConfig()
	log.Main.Info().Msg("Starting API server...")
	server, done, err := api.StartServer(serverConfig)
	if err != nil {
		log.Main.Error().Err(err).Msg("Failed to start API server")
		os.Exit(1)
	}
	log.Main.Info().Msg("API server started successfully")

	runComponents(ctx, &wg)

	<-ctx.Done()
	log.Main.Info().Msg("Shutdown signal received")

	// Create a timeout context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown API server
	log.Main.Info().Msg("Shutting down API server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Main.Error().Err(err).Msg("API server shutdown error")
	} else {
		log.Main.Info().Msg("API server stopped gracefully")
	}
	close(done)

	// Wait for all goroutines to finish
	wg.Wait()
	log.Main.Info().Msg("All processes terminated successfully")
}
