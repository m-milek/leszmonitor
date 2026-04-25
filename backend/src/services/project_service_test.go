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
	svc := newProjectService(base)

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

	t.Run("Fails when project not found", func(t *testing.T) {
		ctx, svc, mockDB := setupTestProjectService()
		defer db.Set(nil)

		mockDB.ProjectsRepo.(*db.MockProjectRepository).On("GetProjectBySlug", ctx, "nonexistent").Return((*models.Project)(nil), db.ErrNotFound)

		project, err := svc.internalGetProjectBySlug(ctx, "nonexistent")

		assert.Nil(t, project)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockDB.ProjectsRepo.(*db.MockProjectRepository).AssertExpectations(t)
	})

	t.Run("Fails when database returns error", func(t *testing.T) {
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
