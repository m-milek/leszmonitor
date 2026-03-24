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

// OrgServiceT handles org-related CRUD operations.
type OrgServiceT struct {
	baseService
}

// newOrgService creates a new instance of OrgServiceT.
func newOrgService() *OrgServiceT {
	return &OrgServiceT{
		baseService{
			authService:   newAuthorizationService(),
			serviceLogger: logging.NewServiceLogger("OrgService"),
		},
	}
}

var OrgService = newOrgService()

type CreateOrgPayload struct {
	Name        string `json:"name"`        // The name of the org
	Description string `json:"description"` // A brief description of the org
}

type CreateOrgResponse struct {
	OrgID string `json:"orgID"` // The DisplayID of the newly created org
}

type UpdateOrgPayload struct {
	Name        string `json:"name"`        // The new name of the org
	Description string `json:"description"` // A new description for the org
}

type AddOrgMemberPayload struct {
	Username string      `json:"username"` // The username of the user to add to the org
	Role     models.Role `json:"role"`     // The role to assign to the user in the org
}

type RemoveOrgMemberPayload struct {
	Username string `json:"username"` // The username of the user to remove from the org
}

type ChangeOrgMemberRolePayload struct {
	Username string      `json:"username"` // The username of the user whose role is to be changed
	Role     models.Role `json:"role"`     // The new role to assign to the user in the org
}

// GetAllOrgs retrieves all orgs from the database. No authentication is required at the moment.
func (s *OrgServiceT) GetAllOrgs(ctx context.Context) ([]models.Org, *ServiceError) {
	logger := s.getMethodLogger("GetAllOrgs")
	logger.Trace().Msg("Retrieving all orgs")

	orgs, err := s.getDB().Orgs().GetAllOrgs(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve orgs")
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error retrieving orgs: %w", err),
		}
	}

	logger.Trace().Int("count", len(orgs)).Msg("Retrieved orgs successfully")
	return orgs, nil
}

// GetOrgByID retrieves a org by its DisplayID, ensuring the requesting user has at least reader permissions.
func (s *OrgServiceT) GetOrgByID(ctx context.Context, orgAuth *middleware.OrgAuth) (*models.Org, *ServiceError) {
	logger := s.getMethodLogger("GetOrgByDisplayID")
	logger.Trace().Str("orgID", orgAuth.OrgID).Msg("Retrieving org by DisplayID")

	org, authErr := s.authService.authorizeOrgAction(ctx, orgAuth, models.PermissionOrgReader)
	if authErr != nil {
		return nil, authErr
	}

	logger.Trace().Str("orgID", org.DisplayID).Msg("Retrieved org successfully")
	return org, nil
}

// CreateOrg creates a new org with the given payload and assigns the owner by username.
func (s *OrgServiceT) CreateOrg(ctx context.Context, ownerUsername string, payload *CreateOrgPayload) (*CreateOrgResponse, *ServiceError) {
	logger := s.getMethodLogger("CreateOrg")
	logger.Trace().Any("payload", payload).Str("username", ownerUsername).Msg("Creating new org")

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

	org, err := models.NewOrg(payload.Name, payload.Description, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create new org model")
		return nil, &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid org data: %w", err),
		}
	}

	_, serviceErr := s.internalCreateOrg(ctx, org)
	if serviceErr != nil {
		return nil, serviceErr
	}

	logger.Trace().Str("orgID", org.DisplayID).Msg("Org created successfully")
	return &CreateOrgResponse{
		OrgID: org.DisplayID,
	}, nil
}

// DeleteOrg deletes a org by its DisplayID.
// Requires admin permissions.
func (s *OrgServiceT) DeleteOrg(ctx context.Context, orgAuth *middleware.OrgAuth) *ServiceError {
	logger := s.getMethodLogger("DeleteOrg")
	logger.Trace().Str("orgID", orgAuth.OrgID).Str("requestorUsername", orgAuth.Username).Msg("Deleting org")

	org, authErr := s.authService.authorizeOrgAction(ctx, orgAuth, models.PermissionOrgAdmin)
	if authErr != nil {
		return authErr
	}

	_, err := db.Get().Orgs().DeleteOrgByID(ctx, org.DisplayID)
	if err != nil {
		logger.Error().Err(err).Str("orgID", org.DisplayID).Msg("Failed to delete org")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to delete org %s: %w", org.DisplayID, err),
		}
	}

	logger.Trace().Str("orgID", org.DisplayID).Msg("Org deleted successfully")
	return nil
}

// UpdateOrg updates the details of a org.
// Requires editor permissions or higher.
func (s *OrgServiceT) UpdateOrg(ctx context.Context, orgAuth *middleware.OrgAuth, payload *UpdateOrgPayload) (*models.Org, *ServiceError) {
	logger := s.getMethodLogger("UpdateOrg")
	logger.Trace().Str("orgID", orgAuth.OrgID).Any("payload", payload).Str("requestorUsername", orgAuth.Username).Msg("Updating org")

	org, authErr := s.authService.authorizeOrgAction(ctx, orgAuth, models.PermissionOrgEditor)
	if authErr != nil {
		return nil, authErr
	}

	org.Name = payload.Name
	org.Description = payload.Description
	org.DisplayIDFromName.Init(org.Name)

	_, err := s.getDB().Orgs().UpdateOrg(ctx, org)

	if err != nil {
		logger.Error().Err(err).Str("orgID", org.DisplayID).Msg("Failed to update org")
		return nil, &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to update org %s: %w", org.DisplayID, err),
		}
	}

	logger.Trace().Str("orgID", org.DisplayID).Msg("Org updated successfully")
	return org, nil
}

