package app

import (
	"context"
	"embed"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/logdyhq/logdy-core/logdy"
	appconfig "github.com/m-milek/leszmonitor/platform/config"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/platform/log"
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
		Port:         os.Getenv(appconfig.APIPort),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// createServer sets up the HTTP server with public and protected routes, applying necessary middleware.
func createServer(
	ctx context.Context,
	config ServerConfig,
	staticFiles embed.FS,
	handlers Handlers,
) (*http.Server, error) {
	router := http.NewServeMux()

	SetupRouters(router, staticFiles, handlers)

	logger := log.FromContext(ctx).With().Str("component", "api_server").Logger()
	ctx = log.WithContext(ctx, &logger)

	if strings.ToLower(os.Getenv(appconfig.EnableLogdy)) != "false" {
		logger.Info().Msg("Initializing Logdy at /_logdy")
		logdyMux := http.NewServeMux()
		logdyLogger := logdy.InitializeLogdy(logdy.Config{
			HttpPathPrefix: "_logdy",
			UiPass:         os.Getenv(appconfig.InstanceAdminPassword),
			ConfigFilePath: os.Getenv(appconfig.LogdyConfigFilePath),
		}, logdyMux)
		log.LogdyGlobalWriter.Logger = logdyLogger
		router.Handle("/_logdy", logdyMux)
		router.Handle("/_logdy/", logdyMux)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", constants.HTTPHeaderContentType},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	handler = httpx.Recoverer(ctx, handler)
	handler = httpx.Logger(ctx, handler)

	server := &http.Server{
		Addr:         net.JoinHostPort(config.Host, config.Port),
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	return server, nil
}

// StartServer initializes and starts the HTTP server based on the provided configuration.
func StartServer(
	ctx context.Context,
	config ServerConfig,
	staticFiles embed.FS,
	handlers Handlers,
) (*http.Server, chan struct{}, error) {
	logger := log.FromContext(ctx).With().Str("component", "api_server").Logger()

	server, err := createServer(ctx, config, staticFiles, handlers)
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
