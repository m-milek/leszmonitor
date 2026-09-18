package monitors_test

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/internal/testsupport"
	"github.com/m-milek/leszmonitor/platform/db"
)

func setupMonitorIntegrationTest(
	t *testing.T,
) (context.Context, *monitors.MonitorService, *users.UserService, *users.User) {
	ctx, userService, user := testsupport.Setup(t)

	monitorService := monitors.NewMonitorService(monitors.MonitorServiceDeps{
		DB: db.Get(),
	})

	return ctx, monitorService, userService, user
}

func insertTestMonitor(ctx context.Context, t *testing.T) *monitors.Monitor {
	return testsupport.InsertTestMonitor(ctx, t)
}
