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
	baseService
}

// newGroupService creates a new instance of GroupServiceT.
func newGroupService(service baseService) *GroupServiceT {
	service.serviceLogger = logging.NewServiceLogger("GroupService")
	return &GroupServiceT{
		baseService: service,
	}
}

var GroupService = newGroupService(NewBaseService(db.Get(), newAuthorizationService(), "GroupService"))

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
	logger := s.getMethodLogger("InsertGroup")

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

	err = db.Get().Groups().InsertGroup(context, monitorGroup)
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

	createdGroup, err := db.Get().Groups().GetGroupByDisplayID(context, monitorGroup.DisplayID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch created monitor group")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to fetch created monitor group: %w", err),
		}
	}

	logger.Info().Str("monitorGroupId", monitorGroup.DisplayID).Msg("Monitor group created successfully")
	return createdGroup, nil
}

// GetTeamMonitorGroups retrieves all monitor groups for the team in the provided TeamAuth.
func (s *GroupServiceT) GetTeamMonitorGroups(context context.Context, teamAuth *middleware.TeamAuth) ([]models.MonitorGroup, *ServiceError) {
	logger := GroupService.getMethodLogger("GetTeamMonitorGroups")

	team, authErr := s.authService.authorizeTeamAction(context, teamAuth, models.PermissionTeamReader)
	if authErr != nil {
		return nil, authErr
	}

	groups, err := db.Get().Groups().GetGroupsByTeamID(context, team)
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

// GetTeamMonitorGroupByID retrieves a specific monitor group by its DisplayID for the team in the provided TeamAuth.
func (s *GroupServiceT) GetTeamMonitorGroupByID(context context.Context, teamAuth *middleware.TeamAuth, groupID string) (*models.MonitorGroup, *ServiceError) {
	logger := s.getMethodLogger("GetTeamMonitorGroupByID")

	if groupID == "" {
		logger.Warn().Msg("Monitor group DisplayID is required to get monitor group")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("monitor group DisplayID is required"),
		}
	}

	_, authErr := s.authService.authorizeTeamAction(context, teamAuth, models.PermissionTeamReader)
	if authErr != nil {
		return nil, authErr
	}

	group, err := s.internalGetMonitorGroupByID(context, groupID)
	if err != nil {
		return nil, err
	}

	logger.Info().Str("groupID", group.DisplayID).Msg("Retrieved monitor group by DisplayID")
	return group, nil
}

// DeleteMonitorGroup deletes a specific monitor group by its DisplayID for the team in the provided TeamAuth.
func (s *GroupServiceT) DeleteMonitorGroup(context context.Context, teamAuth *middleware.TeamAuth, groupID string) *ServiceError {
	logger := s.getMethodLogger("DeleteGroup")

	team, authErr := s.authService.authorizeTeamAction(context, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return authErr
	}

	if groupID == "" {
		logger.Warn().Msg("Monitor group DisplayID is required for deletion")
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("monitor group DisplayID is required"),
		}
	}

	deleted, err := db.Get().Groups().DeleteGroup(context, team, groupID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete monitor group")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to delete monitor group with DisplayID %s: %w", groupID, err),
		}
	}

	if !deleted {
		logger.Warn().Str("groupID", groupID).Msg("Monitor group not found for deletion")
		return &ServiceError{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("monitor group with DisplayID %s not found", groupID),
		}
	}

	logger.Info().Str("groupID", groupID).Msg("Monitor group deleted successfully")
	return nil
}

// UpdateMonitorGroup updates the details of a specific monitor group by its DisplayID for the team in the provided TeamAuth.
func (s *GroupServiceT) UpdateMonitorGroup(ctx context.Context, teamAuth *middleware.TeamAuth, groupID string, payload *UpdateMonitorGroupPayload) (*models.MonitorGroup, *ServiceError) {
	logger := s.getMethodLogger("UpdateGroup")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return nil, authErr
	}

	if groupID == "" {
		logger.Warn().Msg("Monitor group DisplayID is required for update")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("monitor group DisplayID is required"),
		}
	}

	oldGroup, err := s.internalGetMonitorGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	newGroup := *oldGroup
	newGroup.Name = payload.Name
	newGroup.Description = payload.Description
	newGroup.DisplayIDFromName.Init(newGroup.Name)

	_, updateErr := db.Get().Groups().UpdateGroup(ctx, team, oldGroup, &newGroup)
	if updateErr != nil {
		logger.Error().Err(updateErr).Msg("Failed to update monitor group")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to update monitor group with DisplayID %s: %w", groupID, updateErr),
		}
	}

	logger.Info().Str("groupID", oldGroup.DisplayID).Msg("Monitor group updated successfully")
	return &newGroup, nil
}

func (s *GroupServiceT) internalGetMonitorGroupByID(ctx context.Context, groupID string) (*models.MonitorGroup, *ServiceError) {
	logger := s.getMethodLogger("internalGetMonitorGroupByID")

	group, err := db.Get().Groups().GetGroupByDisplayID(ctx, groupID)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("groupID", groupID).Msg("Monitor group not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("monitor group with DisplayID %s not found", groupID),
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
