package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
)

// TeamServiceT handles team-related CRUD operations.
type TeamServiceT struct {
	baseService
}

// newTeamService creates a new instance of TeamServiceT.
func newTeamService() *TeamServiceT {
	return &TeamServiceT{
		baseService{
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
	TeamID string `json:"teamId"` // The DisplayID of the newly created team
}

type TeamUpdatePayload struct {
	Name        string `json:"name"`        // The new name of the team
	Description string `json:"description"` // A new description for the team
}

type TeamAddMemberPayload struct {
	Username string      `json:"username"` // The username of the user to add to the team
	Role     models.Role `json:"role"`     // The role to assign to the user in the team
}

type TeamRemoveMemberPayload struct {
	Username string `json:"username"` // The username of the user to remove from the team
}

type TeamChangeMemberRolePayload struct {
	Username string      `json:"username"` // The username of the user whose role is to be changed
	Role     models.Role `json:"role"`     // The new role to assign to the user in the team
}

// GetAllTeams retrieves all teams from the database. No authentication is required at the moment.
func (s *TeamServiceT) GetAllTeams(ctx context.Context) ([]models.Team, *ServiceError) {
	logger := s.getMethodLogger("GetAllTeams")
	logger.Trace().Msg("Retrieving all teams")

	teams, err := s.getDB().Teams().GetAllTeams(ctx)

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

// GetTeamByID retrieves a team by its DisplayID, ensuring the requesting user has at least reader permissions.
func (s *TeamServiceT) GetTeamByID(ctx context.Context, teamAuth *middleware.TeamAuth) (*models.Team, *ServiceError) {
	logger := s.getMethodLogger("GetTeamByDisplayID")
	logger.Trace().Str("teamId", teamAuth.TeamID).Msg("Retrieving team by DisplayID")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamReader)
	if authErr != nil {
		return nil, authErr
	}

	logger.Trace().Str("teamId", team.DisplayID).Msg("Retrieved team successfully")
	return team, nil
}

// CreateTeam creates a new team with the given payload and assigns the owner by username.
func (s *TeamServiceT) CreateTeam(ctx context.Context, ownerUsername string, payload *TeamCreatePayload) (*TeamCreateResponse, *ServiceError) {
	logger := s.getMethodLogger("CreateTeam")
	logger.Trace().Any("payload", payload).Str("username", ownerUsername).Msg("Creating new team")

	user, err := db.Get().Users().GetUserByUsername(ctx, ownerUsername)
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

	team, err := models.NewTeam(payload.Name, payload.Description, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create new team model")
		return nil, &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid team data: %w", err),
		}
	}

	_, serviceErr := s.internalCreateTeam(ctx, team)
	if serviceErr != nil {
		return nil, serviceErr
	}

	logger.Trace().Str("teamId", team.DisplayID).Msg("Team created successfully")
	return &TeamCreateResponse{
		TeamID: team.DisplayID,
	}, nil
}

// DeleteTeam deletes a team by its DisplayID.
// Requires admin permissions.
func (s *TeamServiceT) DeleteTeam(ctx context.Context, teamAuth *middleware.TeamAuth) *ServiceError {
	logger := s.getMethodLogger("DeleteTeam")
	logger.Trace().Str("teamId", teamAuth.TeamID).Str("requestorUsername", teamAuth.Username).Msg("Deleting team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamAdmin)
	if authErr != nil {
		return authErr
	}

	_, err := db.Get().Teams().DeleteTeamByID(ctx, team.DisplayID)
	if err != nil {
		logger.Error().Err(err).Str("teamId", team.DisplayID).Msg("Failed to delete team")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to delete team %s: %w", team.DisplayID, err),
		}
	}

	logger.Trace().Str("teamId", team.DisplayID).Msg("Team deleted successfully")
	return nil
}

// UpdateTeam updates the details of a team.
// Requires editor permissions or higher.
func (s *TeamServiceT) UpdateTeam(ctx context.Context, teamAuth *middleware.TeamAuth, payload *TeamUpdatePayload) (*models.Team, *ServiceError) {
	logger := s.getMethodLogger("UpdateTeam")
	logger.Trace().Str("teamId", teamAuth.TeamID).Any("payload", payload).Str("requestorUsername", teamAuth.Username).Msg("Updating team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return nil, authErr
	}

	team.Name = payload.Name
	team.Description = payload.Description
	team.DisplayIDFromName.Init(team.Name)

	_, err := s.getDB().Teams().UpdateTeam(ctx, team)

	if err != nil {
		logger.Error().Err(err).Str("teamId", team.DisplayID).Msg("Failed to update team")
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to update team %s: %w", team.DisplayID, err),
		}
	}

	logger.Trace().Str("teamId", team.DisplayID).Msg("Team updated successfully")
	return team, nil
}

