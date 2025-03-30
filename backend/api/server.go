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
		Host:         "localhost",
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

func StartServer(config ServerConfig) {
	log.Api.Info().Msg("Starting server...")
	server, err := createServer(config)

	if err != nil {
		log.Api.Error().Msg("Error creating server: " + err.Error())
		return
	}

	log.Api.Info().Msgf("Server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Api.Error().Msgf("Server error: %v", err)
	}
}
