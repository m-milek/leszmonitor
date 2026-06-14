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

		entry := security.AuditLogEntry{
			Username:  &user,
			ProjectID: &projectID,
			Action:    security.ActionCreateProject,
			IsSuccess: true,
			Summary:   "Created project",
		}

		err := auditLogService.Record(ctx, entry)
		require.NoError(t, err)

		filter := security.AuditLogFilter{ProjectID: &projectID}
		entries, dbErr := db.Get().AuditLog().GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)
		require.Len(t, entries, 1)

		assert.Equal(t, security.ActionCreateProject, entries[0].Action)
		assert.Equal(t, "Created project", entries[0].Summary)
		assert.Equal(t, "testuser", *entries[0].Username)
		assert.NotEqual(t, uuid.Nil, entries[0].ID)
		assert.NotZero(t, entries[0].CreatedAt)
	})
}

func TestIntegration_AuditLogService_GetEntries(t *testing.T) {
	t.Run("Successfully retrieves entries for a project admin", func(t *testing.T) {
		ctx, auditLogService, projectService, _, owner := setupAuditLogIntegrationTest(t)

		project, err := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "Project 1"})
		require.Nil(t, err)

		auditLogService.Record(ctx, security.AuditLogEntry{
			Username:  &owner.Username,
			ProjectID: &project.ID,
			Action:    security.ActionCreateProject,
			IsSuccess: true,
		})
		auditLogService.Record(ctx, security.AuditLogEntry{
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
		require.Len(t, entries, 2)
	})
}
