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
		UsersDAO:    new(db.MockUserDAO),
		ProjectsDAO: new(db.MockProjectDAO),
		AuditLogDAO: new(db.MockAuditLogDAO),
	}
	db.Set(mockDB)

	// We instantiate the real UserService for full DI testing if needed
	userService := NewUserService(UserServiceDeps{
		DB: mockDB,
	})

	svc := NewProjectService(ProjectServiceDeps{
		DB:          mockDB,
		UserService: userService,
	})

	userService.projectService = svc

	return context.Background(), svc, mockDB
}

func TestProjectServiceT_InternalGetProjectBySlug(t *testing.T) {
	t.Run("Returns project successfully", func(t *testing.T) {
		ctx, _, mockDB := setupTestProjectService()
		defer db.Set(nil)

		expected := &models.Project{Description: "Test Description"}
		expected.SlugFromName.Init("test-project")

		mockDB.ProjectsDAO.(*db.MockProjectDAO).On("GetProjectBySlug", ctx, "test-project").
			Return(expected, nil)

		project, err := internalGetProjectBySlug(ctx, mockDB, "test-project")

		assert.Nil(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "test-project", project.Slug)
		mockDB.ProjectsDAO.(*db.MockProjectDAO).AssertExpectations(t)
	})

	t.Run("Fails with 404 when project not found", func(t *testing.T) {
		ctx, _, mockDB := setupTestProjectService()
		defer db.Set(nil)

		mockDB.ProjectsDAO.(*db.MockProjectDAO).On("GetProjectBySlug", ctx, "nonexistent").
			Return((*models.Project)(nil), db.ErrNotFound)

		project, err := internalGetProjectBySlug(ctx, mockDB, "nonexistent")

		assert.Nil(t, project)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockDB.ProjectsDAO.(*db.MockProjectDAO).AssertExpectations(t)
	})

	t.Run("Fails with 500 when database returns error", func(t *testing.T) {
		ctx, _, mockDB := setupTestProjectService()
		defer db.Set(nil)

		mockDB.ProjectsDAO.(*db.MockProjectDAO).On("GetProjectBySlug", ctx, "test-project").
			Return((*models.Project)(nil), errors.New("database error"))

		project, err := internalGetProjectBySlug(ctx, mockDB, "test-project")

		assert.Nil(t, project)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockDB.ProjectsDAO.(*db.MockProjectDAO).AssertExpectations(t)
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

		mockDB.UsersDAO.(*db.MockUserDAO).On("GetUserByUsername", ctx, ownerUsername).Return(mockUser, nil)
		mockDB.ProjectsDAO.(*db.MockProjectDAO).On("InsertProject", ctx, mock.AnythingOfType("*models.Project")).
			Return(nil)
		mockDB.ProjectsDAO.(*db.MockProjectDAO).On("GetProjectBySlug", ctx, expectedProject.Slug).
			Return(expectedProject, nil)
		mockDB.AuditLogDAO.(*db.MockAuditLogDAO).On("Record", ctx, mock.AnythingOfType("security.AuditLogParams")).
			Return(nil)

		project, err := svc.CreateProject(ctx, ownerUsername, payload)

		assert.Nil(t, err)
		assert.NotNil(t, project)
		assert.Equal(t, "test-project", project.Slug)
		assert.Equal(t, "Test Project", project.Name)
		assert.Equal(t, "A test project", project.Description)

		mockDB.UsersDAO.(*db.MockUserDAO).AssertExpectations(t)
		mockDB.ProjectsDAO.(*db.MockProjectDAO).AssertExpectations(t)
		mockDB.AuditLogDAO.(*db.MockAuditLogDAO).AssertExpectations(t)
	})
}

func TestProjectService_GetProjectByID(t *testing.T) {
	t.Run("GetProjectByID success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

	})
}

func TestProjectService_GetProjects(t *testing.T) {
	t.Run("GetProjects success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	t.Run("DeleteProject success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

	})
}

func TestProjectService_UpdateProject(t *testing.T) {
	t.Run("UpdateProject success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

	})
}

func TestProjectService_AddUserToProject(t *testing.T) {
	t.Run("AddUserToProject success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

	})
}

func TestProjectService_RemoveUserFromProject(t *testing.T) {
	t.Run("RemoveUserFromProject success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

	})
}

func TestProjectService_ChangeProjectMemberRole(t *testing.T) {
	t.Run("ChangeProjectMemberRole success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestProjectService()
		// defer db.Set(nil)

	})
}
