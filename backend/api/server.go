package api

import (
	"errors"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/log"
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
		Host:         "127.0.0.1",
		Port:         "7001",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func createServer(config ServerConfig) (*http.Server, error) {
	router := http.NewServeMux()
	SetupRouter(router)

	handler := middleware.Logger(router)

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
		log.Api.Error().Err(err).Msg("Error creating API server")
		return nil, nil, err
	}

	// Create a done channel to signal when server is shut down
	done := make(chan struct{})

	// Start server in a goroutine
	go func() {
		log.Api.Info().Msgf("API server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Api.Error().Err(err).Msgf("Error starting API server: %v", err)
		}
	}()

	return server, done, nil
}
