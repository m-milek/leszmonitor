package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
	"net/http"
)

type GroupServiceT struct {
	BaseService
}

// NewGroupService creates a new instance of GroupServiceT.
func NewGroupService() *GroupServiceT {
	return &GroupServiceT{
		BaseService{
			authService:   newAuthorizationService(),
			serviceLogger: logging.NewServiceLogger("GroupService"),
		},
	}
}

var GroupService = NewGroupService()

type UpdateMonitorGroupPayload struct {
	Name        string `json:"name"`        // Name of the monitor group
	Description string `json:"description"` // Description of the monitor group
}

type CreateMonitorGroupPayload struct {
	Name        string `json:"name"`        // Name of the monitor group
	Description string `json:"description"` // Description of the monitor group
}

// CreateMonitorGroup creates a new monitor group for the team in the provided TeamAuth.
func (s *GroupServiceT) CreateMonitorGroup(context context.Context, teamAuth *middleware.TeamAuth, payload CreateMonitorGroupPayload) (*models.MonitorGroup, *ServiceError) {
	logger := s.getMethodLogger("CreateMonitorGroup")

	team, authErr := s.authService.authorizeTeamAction(context, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return nil, authErr
	}

	monitorGroup, err := models.NewMonitorGroup(payload.Name, payload.Description, team)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create new monitor group")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor group data: %w", err),
		}
	}

	_, err = db.CreateMonitorGroup(context, monitorGroup)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			logger.Warn().Str("monitorGroupId", monitorGroup.DisplayID).Msg("Monitor group already exists")
			return nil, &ServiceError{
				Code: http.StatusConflict,
				Err:  fmt.Errorf("monitor group with DisplayID %s already exists", monitorGroup.DisplayID),
			}
		}
		logger.Error().Err(err).Msg("Failed to create monitor group")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to create monitor group: %w", err),
		}
	}

	logger.Info().Str("monitorGroupId", monitorGroup.DisplayID).Msg("Monitor group created successfully")
	return monitorGroup, nil
}

// GetTeamMonitorGroups retrieves all monitor groups for the team in the provided TeamAuth.
func (s *GroupServiceT) GetTeamMonitorGroups(context context.Context, teamAuth *middleware.TeamAuth) ([]models.MonitorGroup, *ServiceError) {
	logger := GroupService.getMethodLogger("GetTeamMonitorGroups")

	team, authErr := s.authService.authorizeTeamAction(context, teamAuth, models.PermissionTeamReader)
	if authErr != nil {
		return nil, authErr
	}

	groups, err := db.GetMonitorGroupsForTeam(context, team)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get monitor groups for team")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to get monitor groups for team %s: %w", team.DisplayID, err),
		}
	}

	logger.Info().Int("count", len(groups)).Msg("Retrieved monitor groups for team")
	return groups, nil
}

// GetTeamMonitorGroupById retrieves a specific monitor group by its DisplayID for the team in the provided TeamAuth.
func (s *GroupServiceT) GetTeamMonitorGroupById(context context.Context, teamAuth *middleware.TeamAuth, groupId string) (*models.MonitorGroup, *ServiceError) {
	logger := s.getMethodLogger("GetTeamMonitorGroupById")

	if groupId == "" {
		logger.Warn().Msg("Monitor group DisplayID is required to get monitor group")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("monitor group DisplayID is required"),
		}
	}

	team, authErr := s.authService.authorizeTeamAction(context, teamAuth, models.PermissionTeamReader)
	if authErr != nil {
		return nil, authErr
	}

	group, err := s.internalGetMonitorGroupById(context, team.DisplayID, groupId)
	if err != nil {
		return nil, err
	}

	logger.Info().Str("groupId", group.DisplayID).Msg("Retrieved monitor group by DisplayID")
	return group, nil
}

// DeleteMonitorGroup deletes a specific monitor group by its DisplayID for the team in the provided TeamAuth.
func (s *GroupServiceT) DeleteMonitorGroup(context context.Context, teamAuth *middleware.TeamAuth, groupId string) *ServiceError {
	logger := s.getMethodLogger("DeleteMonitorGroup")

	team, authErr := s.authService.authorizeTeamAction(context, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return authErr
	}

	if groupId == "" {
		logger.Warn().Msg("Monitor group DisplayID is required for deletion")
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("monitor group DisplayID is required"),
		}
	}

	deleted, err := db.DeleteMonitorGroup(context, team, groupId)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete monitor group")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to delete monitor group with DisplayID %s: %w", groupId, err),
		}
	}

	if !deleted {
		logger.Warn().Str("groupId", groupId).Msg("Monitor group not found for deletion")
		return &ServiceError{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("monitor group with DisplayID %s not found", groupId),
		}
	}

	logger.Info().Str("groupId", groupId).Msg("Monitor group deleted successfully")
	return nil
}

// UpdateMonitorGroup updates the details of a specific monitor group by its DisplayID for the team in the provided TeamAuth.
func (s *GroupServiceT) UpdateMonitorGroup(ctx context.Context, teamAuth *middleware.TeamAuth, groupId string, payload *UpdateMonitorGroupPayload) (*models.MonitorGroup, *ServiceError) {
	logger := s.getMethodLogger("UpdateMonitorGroup")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return nil, authErr
	}

	if groupId == "" {
		logger.Warn().Msg("Monitor group DisplayID is required for update")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("monitor group DisplayID is required"),
		}
	}

	oldGroup, err := s.internalGetMonitorGroupById(ctx, team.DisplayID, groupId)
	if err != nil {
		return nil, err
	}

	newGroup := *oldGroup
	newGroup.Name = payload.Name
	newGroup.Description = payload.Description
	newGroup.DisplayIDFromName.Init(newGroup.Name)

	_, updateErr := db.UpdateMonitorGroup(ctx, team, oldGroup, &newGroup)
	if updateErr != nil {
		logger.Error().Err(updateErr).Msg("Failed to update monitor group")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to update monitor group with DisplayID %s: %w", groupId, updateErr),
		}
	}

	logger.Info().Str("groupId", oldGroup.DisplayID).Msg("Monitor group updated successfully")
	return &newGroup, nil
}

func (s *GroupServiceT) internalGetMonitorGroupById(ctx context.Context, teamId string, groupId string) (*models.MonitorGroup, *ServiceError) {
	logger := s.getMethodLogger("internalGetMonitorGroupById")

	group, err := db.GetMonitorGroupById(ctx, groupId)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("groupId", groupId).Msg("Monitor group not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("monitor group with DisplayID %s not found", groupId),
			}
		}
		logger.Error().Err(err).Msg("Failed to get monitor group")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to get monitor group: %w", err),
		}
	}

	return group, nil
}
