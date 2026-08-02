package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_AuditLogService_Record(t *testing.T) {
	t.Run("Successfully records an audit log entry", func(t *testing.T) {
		ctx, auditLogService, _, _, _ := setupAuditLogIntegrationTest(t)

		user := "testuser"
		projectID := uuid.New()
		resourceID := uuid.New()

		entry := security.AuditLogParams{
			Username:   &user,
			ProjectID:  &projectID,
			ResourceID: &resourceID,
			Action:     security.ActionCreateMonitor,
			IsSuccess:  true,
			Summary:    "Test log entry",
		}

		err := auditLogService.Record(ctx, entry)
		require.NoError(t, err)

		filter := security.AuditLogFilter{ProjectID: &projectID}
		entries, dbErr := db.Get().AuditLog().GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)
		require.Len(t, entries, 1)

		assert.Equal(t, security.ActionCreateMonitor, entries[0].Action)
		assert.Equal(t, "Test log entry", entries[0].Summary)
		assert.Equal(t, "testuser", *entries[0].Username)
	})
}

func TestIntegration_AuditLogService_GetEntries(t *testing.T) {
	t.Run("Successfully retrieves entries for a project admin", func(t *testing.T) {
		ctx, auditLogService, projectService, _, owner := setupAuditLogIntegrationTest(t)

		project, err := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "Project 1"})
		require.Nil(t, err)

		user1 := "user1"
		resource1 := uuid.New()

		auditLogService.Record(ctx, security.AuditLogParams{
			Username:   &user1,
			ProjectID:  &project.ID,
			ResourceID: &resource1,
			Action:     security.ActionCreateMonitor,
			IsSuccess:  true,
			Summary:    "Entry 1",
		})
		auditLogService.Record(ctx, security.AuditLogParams{
			Username:  &owner.Username,
			ProjectID: &project.ID,
			Action:    security.ActionUpdateProject,
			IsSuccess: false,
		})

		filter := security.AuditLogFilter{
			ProjectID: &project.ID,
		}

		entries, svcErr := auditLogService.GetEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.Nil(t, svcErr)
		require.Len(t, entries, 3)
	})
}
