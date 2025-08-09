package services

import (
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/common"
	"github.com/m-milek/leszmonitor/db"
	"net/http"
)

type TeamServiceT struct{}

// NewUserService creates a new instance of UserServiceT.
func newTeamService() *TeamServiceT {
	return &TeamServiceT{}
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
	Role     common.TeamRole `json:"role"`     // The role to assign to the user in the team
}

type TeamRemoveMemberPayload struct {
	Username string `json:"username"` // The username of the user to remove from the team
}

type TeamChangeMemberRolePayload struct {
	Username string          `json:"username"` // The username of the user whose role is to be changed
	Role     common.TeamRole `json:"role"`     // The new role to assign to the user in the team
}

func (s *TeamServiceT) GetAllTeams() ([]common.Team, *ServiceError) {
	teams, err := db.GetAllTeams()

	if err != nil {
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving teams: %w", err),
		}
	}

	return teams, nil
}

func (s *TeamServiceT) GetTeamById(teamId string) (*common.Team, *ServiceError) {
	team, err := db.GetTeamById(teamId)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{
				Code: 404,
				Err:  fmt.Errorf("team with ID %s not found", teamId),
			}
		}
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving team: %w", err),
		}
	}

	return team, nil
}

func (s *TeamServiceT) CreateTeam(payload *TeamCreatePayload, ownerUsername string) (*TeamCreateResponse, *ServiceError) {
	team := common.NewTeam(payload.Name, payload.Description, ownerUsername)

	_, err := db.CreateTeam(team)

	if err != nil {
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error creating team: %w", err),
		}
	}

	return &TeamCreateResponse{
		TeamId: team.Id,
	}, nil
}

func (s *TeamServiceT) DeleteTeam(teamId string, requestorUsername string) *ServiceError {
	team, err := db.GetTeamById(teamId)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{
				Code: 404,
				Err:  fmt.Errorf("team with ID %s not found", teamId),
			}
		}
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving team: %w", err),
		}
	}

	if !team.IsAdmin(requestorUsername) {
		return &ServiceError{
			Code: 403,
			Err:  fmt.Errorf("user %s is not authorized to delete team %s", requestorUsername, teamId),
		}
	}

	_, err = db.DeleteTeam(teamId)
	if err != nil {
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error deleting team %s: %w", teamId, err),
		}
	}

	return nil
}

func (s *TeamServiceT) UpdateTeam(teamId string, payload *TeamUpdatePayload, requestorUsername string) (*common.Team, *ServiceError) {
	team, err := db.GetTeamById(teamId)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{
				Code: 404,
				Err:  fmt.Errorf("team with ID %s not found", teamId),
			}
		}
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving team: %w", err),
		}
	}

	if !team.IsAdmin(requestorUsername) {
		return nil, &ServiceError{
			Code: 403,
			Err:  fmt.Errorf("user %s is not authorized to update team %s", requestorUsername, teamId),
		}
	}

	team.Name = payload.Name
	team.Description = payload.Description
	_, err = db.UpdateTeam(team)

	if err != nil {
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error updating team %s: %w", teamId, err),
		}
	}

	return team, nil
}

func (s *TeamServiceT) AddUserToTeam(teamId string, payload *TeamAddMemberPayload, requestorUsername string) *ServiceError {
	if teamId == "" || payload.Username == "" || payload.Role == "" {
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("teamId, username, and role are required"),
		}
	}

	team, err := db.GetTeamById(teamId)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{
				Code: 404,
				Err:  fmt.Errorf("team with ID %s not found", teamId),
			}
		}
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving team: %w", err),
		}
	}

	if !team.IsAdmin(requestorUsername) {
		return &ServiceError{
			Code: 403,
			Err:  fmt.Errorf("user %s is not authorized to add users to team %s", requestorUsername, teamId),
		}
	}

	user, err := db.GetUserByUsername(payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{
				Code: 404,
				Err:  fmt.Errorf("user %s not found", payload.Username),
			}
		}
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving user %s: %w", payload.Username, err),
		}
	}

	if user == nil {
		return &ServiceError{
			Code: 404,
			Err:  fmt.Errorf("user %s not found", payload.Username),
		}
	}

	if err := payload.Role.Validate(); err != nil {
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid role: %w", err),
		}
	}

	err = team.AddMember(user.Username, payload.Role)
	if err != nil {
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error adding user %s to team %s: %w", payload.Username, teamId, err),
		}
	}

	_, err = db.UpdateTeam(team)
	if err != nil {
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error adding user %s to team %s: %w", payload.Username, teamId, err),
		}
	}

	return nil
}

func (s *TeamServiceT) RemoveUserFromTeam(teamId string, payload *TeamRemoveMemberPayload, requestorUsername string) *ServiceError {
	team, err := db.GetTeamById(teamId)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{
				Code: 404,
				Err:  fmt.Errorf("team with ID %s not found", teamId),
			}
		}
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving team: %w", err),
		}
	}

	if !team.IsAdmin(requestorUsername) {
		return &ServiceError{
			Code: 403,
			Err:  fmt.Errorf("user %s is not authorized to remove users from team %s", requestorUsername, teamId),
		}
	}

	if !team.IsMember(payload.Username) {
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of team %s", payload.Username, teamId),
		}
	}

	team.RemoveMember(payload.Username)

	_, err = db.UpdateTeam(team)
	if err != nil {
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error removing user %s from team %s: %w", payload.Username, teamId, err),
		}
	}

	return nil
}

func (s *TeamServiceT) ChangeMemberRole(teamId string, payload TeamChangeMemberRolePayload, requestorUsername string) *ServiceError {
	team, err := db.GetTeamById(teamId)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{
				Code: 404,
				Err:  fmt.Errorf("team with ID %s not found", teamId),
			}
		}
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving team: %w", err),
		}
	}

	if !team.IsAdmin(requestorUsername) {
		return &ServiceError{
			Code: 403,
			Err:  fmt.Errorf("user %s is not authorized to change member roles in team %s", requestorUsername, teamId),
		}
	}

	if !team.IsMember(payload.Username) {
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of team %s", payload.Username, teamId),
		}
	}

	if err := payload.Role.Validate(); err != nil {
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid role: %w", err),
		}
	}

	err = team.ChangeMemberRole(payload.Username, payload.Role)
	if err != nil {
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error changing role for user %s in team %s: %w", payload.Username, teamId, err),
		}
	}

	_, err = db.UpdateTeam(team)
	if err != nil {
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error updating team after changing role for user %s in team %s: %w", payload.Username, teamId, err),
		}
	}

	return nil
}