// AddUserToOrg adds a user to a org with a specified role.
// Requires editor permissions or higher.
func (s *OrgServiceT) AddUserToOrg(ctx context.Context, orgAuth *middleware.OrgAuth, payload *AddOrgMemberPayload) *ServiceError {
	logger := s.getMethodLogger("AddUserToOrg")
	logger.Trace().Str("orgID", orgAuth.OrgID).Any("payload", payload).Str("requestorUsername", orgAuth.Username).Msg("Adding user to org")

	org, authErr := s.authService.authorizeOrgAction(ctx, orgAuth, models.PermissionOrgEditor)
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
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to retrieve user for adding to org")
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

	orgMember, err := models.NewOrgMember(user.ID, payload.Role)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Any("role", payload.Role).Msg("Failed to create org member model")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to create org member for user %s: %w", payload.Username, err),
		}
	}

	_, err = s.getDB().Orgs().AddMemberToOrg(ctx, org.DisplayID, orgMember)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			logger.Warn().Str("orgID", org.DisplayID).Str("username", payload.Username).Msg("User already a member of org")
			return &ServiceError{
				Code: 409,
				Err:  fmt.Errorf("user %s is already a member of org %s", payload.Username, org.DisplayID),
			}
		}
		logger.Error().Err(err).Str("orgID", org.DisplayID).Str("username", payload.Username).Msg("Failed to add user to org")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to add user %s to org %s: %w", payload.Username, org.DisplayID, err),
		}
	}

	return nil
}

// RemoveUserFromOrg removes a user from a org.
// Requires editor permissions or higher.
func (s *OrgServiceT) RemoveUserFromOrg(ctx context.Context, orgAuth *middleware.OrgAuth, payload *RemoveOrgMemberPayload) *ServiceError {
	logger := s.getMethodLogger("RemoveUserFromOrg")
	logger.Trace().Str("orgID", orgAuth.OrgID).Any("payload", payload).Str("requestorUsername", orgAuth.Username).Msg("Removing user from org")

	org, authErr := s.authService.authorizeOrgAction(ctx, orgAuth, models.PermissionOrgEditor)
	if authErr != nil {
		return authErr
	}

	user, err := UserService.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		return err
	}

	orgMember := org.GetMember(user.ID)
	if orgMember == nil {
		logger.Warn().Str("orgID", org.DisplayID).Str("username", payload.Username).Msg("User not a member of org")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of org %s", payload.Username, org.DisplayID),
		}
	}

	if orgMember.Role == models.RoleOwner {
		logger.Warn().Str("orgID", org.DisplayID).Str("username", payload.Username).Msg("Cannot remove org owner")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("cannot remove org owner %s from org %s", payload.Username, org.DisplayID),
		}
	}

	removed, dbErr := s.getDB().Orgs().RemoveMemberFromOrg(ctx, org.DisplayID, user.ID)
	if dbErr != nil {
		logger.Error().Err(err).Str("orgID", org.DisplayID).Str("username", payload.Username).Msg("Failed to remove user from org")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("failed to remove user %s from org %s: %w", payload.Username, org.DisplayID, err),
		}
	}

	if !removed {
		logger.Warn().Str("orgID", org.DisplayID).Str("username", payload.Username).Msg("User not a member of org")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of org %s", payload.Username, org.DisplayID),
		}
	}

	return nil
}

// ChangeMemberRole changes the role of an org member.
// Requires admin permissions.
func (s *OrgServiceT) ChangeMemberRole(ctx context.Context, orgAuth *middleware.OrgAuth, payload ChangeOrgMemberRolePayload) *ServiceError {
	logger := s.getMethodLogger("ChangeMemberRole")
	logger.Trace().Str("orgID", orgAuth.OrgID).Any("payload", payload).Str("requestorUsername", orgAuth.Username).Msg("Changing member role in org")

	org, authErr := s.authService.authorizeOrgAction(ctx, orgAuth, models.PermissionOrgAdmin)
	if authErr != nil {
		return authErr
	}

	user, err := UserService.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		return err
	}

	if !org.IsMember(user.ID) {
		logger.Warn().Str("orgID", org.DisplayID).Str("username", payload.Username).Msg("User not a member of org")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("user %s is not a member of org %s", payload.Username, org.DisplayID),
		}
	}

	if err := payload.Role.Validate(); err != nil {
		logger.Warn().Str("orgID", org.DisplayID).Str("username", payload.Username).Any("role", payload.Role).Msg("Invalid role for user")
		return &ServiceError{
			Code: 400,
			Err:  fmt.Errorf("invalid role: %w", err),
		}
	}

	changeRoleErr := org.ChangeMemberRole(user.ID, payload.Role)
	if changeRoleErr != nil {
		logger.Error().Err(changeRoleErr).Str("orgID", org.DisplayID).Str("username", payload.Username).Msg("Error changing role for user in org")
		return &ServiceError{
			Code: 500,
			Err:  fmt.Errorf("error changing role for user %s in org %s: %w", payload.Username, org.DisplayID, changeRoleErr),
		}
	}

	return nil
}

func (s *OrgServiceT) internalGetOrgByID(ctx context.Context, id string) (*models.Org, *ServiceError) {
	org, err := s.getDB().Orgs().GetOrgByDisplayID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("org with DisplayID %s not found", id),
			}
		}
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve org: %w", err),
		}
	}
	return org, nil
}

func (s *OrgServiceT) internalCreateOrg(ctx context.Context, org *models.Org) (*models.Org, *ServiceError) {
	_, err := s.getDB().Orgs().InsertOrg(ctx, org)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return nil, &ServiceError{
				Code: http.StatusConflict,
				Err:  fmt.Errorf("org with DisplayID '%s' already exists", org.DisplayID),
			}
		}
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to create org: %w", err),
		}
	}
	return org, nil
}
