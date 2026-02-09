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
		UsersRepo: new(db.MockUserRepository),
		TeamsRepo: new(db.MockTeamRepository),
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

func createTestTeam(name string, members []models.TeamMember) *models.Team {
	team := &models.Team{
		ID:      pgtype.UUID{Bytes: [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, Valid: true},
		Members: members,
	}
	team.DisplayIDFromName.Init(name)
	return team
}

// Helper to setup team and user with role
func setupTeamWithUser(username string, role models.Role) (*models.User, *models.Team) {
	user := createTestUser(username)
	team := createTestTeam("test-team", []models.TeamMember{{ID: user.ID, Role: role}})
	return user, team
}

// Helper to mock successful team and user retrieval
func mockTeamAndUser(mockDB *db.MockDB, ctx context.Context, team *models.Team, user *models.User) {
	mockDB.TeamsRepo.(*db.MockTeamRepository).On("GetTeamByDisplayID", ctx, team.DisplayID).Return(team, nil)
	mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, user.Username).Return(user, nil)
}

func TestAuthorizationServiceT_AuthorizeTeamAction(t *testing.T) {
	// Test successful authorization for different roles
	roleTests := []struct {
		name       string
		username   string
		role       models.Role
		permission models.Permission
	}{
		{"Authorizes owner with all permissions", "owner", models.RoleOwner, models.PermissionTeamReader},
		{"Authorizes admin with edit permissions", "admin", models.RoleAdmin, models.PermissionTeamEditor},
		{"Authorizes member with monitor edit", "member", models.RoleMember, models.PermissionMonitorEditor},
		{"Authorizes viewer with read permissions", "viewer", models.RoleViewer, models.PermissionTeamReader},
	}

	for _, tt := range roleTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, authService, mockDB := setupTestAuthorizationService()
			defer db.Set(nil)

			user, team := setupTeamWithUser(tt.username, tt.role)
			mockTeamAndUser(mockDB, ctx, team, user)

			resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
				TeamID:   team.DisplayID,
				Username: tt.username,
			}, tt.permission)

			assert.Nil(t, err)
			assert.NotNil(t, resultTeam)
			assert.Equal(t, team.DisplayID, resultTeam.DisplayID)
		})
	}

	// Test database error scenarios
	t.Run("Fails when team does not exist", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		mockDB.TeamsRepo.(*db.MockTeamRepository).On("GetTeamByDisplayID", ctx, "nonexistent-team").
			Return((*models.Team)(nil), db.ErrNotFound)

		resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
			TeamID:   "nonexistent-team",
			Username: "testuser",
		}, models.PermissionTeamReader)

		assert.Nil(t, resultTeam)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})

	t.Run("Fails when team retrieval returns database error", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		mockDB.TeamsRepo.(*db.MockTeamRepository).On("GetTeamByDisplayID", ctx, "test-team").
			Return((*models.Team)(nil), errors.New("database error"))

		resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
			TeamID:   "test-team",
			Username: "testuser",
		}, models.PermissionTeamReader)

		assert.Nil(t, resultTeam)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
	})

	// Test user error scenarios
	t.Run("Fails when user does not exist", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		_, team := setupTeamWithUser("testuser", models.RoleOwner)
		mockDB.TeamsRepo.(*db.MockTeamRepository).On("GetTeamByDisplayID", ctx, team.DisplayID).Return(team, nil)
		mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, "nonexistent").
			Return((*models.User)(nil), db.ErrNotFound)

		resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
			TeamID:   team.DisplayID,
			Username: "nonexistent",
		}, models.PermissionTeamReader)

		assert.Nil(t, resultTeam)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})

	t.Run("Fails when user retrieval returns database error", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user, team := setupTeamWithUser("testuser", models.RoleOwner)
		mockDB.TeamsRepo.(*db.MockTeamRepository).On("GetTeamByDisplayID", ctx, team.DisplayID).Return(team, nil)
		mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, user.Username).
			Return((*models.User)(nil), errors.New("database error"))

		resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
			TeamID:   team.DisplayID,
			Username: user.Username,
		}, models.PermissionTeamReader)

		assert.Nil(t, resultTeam)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
	})

	t.Run("Fails when user is not a member of the team", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user := createTestUser("testuser")
		otherUser := createTestUser("otheruser")
		otherUser.ID = pgtype.UUID{Bytes: [16]byte{99, 98, 97, 96, 95, 94, 93, 92, 91, 90, 89, 88, 87, 86, 85, 84}, Valid: true}

		team := createTestTeam("test-team", []models.TeamMember{{ID: otherUser.ID, Role: models.RoleOwner}})
		mockTeamAndUser(mockDB, ctx, team, user)

		resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
			TeamID:   team.DisplayID,
			Username: user.Username,
		}, models.PermissionTeamReader)

		assert.Nil(t, resultTeam)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.Code)
		assert.Contains(t, err.Err.Error(), "is not a member of team")
	})

	// Test permission failures
	permissionFailTests := []struct {
		name       string
		role       models.Role
		permission models.Permission
	}{
		{"Fails when viewer lacks edit permissions", models.RoleViewer, models.PermissionTeamEditor},
		{"Fails when member lacks admin permissions", models.RoleMember, models.PermissionTeamAdmin},
	}

	for _, tt := range permissionFailTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, authService, mockDB := setupTestAuthorizationService()
			defer db.Set(nil)

			user, team := setupTeamWithUser("testuser", tt.role)
			mockTeamAndUser(mockDB, ctx, team, user)

			resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
				TeamID:   team.DisplayID,
				Username: user.Username,
			}, tt.permission)

			assert.Nil(t, resultTeam)
			assert.NotNil(t, err)
			assert.Equal(t, http.StatusForbidden, err.Code)
			assert.Contains(t, err.Err.Error(), "does not have required permissions")
		})
	}

	// Test multiple permissions
	t.Run("Authorizes with multiple permissions when user has all", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user, team := setupTeamWithUser("owner", models.RoleOwner)
		mockTeamAndUser(mockDB, ctx, team, user)

		resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
			TeamID:   team.DisplayID,
			Username: user.Username,
		}, models.PermissionTeamAdmin, models.PermissionMonitorAdmin)

		assert.Nil(t, err)
		assert.NotNil(t, resultTeam)
	})

	t.Run("Fails with multiple permissions when user lacks one", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user, team := setupTeamWithUser("admin", models.RoleAdmin)
		mockTeamAndUser(mockDB, ctx, team, user)

		resultTeam, err := authService.authorizeTeamAction(ctx, &middleware.TeamAuth{
			TeamID:   team.DisplayID,
			Username: user.Username,
		}, models.PermissionTeamEditor, models.PermissionTeamAdmin)

		assert.Nil(t, resultTeam)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.Code)
	})
}
