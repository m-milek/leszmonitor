package main

import (
	"context"
	"github.com/m-milek/leszmonitor/api"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/logs"
	"github.com/m-milek/leszmonitor/uptime"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup

	var serverConfig = api.DefaultServerConfig()

	// Start the server
	logger.Main.Info().Msg("Starting API server...")
	server, done, err := api.StartServer(serverConfig)
	if err != nil {
		logger.Main.Error().Err(err).Msg("Failed to start API server")
		os.Exit(1)
	}
	logger.Main.Info().Msg("API server started successfully")

	wg.Add(1)
	go func() {
		defer wg.Done()
		uptime.StartUptimeWorker(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logs.StartLogWorker(ctx)
	}()

	<-ctx.Done()
	logger.Main.Info().Msg("Shutdown signal received")

	// Create a timeout context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown API server
	logger.Main.Info().Msg("Shutting down API server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Main.Error().Err(err).Msg("API server shutdown error")
	} else {
		logger.Main.Info().Msg("API server stopped gracefully")
	}
	close(done)

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Main.Info().Msg("All processes terminated successfully")
}
