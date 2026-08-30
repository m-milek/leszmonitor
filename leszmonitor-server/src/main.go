package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/m-milek/leszmonitor/api"
	"github.com/m-milek/leszmonitor/api/controllers"
	config "github.com/m-milek/leszmonitor/appconfig"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/services"
	"github.com/m-milek/leszmonitor/workers"
	"github.com/m-milek/leszmonitor/workers/probesrunner"
	"github.com/m-milek/leszmonitor/workers/resultsprocessor"
)

//go:embed all:static
var staticFiles embed.FS

func runComponents(ctx context.Context, wg *sync.WaitGroup) {
	wg.Go(func() {
		manager := probes.NewManager(db.Get())
		manager.Run(ctx)
	})
	wg.Go(func() {
		workers.StartDataCleanupWorker(ctx)
	})
	wg.Go(func() {
		resultsProcessor := resultsprocessor.NewResultsProcessor(db.Get())
		resultsProcessor.Run(ctx)
	})
}

func main() {
	logger := log.New()

	appCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	appCtx = log.WithContext(appCtx, &logger)

	var wg sync.WaitGroup

	err := config.Validate()
	if err != nil {
		logger.Fatal().Err(err).Msg("Environment variable validation failed")
	}
	logger.Info().Msg("Environment variable validation OK")

	err = db.InitFromEnv(appCtx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize SQLite connection")
	}

	database := db.Get()

	projectService := services.NewProjectService(services.ProjectServiceDeps{
		DB:          database,
		UserService: nil,
	})
	userService := services.NewUserService(services.UserServiceDeps{
		DB:             database,
		ProjectService: projectService,
	})
	projectService.UserService = userService
	monitorService := services.NewMonitorService(services.MonitorServiceDeps{
		DB: database,
	})
	monitorResultService := services.NewMonitorResultsService(services.MonitorResultsServiceDeps{
		DB: database,
	})
	monitorStatsService := services.NewMonitorStatsService(services.MonitorStatsServiceDeps{
		DB: database,
	})
	auditLogService := services.NewAuditLogService(services.AuditLogServiceDeps{
		DB: database,
	})
	instanceMetadataService := services.NewInstanceMetadataService(services.InstanceMetadataServiceDeps{})

	projectAPIController := controllers.NewProjectAPIController(projectService)
	userAPIController := controllers.NewUserAPIController(userService)
	monitorAPIController := controllers.NewMonitorAPIController(monitorService)
	monitorResultsAPIController := controllers.NewMonitorResultsAPIController(monitorResultService)
	monitorStatsAPIController := controllers.NewMonitorStatsAPIController(monitorStatsService)
	auditLogAPIController := controllers.NewAuditLogAPIController(auditLogService)
	instanceMetadataAPIController := controllers.NewInstanceMetadataAPIController(instanceMetadataService)

	authzMiddlewareService := services.NewAuthzMiddlewareService(database)

	handlers := api.Handlers{
		Project:                projectAPIController,
		User:                   userAPIController,
		Monitor:                monitorAPIController,
		MonitorResults:         monitorResultsAPIController,
		MonitorStats:           monitorStatsAPIController,
		AuditLog:               auditLogAPIController,
		InstanceMetadata:       instanceMetadataAPIController,
		AuthzMiddlewareService: authzMiddlewareService,
	}

	// Start the server
	serverConfig := api.DefaultServerConfig()
	logger.Info().Msg("Starting API server...")
	server, done, err := api.StartServer(appCtx, serverConfig, staticFiles, handlers)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to start API server")
		os.Exit(1)
	}
	logger.Info().Msg("API server started successfully")

	runComponents(appCtx, &wg)

	<-appCtx.Done()
	logger.Info().Msg("Shutdown signal received")

	// Create a timeout context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(appCtx, 10*time.Second)
	defer shutdownCancel()

	// Shutdown API server
	logger.Info().Msg("Shutting down API server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("API server shutdown error")
	} else {
		logger.Info().Msg("API server stopped gracefully")
	}
	close(done)

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Info().Msg("All processes terminated successfully")

	database.Close()
	logger.Info().Msg("Database connection closed")
}
