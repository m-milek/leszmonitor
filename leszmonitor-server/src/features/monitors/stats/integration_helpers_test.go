package stats_test

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/internal/testsupport"
	"github.com/m-milek/leszmonitor/platform/db"
)

func setupMonitorStatsIntegrationTest(
	t *testing.T,
) (context.Context, *stats.MonitorStatsService, *users.UserService, *users.User) {
	ctx, userService, user := testsupport.Setup(t)

	monitorStatsService := stats.NewMonitorStatsService(stats.MonitorStatsServiceDeps{
		DB: db.Get(),
	})

	return ctx, &monitorStatsService, userService, user
}

func insertTestMonitor(ctx context.Context, t *testing.T) *monitors.Monitor {
	return testsupport.InsertTestMonitor(ctx, t)
}
