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

type TeamServiceT struct {
	BaseService
}

// newTeamService creates a new instance of TeamServiceT.
func newTeamService() *TeamServiceT {
	return &TeamServiceT{
		BaseService{
			authService:   newAuthorizationService(),
			serviceLogger: logging.NewServiceLogger("TeamService"),
		},
	}
}

var TeamService = newTeamService()

type TeamCreatePayload struct {
	Name        string `json:"name"`        // The name of the team
	Description string `json:"description"` // A brief description of the team
}

type TeamCreateResponse struct {
	TeamId string `json:"teamId"` // The ID of the newly created team
}

type TeamUpdatePayload struct {
	Name        string `json:"name"`        // The new name of the team
	Description string `json:"description"` // A new description for the team
}

type TeamAddMemberPayload struct {
	Username string          `json:"username"` // The username of the user to add to the team
	Role     models.TeamRole `json:"role"`     // The role to assign to the user in the team
}

type TeamRemoveMemberPayload struct {
	Username string `json:"username"` // The username of the user to remove from the team
}

type TeamChangeMemberRolePayload struct {
	Username string          `json:"username"` // The username of the user whose role is to be changed
	Role     models.TeamRole `json:"role"`     // The new role to assign to the user in the team
}

func (s *TeamServiceT) GetAllTeams(ctx context.Context) ([]models.Team, *ServiceError) {
	logger := s.getMethodLogger("GetAllTeams")
	logger.Trace().Msg("Retrieving all teams")

	teams, err := db.GetAllTeams(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve teams")
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving teams: %w", err),
		}
	}

	logger.Trace().Int("count", len(teams)).Msg("Retrieved teams successfully")
	return teams, nil
}

func (s *TeamServiceT) GetTeamById(ctx context.Context, teamAuth *middleware.TeamAuth) (*models.Team, *ServiceError) {
	logger := s.getMethodLogger("GetTeamById")
	logger.Trace().Str("teamId", teamAuth.TeamId).Msg("Retrieving team by ID")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamReader)
	if authErr != nil {
		return nil, authErr
	}

	logger.Trace().Str("teamId", team.Id).Msg("Retrieved team successfully")
	return team, nil
}

func (s *TeamServiceT) CreateTeam(ctx context.Context, ownerUsername string, payload *TeamCreatePayload) (*TeamCreateResponse, *ServiceError) {
	logger := s.getMethodLogger("CreateTeam")
	logger.Trace().Any("payload", payload).Str("username", ownerUsername).Msg("Creating new team")

	user, err := db.GetUserByUsername(ctx, ownerUsername)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("username", ownerUsername).Msg("Requesting user not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("user %s not found", ownerUsername),
			}
		}
		logger.Error().Err(err).Msg("Failed to retrieve requesting user")
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to retrieve user %s: %w", ownerUsername, err),
		}
	}

	team, err := models.NewTeam(payload.Name, payload.Description, user.Username)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create new team model")
		return nil, &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid team data: %w", err),
		}
	}

	_, err = db.CreateTeam(ctx, team)

	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			logger.Warn().Str("ID", team.Id).Msg("Team already exists")
			return nil, &ServiceError{
				Code: http.StatusConflict,
				Err:  fmt.Errorf("team with ID '%s' already exists", team.Id),
			}
		}
		logger.Error().Err(err).Msg("Failed to create team")
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to create: %w", err),
		}
	}

	logger.Trace().Str("teamId", team.Id).Msg("Team created successfully")
	return &TeamCreateResponse{
		TeamId: team.Id,
	}, nil
}

func (s *TeamServiceT) DeleteTeam(ctx context.Context, teamAuth *middleware.TeamAuth) *ServiceError {
	logger := s.getMethodLogger("DeleteTeam")
	logger.Trace().Str("teamId", teamAuth.TeamId).Str("requestorUsername", teamAuth.Username).Msg("Deleting team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamAdmin)
	if authErr != nil {
		return authErr
	}

	_, err := db.DeleteTeam(ctx, team.Id)
	if err != nil {
		logger.Error().Err(err).Str("teamId", team.Id).Msg("Failed to delete team")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to delete team %s: %w", team.Id, err),
		}
	}

	logger.Trace().Str("teamId", team.Id).Msg("Team deleted successfully")
	return nil
}

