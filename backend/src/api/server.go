package api

import (
	"context"
	"embed"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/util"
	"github.com/rs/cors"
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

// createServer sets up the HTTP server with public and protected routes, applying necessary middleware.
func createServer(ctx context.Context, config ServerConfig, staticFiles embed.FS) (*http.Server, error) {
	publicRouter := http.NewServeMux()
	protectedRouter := http.NewServeMux()

	SetupRouters(publicRouter, protectedRouter, staticFiles)

	logger := log.FromContext(ctx).With().Str("component", "api_server").Logger()
	ctx = log.WithContext(ctx, &logger)

	combinedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if !strings.HasPrefix(path, "/api/") {
			publicRouter.ServeHTTP(w, r)
			return
		}

		publicAPIPaths := []string{
			"/api/v1/auth/register",
			"/api/v1/auth/login",
			"/api/ws",
		}
		if util.SliceContains(publicAPIPaths, path) {
			publicRouter.ServeHTTP(w, r)
			return
		}

		// Apply JWT auth to protected API paths
		middleware.JwtAuth(protectedRouter).ServeHTTP(w, r)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(combinedHandler)

	handler = middleware.Logger(ctx, handler)

	server := &http.Server{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	return server, nil
}

// StartServer initializes and starts the HTTP server based on the provided configuration.
func StartServer(ctx context.Context, config ServerConfig, staticFiles embed.FS) (*http.Server, chan struct{}, error) {
	logger := log.FromContext(ctx).With().Str("component", "api_server").Logger()

	server, err := createServer(ctx, config, staticFiles)
	if err != nil {
		logger.Error().Err(err).Msg("Error creating API server")
		return nil, nil, err
	}

	// Create a done channel to signal when server is shut down
	done := make(chan struct{})

	// Start server in a goroutine
	go func() {
		logger.Info().Msgf("API server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Error starting API server")
		}
	}()

	return server, done, nil
}
