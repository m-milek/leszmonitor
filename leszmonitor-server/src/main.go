package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/m-milek/leszmonitor/features/auditlog"
	"github.com/m-milek/leszmonitor/features/instance"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/users"

	"github.com/m-milek/leszmonitor/app"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/features/monitors/workers"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/config"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/downtime"
	"github.com/m-milek/leszmonitor/platform/log"
)

//go:embed all:static
var staticFiles embed.FS

func runComponents(ctx context.Context, wg *sync.WaitGroup) {
	wg.Go(func() {
		manager := workers.NewManager(db.Get())
		manager.Run(ctx)
	})
	wg.Go(func() {
		workers.StartDataCleanupWorker(ctx)
	})
	wg.Go(func() {
		resultsProcessor := workers.NewResultsProcessor(db.Get())
		resultsProcessor.Run(ctx)
	})
	wg.Go(func() {
		heartbeatWorker := downtime.NewHeartbeatWorker(db.Get())
		heartbeatWorker.Run(ctx)
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

	userService := users.NewUserService(users.UserServiceDeps{
		DB: database,
	})
	monitorService := monitors.NewMonitorService(monitors.MonitorServiceDeps{
		DB: database,
	})
	monitorResultService := results.NewMonitorResultsService(results.MonitorResultsServiceDeps{
		DB: database,
	})
	monitorStatsService := stats.NewMonitorStatsService(stats.MonitorStatsServiceDeps{
		DB: database,
	})
	auditLogService := auditlog.NewAuditLogService(auditlog.AuditLogServiceDeps{
		DB: database,
	})
	tagService := tags.NewTagService(tags.TagServiceDeps{
		DB: database,
	})
	instanceMetadataService := instance.NewInstanceMetadataService(instance.InstanceMetadataServiceDeps{})

	userAPIController := users.NewUserAPIController(userService)
	monitorAPIController := monitors.NewMonitorAPIController(monitorService)
	monitorResultsAPIController := results.NewMonitorResultsAPIController(monitorResultService)
	monitorStatsAPIController := stats.NewMonitorStatsAPIController(monitorStatsService)
	auditLogAPIController := auditlog.NewAuditLogAPIController(auditLogService)
	tagAPIController := tags.NewTagAPIController(tagService)
	instanceMetadataAPIController := instance.NewInstanceMetadataAPIController(instanceMetadataService)

	authzMiddlewareService := users.NewAuthzMiddlewareService(database)

	handlers := app.Handlers{
		User:                   userAPIController,
		Monitor:                monitorAPIController,
		MonitorResults:         monitorResultsAPIController,
		MonitorStats:           monitorStatsAPIController,
		AuditLog:               auditLogAPIController,
		Tag:                    tagAPIController,
		InstanceMetadata:       instanceMetadataAPIController,
		AuthzMiddlewareService: authzMiddlewareService,
	}

	svcErr := userService.EnsureAdminUserExists(appCtx)
	if svcErr != nil {
		logger.Fatal().Err(svcErr).Msg("Failed to ensure admin user exists")
	}

	// Start the server
	serverConfig := app.DefaultServerConfig()
	logger.Info().Msg("Starting API server...")
	server, done, err := app.StartServer(appCtx, serverConfig, staticFiles, handlers)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to start API server")
		os.Exit(1)
	}
	logger.Info().Msg("API server started successfully")

	runComponents(appCtx, &wg)

	<-appCtx.Done()
	logger.Info().Msg("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(appCtx, 10*time.Second)
	defer shutdownCancel()

	logger.Info().Msg("Shutting down API server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("API server shutdown error")
	} else {
		logger.Info().Msg("API server stopped gracefully")
	}
	close(done)

	wg.Wait()
	logger.Info().Msg("All processes terminated successfully")

	database.Close()
	logger.Info().Msg("Database connection closed")
}
