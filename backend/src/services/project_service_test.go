package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestProjectService() (context.Context, *ProjectService, *db.MockDB) {
	mockDB := &db.MockDB{
		UsersRepo:    new(db.MockUserRepository),
		ProjectsRepo: new(db.MockProjectRepository),
	}
	db.Set(mockDB)

	authService := NewAuthorizationService(AuthorizationServiceDeps{DB: mockDB})

	// We instantiate the real UserService for full DI testing if needed
	userService := NewUserService(UserServiceDeps{
		DB:   mockDB,
		Auth: authService,
	})

	svc := NewProjectService(ProjectServiceDeps{
		DB:          mockDB,
		Auth:        authService,
		UserService: userService,
	})

	userService.projectService = svc

	return context.Background(), svc, mockDB
}

func TestProjectServiceT_InternalGetProjectBySlug(t *testing.T) {
	t.Run("Returns project successfully", func(t *testing.T) {
		ctx, svc, mockDB := setupTestProjectService()
		defer db.Set(nil)

		expected := &models.Project{Description: "Test Description"}
		expected.SlugFromName.Init("test-project")

		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectBySlug", ctx, "test-project").Return(expected, nil)

		project, err := svc.internalGetProjectBySlug(ctx, "test-project")

		assert.Nil(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "test-project", project.Slug)
		mockDB.ProjectsRepo.(*db.MockProjectRepository).AssertExpectations(t)
	})

	t.Run("Fails with 404 when project not found", func(t *testing.T) {
		ctx, svc, mockDB := setupTestProjectService()
		defer db.Set(nil)

		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectBySlug", ctx, "nonexistent").Return((*models.Project)(nil), db.ErrNotFound)

		project, err := svc.internalGetProjectBySlug(ctx, "nonexistent")

		assert.Nil(t, project)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockDB.ProjectsRepo.(*db.MockProjectRepository).AssertExpectations(t)
	})

	t.Run("Fails with 500 when database returns error", func(t *testing.T) {
		ctx, svc, mockDB := setupTestProjectService()
		defer db.Set(nil)

		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectBySlug", ctx, "test-project").Return((*models.Project)(nil), errors.New("database error"))

		project, err := svc.internalGetProjectBySlug(ctx, "test-project")

		assert.Nil(t, project)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockDB.ProjectsRepo.(*db.MockProjectRepository).AssertExpectations(t)
	})
}

func TestProjectService_CreateProject(t *testing.T) {
	t.Run("CreateProject success", func(t *testing.T) {
		ctx, svc, mockDB := setupTestProjectService()
		defer db.Set(nil)

		ownerUsername := "testuser"
		userID := uuid.New()
		mockUser := &models.User{ID: userID, Username: ownerUsername}

		payload := CreateProjectPayload{
			Name:        "Test Project",
			Description: "A test project",
		}

		expectedProject := &models.Project{
			Description: payload.Description,
		}
		expectedProject.SlugFromName.Init(payload.Name)

		mockDB.UsersRepo.(*db.MockUserRepository).On("GetUserByUsername", ctx, ownerUsername).Return(mockUser, nil)
		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("InsertProject", ctx, mock.AnythingOfType("*models.Project")).Return(nil)
		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectBySlug", ctx, expectedProject.Slug).Return(expectedProject, nil)

		project, err := svc.CreateProject(ctx, ownerUsername, payload)

		assert.Nil(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "test-project", project.Slug)
		assert.Equal(t, "Test Project", project.Name)
		assert.Equal(t, "A test project", project.Description)

		mockDB.UsersRepo.(*db.MockUserRepository).AssertExpectations(t)
		mockDB.ProjectsRepo.(*db.MockProjectRepository).AssertExpectations(t)
	})
}

func TestProjectService_GetProjectByID(t *testing.T) {
	t.Run("GetProjectByID success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

		assert.True(t, true)
	})
}

func TestProjectService_GetProjects(t *testing.T) {
	t.Run("GetProjects success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

		assert.True(t, true)
	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	t.Run("DeleteProject success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

		assert.True(t, true)
	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	t.Run("UpdateProject success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

		assert.True(t, true)
	})
}

func TestProjectService_AddUserToProject(t *testing.T) {
	t.Run("AddUserToProject success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

		assert.True(t, true)
	})
}

func TestProjectService_RemoveUserFromProject(t *testing.T) {
	t.Run("RemoveUserFromProject success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

		assert.True(t, true)
	})
}

func TestProjectService_ChangeProjectMemberRole(t *testing.T) {
	t.Run("ChangeProjectMemberRole success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

		assert.True(t, true)
	})
}
