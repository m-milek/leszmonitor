package main

import (
	"context"
	"github.com/m-milek/leszmonitor/api"
	"github.com/m-milek/leszmonitor/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var serverConfig = api.DefaultServerConfig()

	// Start the server
	log.Main.Info().Msg("Starting API server...")
	server, done, err := api.StartServer(serverConfig)
	if err != nil {
		log.Main.Error().Err(err).Msg("Failed to start API server")
		os.Exit(1)
	} else {
		log.Main.Info().Msg("API server started successfully")
	}

	// Create a channel to listen for interrupt signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-signalCh
	log.Main.Info().Msg("Shutdown signal received")

	// Create a timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	log.Main.Info().Msg("Shutting down API server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Main.Error().Err(err).Msg("API server shutdown error")
		os.Exit(1)
	}

	// Signal that we're done
	close(done)
	log.Main.Info().Msg("API server stopped gracefully")
}
