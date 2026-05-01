package api

import (
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
func createServer(config ServerConfig, staticFiles embed.FS) (*http.Server, error) {
	publicRouter := http.NewServeMux()
	protectedRouter := http.NewServeMux()

	SetupRouters(publicRouter, protectedRouter, staticFiles)

	combinedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if !strings.HasPrefix(path, "/api/") {
			publicRouter.ServeHTTP(w, r)
			return
		}

		publicAPIPaths := []string{
			"/api/auth/register",
			"/api/auth/login",
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

	handler = middleware.Logger(handler)

	server := &http.Server{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	return server, nil
}

// StartServer initializes and starts the HTTP server based on the provided configuration.
func StartServer(config ServerConfig, staticFiles embed.FS) (*http.Server, chan struct{}, error) {
	server, err := createServer(config, staticFiles)

	if err != nil {
		log.Api.Error().Err(err).Msg("Error creating API server")
		return nil, nil, err
	}

	// Create a done channel to signal when server is shut down
	done := make(chan struct{})

	// Start server in a goroutine
	go func() {
		log.Api.Info().Msgf("API server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Api.Fatal().Err(err).Msgf("Error starting API server: %v", err)
		}
	}()

	return server, done, nil
}
