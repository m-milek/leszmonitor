package results_test

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/internal/testsupport"
	"github.com/m-milek/leszmonitor/platform/db"
)

func setupMonitorResultsIntegrationTest(
	t *testing.T,
) (context.Context, *results.MonitorResultsService, *users.UserService, *users.User) {
	ctx, userService, user := testsupport.Setup(t)

	service := results.NewMonitorResultsService(results.MonitorResultsServiceDeps{
		DB: db.Get(),
	})

	return ctx, service, userService, user
}

func insertTestMonitor(ctx context.Context, t *testing.T) *monitors.Monitor {
	return testsupport.InsertTestMonitor(ctx, t)
}
