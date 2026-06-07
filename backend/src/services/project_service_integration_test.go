package services

import (
	"context"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupIntegrationTest initializes a temporary SQLite DB, sets up services, and registers a test user.
func setupIntegrationTest(t *testing.T) (context.Context, *ProjectService, *UserService, *models.User) {
	ctx := context.Background()

	// Use a temporary directory for the sqlite file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "testdb.sqlite")

	// DSN format for SQLite
	dsn := "file:" + dbPath + "?_pragma=foreign_keys(1)"

	realDB, err := db.New(ctx, dsn)
	require.NoError(t, err)

	// Set globally because parts of the system (like authorization) might still rely on db.Get()
	db.Set(realDB)

	t.Cleanup(func() {
		realDB.Close()
		db.Set(nil)
	})

	authService := NewAuthorizationService(AuthorizationServiceDeps{DB: realDB})

	userService := NewUserService(UserServiceDeps{
		DB:   realDB,
		Auth: authService,
	})

	projectService := NewProjectService(ProjectServiceDeps{
		DB:          realDB,
		Auth:        authService,
		UserService: userService,
	})

	userService.projectService = projectService

	// Setup Phase: Create a real user in the DB
	registerPayload := &UserRegisterPayload{
		Username:        "integration_user",
		Password:        "Password123!",
		PasswordConfirm: "Password123!",
	}
	svcErr := userService.RegisterUser(ctx, registerPayload)
	require.Nil(t, svcErr)

	user, svcErr := userService.GetUserByUsername(ctx, "integration_user")
	require.Nil(t, svcErr)
	require.NotNil(t, user)

	return ctx, projectService, userService, user
}

func TestIntegration_ProjectService_CreateProject(t *testing.T) {
	t.Run("Successfully creates a project", func(t *testing.T) {
		ctx, projectService, _, user := setupIntegrationTest(t)

		createPayload := CreateProjectPayload{
			Name:        "Integration Test Project",
			Description: "A real project stored in SQLite",
		}
		project, svcErr := projectService.CreateProject(ctx, user.Username, createPayload)
		require.Nil(t, svcErr)
		require.NotNil(t, project)

		assert.Equal(t, "Integration Test Project", project.Name)
		assert.Equal(t, "integration-test-project", project.Slug)
		assert.Equal(t, "A real project stored in SQLite", project.Description)

		// Verify memberships exist in the database model
		require.Len(t, project.Members, 1)
		assert.Equal(t, user.ID, project.Members[0].ID)
		assert.Equal(t, models.RoleOwner, project.Members[0].Role)

		// Directly verify against the DB layer to ensure it persisted correctly
		dbProject, err := db.Get().Projects().GetProjectBySlug(ctx, project.Slug)
		require.NoError(t, err)
		assert.Equal(t, "integration-test-project", dbProject.Slug)
		require.Len(t, dbProject.Members, 1)
		assert.Equal(t, models.RoleOwner, dbProject.Members[0].Role)
	})

	t.Run("Fails with 404 when user doesn't exist", func(t *testing.T) {
		ctx, projectService, _, _ := setupIntegrationTest(t)

		createPayload := CreateProjectPayload{
			Name:        "Integration Test Project",
			Description: "A real project stored in SQLite",
		}

		project, svcErr := projectService.CreateProject(ctx, "nonexistent_user", createPayload)

		assert.Nil(t, project)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

	t.Run("Fails with 409 Conflict when creating a duplicate project", func(t *testing.T) {
		ctx, projectService, _, user := setupIntegrationTest(t)

		createPayload := CreateProjectPayload{
			Name:        "Integration Test Project",
			Description: "A real project stored in SQLite",
		}

		_, svcErr := projectService.CreateProject(ctx, user.Username, createPayload)
		require.Nil(t, svcErr)

		// Attempt to create it again to trigger UNIQUE constraint in SQLite
		duplicateProject, duplicateErr := projectService.CreateProject(ctx, user.Username, createPayload)

		assert.Nil(t, duplicateProject)
		require.NotNil(t, duplicateErr)
		assert.Equal(t, http.StatusConflict, duplicateErr.Code)
	})
}

func TestIntegration_ProjectService_GetProjects(t *testing.T) {
	t.Run("Returns all projects for a user when no query is provided", func(t *testing.T) {
		ctx, projectService, _, user1 := setupIntegrationTest(t)

		p1, err := projectService.CreateProject(ctx, user1.Username, CreateProjectPayload{Name: "P1", Description: "D1"})
		require.Nil(t, err)

		p2, err := projectService.CreateProject(ctx, user1.Username, CreateProjectPayload{Name: "P2", Description: "D2"})
		require.Nil(t, err)

		projects, err := projectService.GetProjects(ctx, user1.Username, "")
		require.Nil(t, err)
		require.Len(t, projects, 3) // Includes auto-generated Sandbox project

		var slugs []string
		for _, p := range projects {
			slugs = append(slugs, p.Slug)
		}
		assert.Contains(t, slugs, p1.Slug)
		assert.Contains(t, slugs, p2.Slug)
		assert.Contains(t, slugs, "integrationusers-sandbox")
	})

	t.Run("Filters projects by shared member username", func(t *testing.T) {
		ctx, projectService, userService, user1 := setupIntegrationTest(t)

		require.Nil(t, userService.RegisterUser(ctx, &UserRegisterPayload{Username: "user2", Password: "Password123!", PasswordConfirm: "Password123!"}))
		user2, _ := userService.GetUserByUsername(ctx, "user2")

		_, err := projectService.CreateProject(ctx, user1.Username, CreateProjectPayload{Name: "Only User1", Description: "D1"})
		require.Nil(t, err)

		p2, err := projectService.CreateProject(ctx, user1.Username, CreateProjectPayload{Name: "Shared Project", Description: "D2"})
		require.Nil(t, err)

		_, repoErr := db.Get().Projects().AddMemberToProject(ctx, p2.Slug, &models.ProjectMember{ID: user2.ID, Role: models.RoleViewer})
		require.NoError(t, repoErr)

		_, err = projectService.CreateProject(ctx, user2.Username, CreateProjectPayload{Name: "Only User2", Description: "D3"})
		require.Nil(t, err)

		projects, err := projectService.GetProjects(ctx, user1.Username, "user2")
		require.Nil(t, err)

		require.Len(t, projects, 1)
		assert.Equal(t, p2.Slug, projects[0].Slug)
	})

	t.Run("Returns 404 when requestor user is not found", func(t *testing.T) {
		ctx, projectService, _, _ := setupIntegrationTest(t)

		projects, err := projectService.GetProjects(ctx, "nonexistent_user", "")
		assert.Nil(t, projects)
		require.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
	})
}

func TestIntegration_ProjectService_GetProjectByID(t *testing.T) {
	t.Run("Successfully retrieves a project by ID", func(t *testing.T) {
		ctx, projectService, _, user := setupIntegrationTest(t)

		createPayload := CreateProjectPayload{
			Name:        "Integration Test Project",
			Description: "A real project stored in SQLite",
		}
		project, svcErr := projectService.CreateProject(ctx, user.Username, createPayload)
		require.Nil(t, svcErr)
		require.NotNil(t, project)

		projectAuth := authorization.ProjectAuthorization{ProjectID: project.ID, Username: "integration_user"}

		retrievedProject, getErr := projectService.GetProjectByID(ctx, &projectAuth)

		assert.Nil(t, getErr)
		require.NotNil(t, retrievedProject)
		assert.Equal(t, project.ID, retrievedProject.ID)
		assert.Equal(t, "Integration Test Project", retrievedProject.Name)
		assert.Equal(t, "A real project stored in SQLite", retrievedProject.Description)
	})

	t.Run("Fails with 404 when project doesn't exist", func(t *testing.T) {
		ctx, projectService, _, _ := setupIntegrationTest(t)

		projectAuth := authorization.ProjectAuthorization{ProjectID: uuid.New(), Username: "integration_user"}

		retrievedProject, getErr := projectService.GetProjectByID(ctx, &projectAuth)

		assert.Nil(t, retrievedProject)
		require.NotNil(t, getErr)
		assert.Equal(t, http.StatusNotFound, getErr.Code)
	})
}

func TestIntegration_ProjectService_DeleteProject(t *testing.T) {
	t.Run("Successfully deletes a project when user is admin", func(t *testing.T) {
		ctx, projectService, _, user := setupIntegrationTest(t)

		project, err := projectService.CreateProject(ctx, user.Username, CreateProjectPayload{Name: "To Be Deleted"})
		require.Nil(t, err)

		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  user.Username,
		}

		svcErr := projectService.DeleteProject(ctx, auth)
		require.Nil(t, svcErr)

		// Verify it's actually deleted from DB
		deletedProject, dbErr := db.Get().Projects().GetProjectBySlug(ctx, project.Slug)
		require.NotNil(t, dbErr)
		assert.ErrorIs(t, dbErr, db.ErrNotFound)
		assert.Nil(t, deletedProject)
	})

	t.Run("Fails with 403 Forbidden when user does not have admin permissions", func(t *testing.T) {
		ctx, projectService, userService, owner := setupIntegrationTest(t)

		// Setup: Create a project
		project, err := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "Not Yours"})
		require.Nil(t, err)

		// Setup: Register a second user
		require.Nil(t, userService.RegisterUser(ctx, &UserRegisterPayload{Username: "viewer", Password: "Password123!", PasswordConfirm: "Password123!"}))
		viewer, _ := userService.GetUserByUsername(ctx, "viewer")

		// Setup: Add viewer to project as a Viewer
		_, repoErr := db.Get().Projects().AddMemberToProject(ctx, project.Slug, &models.ProjectMember{ID: viewer.ID, Role: models.RoleViewer})
		require.NoError(t, repoErr)

		// Action: Attempt deletion as viewer
		auth := &authorization.ProjectAuthorization{
			ProjectID: project.ID,
			Username:  viewer.Username,
		}
		svcErr := projectService.DeleteProject(ctx, auth)

		// Assert: Verify forbidden error
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusForbidden, svcErr.Code)

		// Assert: Verify project still exists in DB
		dbProject, dbErr := db.Get().Projects().GetProjectBySlug(ctx, project.Slug)
		require.NoError(t, dbErr)
		assert.NotNil(t, dbProject)
	})

	t.Run("Returns 404 Not Found when project does not exist", func(t *testing.T) {
		ctx, projectService, _, user := setupIntegrationTest(t)

		auth := &authorization.ProjectAuthorization{
			ProjectID: uuid.New(),
			Username:  user.Username,
		}
		svcErr := projectService.DeleteProject(ctx, auth)

		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})
}
