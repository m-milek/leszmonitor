package auditlog_test

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/features/auditlog"
	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/internal/testsupport"
	"github.com/m-milek/leszmonitor/platform/db"
)

func setupAuditLogIntegrationTest(
	t *testing.T,
) (context.Context, auditlog.AuditLogService, *users.UserService, *users.User) {
	ctx, userService, user := testsupport.Setup(t)

	auditLogService := auditlog.NewAuditLogService(auditlog.AuditLogServiceDeps{
		DB: db.Get(),
	})

	return ctx, auditLogService, userService, user
}
