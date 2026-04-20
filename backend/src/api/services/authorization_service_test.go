package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/assert"
)

func setupTestAuthorizationService() (context.Context, *authorizationServiceT, *db.MockDB) {
	mockDB := &db.MockDB{
		UsersRepo:    new(db.MockUserRepository),
		ProjectsRepo: new(db.MockProjectRepository),
	}
	db.Set(mockDB)

	authService := newAuthorizationService()

	return context.Background(), authService, mockDB
}

func createTestUser(username string) *models.User {
	return &models.User{
		ID:       pgtype.UUID{Bytes: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, Valid: true},
		Username: username,
	}
}

func createTestProject(name string, members []models.ProjectMember) *models.Project {
	project := &models.Project{
		ID:      pgtype.UUID{Bytes: [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, Valid: true},
		Members: members,
	}
	project.DisplayIDFromName.Init(name)
	return project
}

// Helper to set up project and user with role
func setupProjectWithUser(username string, role models.Role) (*models.User, *models.Project) {
	user := createTestUser(username)
	project := createTestProject("test-project", []models.ProjectMember{{ID: user.ID, Role: role}})
	return user, project
}

// Helper to mock successful project and user retrieval
func mockProjectAndUser(mockDB *db.MockDB, ctx context.Context, project *models.Project, user *models.User) {
	mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectByDisplayID", ctx, project.DisplayID).Return(project, nil)
	mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, user.Username).Return(user, nil)
}

func TestAuthorizationServiceT_AuthorizeProjectAction(t *testing.T) {
	roleTests := []struct {
		name       string
		username   string
		role       models.Role
		permission models.Permission
	}{
		{"Authorizes owner with all permissions", "owner", models.RoleOwner, models.PermissionProjectReader},
		{"Authorizes admin with edit permissions", "admin", models.RoleAdmin, models.PermissionProjectEditor},
		{"Authorizes member with monitor edit", "member", models.RoleMember, models.PermissionMonitorEditor},
		{"Authorizes viewer with read permissions", "viewer", models.RoleViewer, models.PermissionProjectReader},
	}

	for _, tt := range roleTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, authService, mockDB := setupTestAuthorizationService()
			defer db.Set(nil)

			user, project := setupProjectWithUser(tt.username, tt.role)
			mockProjectAndUser(mockDB, ctx, project, user)

			resultProject, err := authService.authorizeProjectAction(ctx, &middleware.ProjectAuth{
				ProjectID: project.DisplayID,
				Username:  tt.username,
			}, tt.permission)

			assert.Nil(t, err)
			assert.NotNil(t, resultProject)
			assert.Equal(t, project.DisplayID, resultProject.DisplayID)
		})
	}

	t.Run("Fails when project does not exist", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectByDisplayID", ctx, "nonexistent").
			Return((*models.Project)(nil), db.ErrNotFound)

		resultProject, err := authService.authorizeProjectAction(ctx, &middleware.ProjectAuth{
			ProjectID: "nonexistent",
			Username:  "testuser",
		}, models.PermissionProjectReader)

		assert.Nil(t, resultProject)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})

	t.Run("Fails when project retrieval returns database error", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectByDisplayID", ctx, "test-project").
			Return((*models.Project)(nil), errors.New("database error"))

		resultProject, err := authService.authorizeProjectAction(ctx, &middleware.ProjectAuth{
			ProjectID: "test-project",
			Username:  "testuser",
		}, models.PermissionProjectReader)

		assert.Nil(t, resultProject)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
	})

	t.Run("Fails when user does not exist", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		_, project := setupProjectWithUser("testuser", models.RoleOwner)
		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectByDisplayID", ctx, project.DisplayID).Return(project, nil)
		mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, "nonexistent").
			Return((*models.User)(nil), db.ErrNotFound)

		resultProject, err := authService.authorizeProjectAction(ctx, &middleware.ProjectAuth{
			ProjectID: project.DisplayID,
			Username:  "nonexistent",
		}, models.PermissionProjectReader)

		assert.Nil(t, resultProject)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})

	t.Run("Fails when user is not a member of the project", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user := createTestUser("testuser")
		otherUser := createTestUser("otheruser")
		otherUser.ID = pgtype.UUID{Bytes: [16]byte{99, 98, 97, 96, 95, 94, 93, 92, 91, 90, 89, 88, 87, 86, 85, 84}, Valid: true}

		project := createTestProject("test-project", []models.ProjectMember{{ID: otherUser.ID, Role: models.RoleOwner}})
		mockProjectAndUser(mockDB, ctx, project, user)

		resultProject, err := authService.authorizeProjectAction(ctx, &middleware.ProjectAuth{
			ProjectID: project.DisplayID,
			Username:  user.Username,
		}, models.PermissionProjectReader)

		assert.Nil(t, resultProject)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.Code)
		assert.Contains(t, err.Err.Error(), "is not a member")
	})

	permissionFailTests := []struct {
		name       string
		role       models.Role
		permission models.Permission
	}{
		{"Fails when viewer lacks edit permissions", models.RoleViewer, models.PermissionProjectEditor},
		{"Fails when member lacks admin permissions", models.RoleMember, models.PermissionProjectAdmin},
	}

	for _, tt := range permissionFailTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, authService, mockDB := setupTestAuthorizationService()
			defer db.Set(nil)

			user, project := setupProjectWithUser("testuser", tt.role)
			mockProjectAndUser(mockDB, ctx, project, user)

			resultProject, err := authService.authorizeProjectAction(ctx, &middleware.ProjectAuth{
				ProjectID: project.DisplayID,
				Username:  user.Username,
			}, tt.permission)

			assert.Nil(t, resultProject)
			assert.NotNil(t, err)
			assert.Equal(t, http.StatusForbidden, err.Code)
			assert.Contains(t, err.Err.Error(), "does not have required permissions")
		})
	}

	t.Run("Authorizes with multiple permissions when user has all", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user, project := setupProjectWithUser("owner", models.RoleOwner)
		mockProjectAndUser(mockDB, ctx, project, user)

		resultProject, err := authService.authorizeProjectAction(ctx, &middleware.ProjectAuth{
			ProjectID: project.DisplayID,
			Username:  user.Username,
		}, models.PermissionProjectAdmin, models.PermissionMonitorAdmin)

		assert.Nil(t, err)
		assert.NotNil(t, resultProject)
	})

	t.Run("Fails with multiple permissions when user lacks one", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user, project := setupProjectWithUser("admin", models.RoleAdmin)
		mockProjectAndUser(mockDB, ctx, project, user)

		resultProject, err := authService.authorizeProjectAction(ctx, &middleware.ProjectAuth{
			ProjectID: project.DisplayID,
			Username:  user.Username,
		}, models.PermissionProjectEditor, models.PermissionProjectAdmin)

		assert.Nil(t, resultProject)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.Code)
	})
}
