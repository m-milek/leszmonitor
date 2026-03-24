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
		OrgsRepo:  new(db.MockOrgRepository),
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

func createTestOrg(name string, members []models.OrgMember) *models.Org {
	org := &models.Org{
		ID:      pgtype.UUID{Bytes: [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, Valid: true},
		Members: members,
	}
	org.DisplayIDFromName.Init(name)
	return org
}

// Helper to set up org and user with role
func setupOrgWithUser(username string, role models.Role) (*models.User, *models.Org) {
	user := createTestUser(username)
	org := createTestOrg("test-org", []models.OrgMember{{ID: user.ID, Role: role}})
	return user, org
}

// Helper to mock successful org and user retrieval
func mockOrgAndUser(mockDB *db.MockDB, ctx context.Context, org *models.Org, user *models.User) {
	mockDB.OrgsRepo.(*db.MockOrgRepository).On("GetOrgByDisplayID", ctx, org.DisplayID).Return(org, nil)
	mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, user.Username).Return(user, nil)
}

func TestAuthorizationServiceT_AuthorizeOrgAction(t *testing.T) {
	// Test successful authorization for different roles
	roleTests := []struct {
		name       string
		username   string
		role       models.Role
		permission models.Permission
	}{
		{"Authorizes owner with all permissions", "owner", models.RoleOwner, models.PermissionOrgReader},
		{"Authorizes admin with edit permissions", "admin", models.RoleAdmin, models.PermissionOrgEditor},
		{"Authorizes member with monitor edit", "member", models.RoleMember, models.PermissionMonitorEditor},
		{"Authorizes viewer with read permissions", "viewer", models.RoleViewer, models.PermissionOrgReader},
	}

	for _, tt := range roleTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, authService, mockDB := setupTestAuthorizationService()
			defer db.Set(nil)

			user, org := setupOrgWithUser(tt.username, tt.role)
			mockOrgAndUser(mockDB, ctx, org, user)

			resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
				OrgID:    org.DisplayID,
				Username: tt.username,
			}, tt.permission)

			assert.Nil(t, err)
			assert.NotNil(t, resultOrg)
			assert.Equal(t, org.DisplayID, resultOrg.DisplayID)
		})
	}

	// Test database error scenarios
	t.Run("Fails when org does not exist", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		mockDB.OrgsRepo.(*db.MockOrgRepository).On("GetOrgByDisplayID", ctx, "nonexistent-org").
			Return((*models.Org)(nil), db.ErrNotFound)

		resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
			OrgID:    "nonexistent-org",
			Username: "testuser",
		}, models.PermissionOrgReader)

		assert.Nil(t, resultOrg)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})

	t.Run("Fails when org retrieval returns database error", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		mockDB.OrgsRepo.(*db.MockOrgRepository).On("GetOrgByDisplayID", ctx, "test-org").
			Return((*models.Org)(nil), errors.New("database error"))

		resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
			OrgID:    "test-org",
			Username: "testuser",
		}, models.PermissionOrgReader)

		assert.Nil(t, resultOrg)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
	})

	// Test user error scenarios
	t.Run("Fails when user does not exist", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		_, org := setupOrgWithUser("testuser", models.RoleOwner)
		mockDB.OrgsRepo.(*db.MockOrgRepository).On("GetOrgByDisplayID", ctx, org.DisplayID).Return(org, nil)
		mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, "nonexistent").
			Return((*models.User)(nil), db.ErrNotFound)

		resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
			OrgID:    org.DisplayID,
			Username: "nonexistent",
		}, models.PermissionOrgReader)

		assert.Nil(t, resultOrg)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})

	t.Run("Fails when user retrieval returns database error", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user, org := setupOrgWithUser("testuser", models.RoleOwner)
		mockDB.OrgsRepo.(*db.MockOrgRepository).On("GetOrgByDisplayID", ctx, org.DisplayID).Return(org, nil)
		mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, user.Username).
			Return((*models.User)(nil), errors.New("database error"))

		resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
			OrgID:    org.DisplayID,
			Username: user.Username,
		}, models.PermissionOrgReader)

		assert.Nil(t, resultOrg)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
	})

	t.Run("Fails when user is not a member of the org", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user := createTestUser("testuser")
		otherUser := createTestUser("otheruser")
		otherUser.ID = pgtype.UUID{Bytes: [16]byte{99, 98, 97, 96, 95, 94, 93, 92, 91, 90, 89, 88, 87, 86, 85, 84}, Valid: true}

		org := createTestOrg("test-org", []models.OrgMember{{ID: otherUser.ID, Role: models.RoleOwner}})
		mockOrgAndUser(mockDB, ctx, org, user)

		resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
			OrgID:    org.DisplayID,
			Username: user.Username,
		}, models.PermissionOrgReader)

		assert.Nil(t, resultOrg)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.Code)
		assert.Contains(t, err.Err.Error(), "is not a member")
	})

	// Test permission failures
	permissionFailTests := []struct {
		name       string
		role       models.Role
		permission models.Permission
	}{
		{"Fails when viewer lacks edit permissions", models.RoleViewer, models.PermissionOrgEditor},
		{"Fails when member lacks admin permissions", models.RoleMember, models.PermissionOrgAdmin},
	}

	for _, tt := range permissionFailTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, authService, mockDB := setupTestAuthorizationService()
			defer db.Set(nil)

			user, org := setupOrgWithUser("testuser", tt.role)
			mockOrgAndUser(mockDB, ctx, org, user)

			resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
				OrgID:    org.DisplayID,
				Username: user.Username,
			}, tt.permission)

			assert.Nil(t, resultOrg)
			assert.NotNil(t, err)
			assert.Equal(t, http.StatusForbidden, err.Code)
			assert.Contains(t, err.Err.Error(), "does not have required permissions")
		})
	}

	// Test multiple permissions
	t.Run("Authorizes with multiple permissions when user has all", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user, org := setupOrgWithUser("owner", models.RoleOwner)
		mockOrgAndUser(mockDB, ctx, org, user)

		resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
			OrgID:    org.DisplayID,
			Username: user.Username,
		}, models.PermissionOrgAdmin, models.PermissionMonitorAdmin)

		assert.Nil(t, err)
		assert.NotNil(t, resultOrg)
	})

	t.Run("Fails with multiple permissions when user lacks one", func(t *testing.T) {
		ctx, authService, mockDB := setupTestAuthorizationService()
		defer db.Set(nil)

		user, org := setupOrgWithUser("admin", models.RoleAdmin)
		mockOrgAndUser(mockDB, ctx, org, user)

		resultOrg, err := authService.authorizeOrgAction(ctx, &middleware.OrgAuth{
			OrgID:    org.DisplayID,
			Username: user.Username,
		}, models.PermissionOrgEditor, models.PermissionOrgAdmin)

		assert.Nil(t, resultOrg)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusForbidden, err.Code)
	})
}
