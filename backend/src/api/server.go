package api

import (
	"errors"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/logger"
	"net"
	"net/http"
	"time"
)

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Host:         "0.0.0.0",
		Port:         "7001",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func createServer(config ServerConfig) (*http.Server, error) {
	publicRouter := http.NewServeMux()
	protectedRouter := http.NewServeMux()

	SetupRouters(publicRouter, protectedRouter)

	combinedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Check if this is a public path
		if path == "/auth/register" || path == "/auth/login" {
			publicRouter.ServeHTTP(w, r)
			return
		}

		// Apply JWT auth to all other paths
		middleware.JwtAuth(protectedRouter).ServeHTTP(w, r)
	})

	handler := middleware.Logger(combinedHandler)

	server := &http.Server{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	return server, nil
}

func StartServer(config ServerConfig) (*http.Server, chan struct{}, error) {
	server, err := createServer(config)

	if err != nil {
		logger.Api.Error().Err(err).Msg("Error creating API server")
		return nil, nil, err
	}

	// Create a done channel to signal when server is shut down
	done := make(chan struct{})

	// Start server in a goroutine
	go func() {
		logger.Api.Info().Msgf("API server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Api.Fatal().Err(err).Msgf("Error starting API server: %v", err)
		}
	}()

	return server, done, nil
}
