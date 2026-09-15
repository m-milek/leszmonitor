package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/security"
)

type ITagService interface {
	CreateTag(ctx context.Context, tag models.Tag) (*models.Tag, *ServiceError)
	GetTagByID(ctx context.Context, id string) (*models.Tag, *ServiceError)
	GetAllTags(ctx context.Context) ([]models.Tag, *ServiceError)
	UpdateTag(ctx context.Context, tag models.Tag) (*models.Tag, *ServiceError)
	DeleteTag(ctx context.Context, id string) *ServiceError
}

// TagService handles tag-related CRUD operations.
type TagService struct {
	db db.DB
}

type TagServiceDeps struct {
	DB db.DB
}

func NewTagService(deps TagServiceDeps) *TagService {
	return &TagService{
		db: deps.DB,
	}
}

// CreateTag creates a new tag. The ID is always assigned server-side.
func (s *TagService) CreateTag(ctx context.Context, tag models.Tag) (*models.Tag, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameTag, "CreateTag")
	logger.Trace().Interface("tag", tag).Msg("Creating tag")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, NewUnauthorizedError("user claims not found in context")
	}

	newTag := models.Tag{
		ID:          uuid.New(),
		Name:        tag.Name,
		Description: tag.Description,
		ColorHex:    tag.ColorHex,
	}
	newTag.Normalize()

	if err := newTag.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid tag")
		return nil, NewBadRequestError("invalid tag: %w", err)
	}

	tagFromDB, txErr := db.WithAuditedTx(ctx, s.db, func(tx db.DB) (*models.Tag, *security.AuditLogParams, error) {
		t, createErr := tx.Tags().InsertTag(ctx, newTag)
		if createErr != nil {
			return nil, nil, createErr
		}

		params := &security.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &t.ID,
			Action:     security.ActionCreateTag,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Tag with ID %s created", t.ID),
			After:      t,
		}
		return t, params, nil
	})
	if txErr != nil {
		logger.Error().Err(txErr).Msg("Failed to create tag within transaction")
		return nil, NewInternalError("failed to create tag within transaction: %w", txErr)
	}

	logger.Debug().Str("id", tagFromDB.ID.String()).Msg("Tag created")
	return tagFromDB, nil
}

// GetTagByID retrieves a single tag by its ID.
func (s *TagService) GetTagByID(ctx context.Context, id string) (*models.Tag, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameTag, "GetTagByID")
	logger.Trace().Str("id", id).Msg("Retrieving tag by ID")

	tagUUID, parseErr := uuid.Parse(id)
	if parseErr != nil {
		logger.Error().Str("id", id).Msg("Invalid tag ID format")
		return nil, NewBadRequestError("invalid tag ID format: %w", parseErr)
	}

	tag, err := s.db.Tags().GetTagByID(ctx, tagUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Tag not found in database")
			return nil, NewNotFoundError("tag with ID %s not found", id)
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve tag from database")
		return nil, NewInternalError("failed to retrieve tag: %w", err)
	}

	logger.Debug().Str("id", id).Msg("Tag retrieved successfully")
	return tag, nil
}

// GetAllTags retrieves every tag in the instance.
func (s *TagService) GetAllTags(ctx context.Context) ([]models.Tag, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameTag, "GetAllTags")
	logger.Trace().Msg("Retrieving all tags")

	allTags, err := s.db.Tags().GetAllTags(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve tags from database")
		return nil, NewInternalError("failed to retrieve tags: %w", err)
	}

	logger.Debug().Int("count", len(allTags)).Msg("Tags retrieved successfully")
	return allTags, nil
}

// UpdateTag updates an existing tag.
func (s *TagService) UpdateTag(ctx context.Context, tag models.Tag) (*models.Tag, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameTag, "UpdateTag")
	logger.Trace().Interface("tag", tag).Msg("Updating tag")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, NewUnauthorizedError("user claims not found in context")
	}

	tag.Normalize()
	if err := tag.Validate(); err != nil {
		logger.Error().Err(err).Str("id", tag.ID.String()).Msg("Invalid tag")
		return nil, NewBadRequestError("invalid tag: %w", err)
	}

	updatedTag, txErr := db.WithAuditedTx(ctx, s.db, func(tx db.DB) (*models.Tag, *security.AuditLogParams, error) {
		existingTag, err := tx.Tags().GetTagByID(ctx, tag.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", tag.ID.String()).Msg("Tag not found")
				return nil, nil, NewNotFoundError("tag with ID %s not found", tag.ID)
			}
			logger.Error().Err(err).Str("id", tag.ID.String()).Msg("Failed to retrieve existing tag for update")
			return nil, nil, fmt.Errorf("failed to retrieve existing tag for update: %w", err)
		}

		t, updateErr := tx.Tags().UpdateTag(ctx, tag)
		if updateErr != nil {
			logger.Error().Err(updateErr).Str("id", tag.ID.String()).Msg("Failed to update tag in database")
			return nil, nil, updateErr
		}

		params := &security.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &t.ID,
			Action:     security.ActionUpdateTag,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Tag with ID %s updated", t.ID),
			Before:     existingTag,
			After:      t,
		}
		return t, params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*ServiceError](txErr); ok2 {
			return nil, serviceErr
		}
		logger.Error().Err(txErr).Str("id", tag.ID.String()).Msg("Failed to update tag within transaction")
		return nil, NewInternalError("failed to update tag within transaction: %w", txErr)
	}

	logger.Debug().Str("id", updatedTag.ID.String()).Msg("Tag updated")
	return updatedTag, nil
}

// DeleteTag deletes a tag by its ID.
func (s *TagService) DeleteTag(ctx context.Context, id string) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameTag, "DeleteTag")
	logger.Trace().Str("id", id).Msg("Deleting tag")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return NewUnauthorizedError("user claims not found in context")
	}

	tagUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error().Str("id", id).Msg("Invalid tag ID format")
		return NewBadRequestError("invalid tag ID format: %w", err)
	}

	txErr := db.WithAuditedVoidTx(ctx, s.db, func(tx db.DB) (*security.AuditLogParams, error) {
		tagBeforeDelete, err := tx.Tags().GetTagByID(ctx, tagUUID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", id).Msg("Tag not found in database")
				return nil, NewNotFoundError("tag with ID %s not found", id)
			}
			logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve tag before deletion")
			return nil, fmt.Errorf("failed to retrieve tag before deletion: %w", err)
		}

		if _, delErr := tx.Tags().DeleteTagByID(ctx, tagUUID); delErr != nil {
			logger.Error().Err(delErr).Str("id", id).Msg("Failed to delete tag in database")
			return nil, delErr
		}

		params := &security.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &tagUUID,
			Action:     security.ActionDeleteTag,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Tag with ID %s deleted", tagUUID),
			Before:     tagBeforeDelete,
		}
		return params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*ServiceError](txErr); ok2 {
			return serviceErr
		}
		if errors.Is(txErr, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Tag not found or already deleted")
			return NewNotFoundError("tag not found or already deleted")
		}
		logger.Error().Err(txErr).Str("id", id).Msg("Failed to delete tag")
		return NewInternalError("failed to delete tag: %w", txErr)
	}

	logger.Debug().Str("id", id).Msg("Tag deleted")
	return nil
}