// AddUserToTeam adds a user to a team with a specified role.
// Requires editor permissions or higher.
func (s *TeamServiceT) AddUserToTeam(ctx context.Context, teamAuth *middleware.TeamAuth, payload *TeamAddMemberPayload) *ServiceError {
	logger := s.getMethodLogger("AddUserToTeam")
	logger.Trace().Str("teamId", teamAuth.TeamID).Any("payload", payload).Str("requestorUsername", teamAuth.Username).Msg("Adding user to team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return authErr
	}

	user, err := db.Get().Users().GetUserByUsername(ctx, payload.Username)
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

	teamMember, err := models.NewTeamMember(user.ID, payload.Role)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Any("role", payload.Role).Msg("Failed to create team member model")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to create team member for user %s: %w", payload.Username, err),
		}
	}

	_, err = s.getDB().Teams().AddMemberToTeam(ctx, team.DisplayID, teamMember)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			logger.Warn().Str("teamId", team.DisplayID).Str("username", payload.Username).Msg("User already a member of team")
			return &ServiceError{
				Code: 409,
				Err:  fmt.Errorf("user %s is already a member of team %s", payload.Username, team.DisplayID),
			}
		}
		logger.Error().Err(err).Str("teamId", team.DisplayID).Str("username", payload.Username).Msg("Failed to add user to team")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to add user %s to team %s: %w", payload.Username, team.DisplayID, err),
		}
	}

	return nil
}

// RemoveUserFromTeam removes a user from a team.
// Requires editor permissions or higher.
func (s *TeamServiceT) RemoveUserFromTeam(ctx context.Context, teamAuth *middleware.TeamAuth, payload *TeamRemoveMemberPayload) *ServiceError {
	logger := s.getMethodLogger("RemoveUserFromTeam")
	logger.Trace().Str("teamId", teamAuth.TeamID).Any("payload", payload).Str("requestorUsername", teamAuth.Username).Msg("Removing user from team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamEditor)
	if authErr != nil {
		return authErr
	}

	user, err := UserService.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		return err
	}

	if !team.IsMember(user.ID) {
		logger.Warn().Str("teamId", team.DisplayID).Str("username", payload.Username).Msg("User not a member of team")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of team %s", payload.Username, team.DisplayID),
		}
	}

	removed, dbErr := s.getDB().Teams().RemoveMemberFromTeam(ctx, team.DisplayID, user.ID)
	if dbErr != nil {
		logger.Error().Err(err).Str("teamId", team.DisplayID).Str("username", payload.Username).Msg("Failed to remove user from team")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to remove user %s from team %s: %w", payload.Username, team.DisplayID, err),
		}
	}

	if !removed {
		logger.Warn().Str("teamId", team.DisplayID).Str("username", payload.Username).Msg("User not a member of team")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of team %s", payload.Username, team.DisplayID),
		}
	}

	return nil
}

// ChangeMemberRole changes the role of a team member.
// Requires admin permissions.
func (s *TeamServiceT) ChangeMemberRole(ctx context.Context, teamAuth *middleware.TeamAuth, payload TeamChangeMemberRolePayload) *ServiceError {
	logger := s.getMethodLogger("ChangeMemberRole")
	logger.Trace().Str("teamId", teamAuth.TeamID).Any("payload", payload).Str("requestorUsername", teamAuth.Username).Msg("Changing member role in team")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionTeamAdmin)
	if authErr != nil {
		return authErr
	}

	user, err := UserService.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		return err
	}

	if !team.IsMember(user.ID) {
		logger.Warn().Str("teamId", team.DisplayID).Str("username", payload.Username).Msg("User not a member of team")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of team %s", payload.Username, team.DisplayID),
		}
	}

	if err := payload.Role.Validate(); err != nil {
		logger.Warn().Str("teamId", team.DisplayID).Str("username", payload.Username).Any("role", payload.Role).Msg("Invalid role for user")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid role: %w", err),
		}
	}

	changeRoleErr := team.ChangeMemberRole(user.ID, payload.Role)
	if changeRoleErr != nil {
		logger.Error().Err(changeRoleErr).Str("teamId", team.DisplayID).Str("username", payload.Username).Msg("Error changing role for user in team")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error changing role for user %s in team %s: %w", payload.Username, team.DisplayID, changeRoleErr),
		}
	}

	return nil
}

func (s *TeamServiceT) internalGetTeamByID(ctx context.Context, id string) (*models.Team, *ServiceError) {
	team, err := s.getDB().Teams().GetTeamByDisplayID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("team with DisplayID %s not found", id),
			}
		}
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve team: %w", err),
		}
	}
	return team, nil
}

func (s *TeamServiceT) internalCreateTeam(ctx context.Context, team *models.Team) (*models.Team, *ServiceError) {
	_, err := s.getDB().Teams().InsertTeam(ctx, team)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return nil, &ServiceError{
				Code: http.StatusConflict,
				Err:  fmt.Errorf("team with DisplayID '%s' already exists", team.DisplayID),
			}
		}
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to create team: %w", err),
		}
	}
	return team, nil
}
