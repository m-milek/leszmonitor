package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuditLogIntegrationTest(t *testing.T) (context.Context, AuditLogService, *ProjectService, *UserService, *models.User) {
	ctx, projectService, userService, user := setupIntegrationTest(t)

	auditLogService := NewAuditLogService(AuditLogServiceDeps{
		DB:          db.Get(),
		AuthService: NewAuthorizationService(AuthorizationServiceDeps{DB: db.Get()}),
	})

	return ctx, auditLogService, projectService, userService, user
}

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

		userClaims := &auth.UserClaims{
			Username: owner.Username,
		}
		filter := security.AuditLogFilter{
			ProjectID: &project.ID,
		}

		entries, svcErr := auditLogService.GetEntries(ctx, userClaims, filter, util.Pagination{Page: 1, PerPage: 10})
		require.Nil(t, svcErr)
		require.Len(t, entries, 2)
	})

	t.Run("Fails with 403 Forbidden for a project member without admin rights", func(t *testing.T) {
		ctx, auditLogService, projectService, userService, owner := setupAuditLogIntegrationTest(t)

		project, err := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "Project 1"})
		require.Nil(t, err)

		require.Nil(t, userService.RegisterUser(ctx, &UserRegisterPayload{Username: "viewer", Password: "Password123!", PasswordConfirm: "Password123!"}))
		viewer, _ := userService.GetUserByUsername(ctx, "viewer")
		_, repoErr := db.Get().Projects().AddMemberToProject(ctx, project.Slug, &models.ProjectMember{ID: viewer.ID, Role: models.RoleViewer})
		require.NoError(t, repoErr)

		userClaims := &auth.UserClaims{
			Username: viewer.Username,
		}
		filter := security.AuditLogFilter{
			ProjectID: &project.ID,
		}

		entries, svcErr := auditLogService.GetEntries(ctx, userClaims, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusForbidden, svcErr.Code)
		assert.Nil(t, entries)
	})

	t.Run("Fails with 400 Bad Request when non-instance admin tries to query without ProjectID filter", func(t *testing.T) {
		ctx, auditLogService, _, _, user := setupAuditLogIntegrationTest(t)

		userClaims := &auth.UserClaims{
			Username: user.Username,
		}
		filter := security.AuditLogFilter{}

		entries, svcErr := auditLogService.GetEntries(ctx, userClaims, filter, util.Pagination{Page: 1, PerPage: 10})

		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
		assert.Nil(t, entries)
	})
}
