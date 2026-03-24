package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/assert"
)

func setupTestProjectService() (context.Context, *ProjectServiceT, *db.MockDB) {
	mockDB := &db.MockDB{
		ProjectsRepo: new(db.MockProjectRepository),
	}
	db.Set(mockDB)

	base := newBaseService(nil, "ProjectServiceTest")
	orgService := newProjectService(base)

	ctx := context.Background()

	return ctx, orgService, mockDB
}

func TestProjectServiceT_InternalGetProjectByID(t *testing.T) {
	t.Run("Returns org successfully", func(t *testing.T) {
		ctx, orgService, mockDB := setupTestProjectService()
		defer db.Set(nil)

		expectedProject := &models.Project{
			Description: "Test Description",
		}
		expectedProject.DisplayIDFromName.Init("test-org")

		mockProjectRepo := mockDB.ProjectsRepo.(*db.MockProjectRepository)
		mockProjectRepo.On("GetProjectByDisplayID", ctx, "test-org").Return(expectedProject, nil)

		org, err := orgService.internalGetProjectByDisplayID(ctx, "test-org")

		assert.Nil(t, err)
		assert.NotNil(t, org)
		assert.Equal(t, "test-org", org.DisplayID)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("Fails when project not found", func(t *testing.T) {
		ctx, orgService, mockDB := setupTestProjectService()
		defer db.Set(nil)

		mockProjectRepo := mockDB.ProjectsRepo.(*db.MockProjectRepository)
		mockProjectRepo.On("GetProjectByDisplayID", ctx, "nonexistent").Return((*models.Project)(nil), db.ErrNotFound)

		org, err := orgService.internalGetProjectByDisplayID(ctx, "nonexistent")

		assert.Nil(t, org)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockProjectRepo.AssertExpectations(t)
	})

	t.Run("Fails when database returns error", func(t *testing.T) {
		ctx, orgService, mockDB := setupTestProjectService()
		defer db.Set(nil)

		mockProjectRepo := mockDB.ProjectsRepo.(*db.MockProjectRepository)
		mockProjectRepo.On("GetProjectByDisplayID", ctx, "test-org").Return((*models.Project)(nil), errors.New("database error"))

		org, err := orgService.internalGetProjectByDisplayID(ctx, "test-org")

		assert.Nil(t, org)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockProjectRepo.AssertExpectations(t)
	})
}