func (s *TeamServiceT) UpdateTeam(ctx context.Context, teamAuth *middleware.TeamAuth, payload *TeamUpdatePayload) (*models.Team, *ServiceError) {
	logger := s.getMethodLogger("UpdateTeam")
	logger.Trace().Str("teamId", teamAuth.TeamId).Any("payload", payload).Str("requestorUsername", teamAuth.Username).Msg("Updating team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return nil, authErr
	}

	team.Name = payload.Name
	team.Description = payload.Description
	team.GenerateId()

	_, err := db.UpdateTeam(ctx, team)

	if err != nil {
		logger.Error().Err(err).Str("teamId", team.Id).Msg("Failed to update team")
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to update team %s: %w", team.Id, err),
		}
	}

	logger.Trace().Str("teamId", team.Id).Msg("Team updated successfully")
	return team, nil
}

func (s *TeamServiceT) AddUserToTeam(ctx context.Context, teamAuth *middleware.TeamAuth, payload *TeamAddMemberPayload) *ServiceError {
	logger := s.getMethodLogger("AddUserToTeam")
	logger.Trace().Str("teamId", teamAuth.TeamId).Any("payload", payload).Str("requestorUsername", teamAuth.Username).Msg("Adding user to team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return authErr
	}

	user, err := db.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("username", payload.Username).Msg("User not found")
			return &ServiceError{
				Code: 404,
				Err:  fmt.Errorf("user %s not found", payload.Username),
			}
		}
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to retrieve user for adding to team")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to retrieve user %s: %w", payload.Username, err),
		}
	}

	if err := payload.Role.Validate(); err != nil {
		logger.Warn().Str("username", payload.Username).Any("role", payload.Role).Msg("Invalid role for user")
		return &ServiceError{
			Code: 400,
			Err:  err,
		}
	}

	err = team.AddMember(user.Username, payload.Role)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Str("teamId", team.Id).Msg("Failed to add user to team")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("failed to add user %s to team %s: %w", payload.Username, team.Id, err),
		}
	}

	_, err = db.UpdateTeam(ctx, team)
	if err != nil {
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to update team %s: %w", team.Id, err),
		}
	}

	return nil
}

func (s *TeamServiceT) RemoveUserFromTeam(ctx context.Context, teamAuth *middleware.TeamAuth, payload *TeamRemoveMemberPayload) *ServiceError {
	logger := s.getMethodLogger("RemoveUserFromTeam")
	logger.Trace().Str("teamId", teamAuth.TeamId).Any("payload", payload).Str("requestorUsername", teamAuth.Username).Msg("Removing user from team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return authErr
	}

	if !team.IsMember(payload.Username) {
		logger.Warn().Str("teamId", team.Id).Str("username", payload.Username).Msg("User not a member of team")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of team %s", payload.Username, team.Id),
		}
	}

	team.RemoveMember(payload.Username)

	_, err := db.UpdateTeam(ctx, team)
	if err != nil {
		logger.Error().Err(err).Str("teamId", team.Id).Str("username", payload.Username).Msg("Failed to update team after removing user")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to update team %s after removing user %s: %w", team.Id, payload.Username, err),
		}
	}

	return nil
}

func (s *TeamServiceT) ChangeMemberRole(ctx context.Context, teamAuth *middleware.TeamAuth, payload TeamChangeMemberRolePayload) *ServiceError {
	logger := s.getMethodLogger("ChangeMemberRole")
	logger.Trace().Str("teamId", teamAuth.TeamId).Any("payload", payload).Str("requestorUsername", teamAuth.Username).Msg("Changing member role in team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamAdmin)
	if authErr != nil {
		return authErr
	}

	if !team.IsMember(payload.Username) {
		logger.Warn().Str("teamId", team.Id).Str("username", payload.Username).Msg("User not a member of team")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of team %s", payload.Username, team.Id),
		}
	}

	if err := payload.Role.Validate(); err != nil {
		logger.Warn().Str("teamId", team.Id).Str("username", payload.Username).Any("role", payload.Role).Msg("Invalid role for user")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid role: %w", err),
		}
	}

	err := team.ChangeMemberRole(payload.Username, payload.Role)
	if err != nil {
		logger.Error().Err(err).Str("teamId", team.Id).Str("username", payload.Username).Msg("Error changing role for user in team")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error changing role for user %s in team %s: %w", payload.Username, team.Id, err),
		}
	}

	_, err = db.UpdateTeam(ctx, team)
	if err != nil {
		logger.Error().Err(err).Str("teamId", team.Id).Str("username", payload.Username).Msg("Error updating team after changing role for user")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error updating team after changing role for user %s in team %s: %w", payload.Username, team.Id, err),
		}
	}

	return nil
}
