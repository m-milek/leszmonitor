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

func setupTestGroupService() (context.Context, *GroupServiceT, *db.MockDB) {
	mockDB := &db.MockDB{
		GroupsRepo: new(db.MockGroupRepository),
	}
	db.Set(mockDB)

	base := newBaseService(nil, "GroupServiceTest")
	groupService := newGroupService(base)

	ctx := context.Background()

	return ctx, groupService, mockDB
}

func TestGroupServiceT_InternalGetMonitorGroupByID(t *testing.T) {
	t.Run("Returns monitor group successfully", func(t *testing.T) {
		ctx, groupService, mockDB := setupTestGroupService()
		defer db.Set(nil)

		expectedGroup := &models.MonitorGroup{
			Description: "Test Description",
		}
		expectedGroup.DisplayIDFromName.Init("test-group")

		mockGroupRepo := mockDB.GroupsRepo.(*db.MockGroupRepository)
		mockGroupRepo.On("GetGroupByDisplayID", ctx, "test-group").Return(expectedGroup, nil)

		group, err := groupService.internalGetMonitorGroupByID(ctx, "test-group")

		assert.Nil(t, err)
		assert.NotNil(t, group)
		assert.Equal(t, "test-group", group.DisplayID)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("Fails when monitor group not found", func(t *testing.T) {
		ctx, groupService, mockDB := setupTestGroupService()
		defer db.Set(nil)

		mockGroupRepo := mockDB.GroupsRepo.(*db.MockGroupRepository)
		mockGroupRepo.On("GetGroupByDisplayID", ctx, "nonexistent").Return((*models.MonitorGroup)(nil), db.ErrNotFound)

		group, err := groupService.internalGetMonitorGroupByID(ctx, "nonexistent")

		assert.Nil(t, group)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("Fails when database returns error", func(t *testing.T) {
		ctx, groupService, mockDB := setupTestGroupService()
		defer db.Set(nil)

		mockGroupRepo := mockDB.GroupsRepo.(*db.MockGroupRepository)
		mockGroupRepo.On("GetGroupByDisplayID", ctx, "test-group").Return((*models.MonitorGroup)(nil), errors.New("database error"))

		group, err := groupService.internalGetMonitorGroupByID(ctx, "test-group")

		assert.Nil(t, group)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockGroupRepo.AssertExpectations(t)
	})
}
