package tags

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
)

type ITagService interface {
	CreateTag(ctx context.Context, tag Tag) (*Tag, *apperr.ServiceError)
	GetTagByID(ctx context.Context, id string) (*Tag, *apperr.ServiceError)
	GetAllTags(ctx context.Context) ([]Tag, *apperr.ServiceError)
	UpdateTag(ctx context.Context, tag Tag) (*Tag, *apperr.ServiceError)
	DeleteTag(ctx context.Context, id string) *apperr.ServiceError
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
func (s *TagService) CreateTag(ctx context.Context, tag Tag) (*Tag, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameTag, "CreateTag")
	logger.Trace().Interface("tag", tag).Msg("Creating tag")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewUnauthorizedError("user claims not found in context")
	}

	newTag := Tag{
		ID:          uuid.New(),
		Name:        tag.Name,
		Description: tag.Description,
		ColorHex:    tag.ColorHex,
	}
	newTag.Normalize()

	if err := newTag.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid tag")
		return nil, apperr.NewBadRequestError("invalid tag: %w", err)
	}

	tagFromDB, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*Tag, *audit.AuditLogParams, error) {
		t, createErr := NewTagDAO(q).InsertTag(ctx, newTag)
		if createErr != nil {
			return nil, nil, createErr
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &t.ID,
			Action:     audit.ActionCreateTag,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Tag with ID %s created", t.ID),
			After:      t,
		}
		return t, params, nil
	})
	if txErr != nil {
		logger.Error().Err(txErr).Msg("Failed to create tag within transaction")
		return nil, apperr.NewInternalError("failed to create tag within transaction: %w", txErr)
	}

	logger.Debug().Str("id", tagFromDB.ID.String()).Msg("Tag created")
	return tagFromDB, nil
}

// GetTagByID retrieves a single tag by its ID.
func (s *TagService) GetTagByID(ctx context.Context, id string) (*Tag, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameTag, "GetTagByID")
	logger.Trace().Str("id", id).Msg("Retrieving tag by ID")

	tagUUID, parseErr := uuid.Parse(id)
	if parseErr != nil {
		logger.Error().Str("id", id).Msg("Invalid tag ID format")
		return nil, apperr.NewBadRequestError("invalid tag ID format: %w", parseErr)
	}

	tag, err := NewTagDAO(s.db.Querier()).GetTagByID(ctx, tagUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Tag not found in database")
			return nil, apperr.NewNotFoundError("tag with ID %s not found", id)
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve tag from database")
		return nil, apperr.NewInternalError("failed to retrieve tag: %w", err)
	}

	logger.Debug().Str("id", id).Msg("Tag retrieved successfully")
	return tag, nil
}

// GetAllTags retrieves every tag in the instance.
func (s *TagService) GetAllTags(ctx context.Context) ([]Tag, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameTag, "GetAllTags")
	logger.Trace().Msg("Retrieving all tags")

	allTags, err := NewTagDAO(s.db.Querier()).GetAllTags(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve tags from database")
		return nil, apperr.NewInternalError("failed to retrieve tags: %w", err)
	}

	logger.Debug().Int("count", len(allTags)).Msg("Tags retrieved successfully")
	return allTags, nil
}

// UpdateTag updates an existing tag.
func (s *TagService) UpdateTag(ctx context.Context, tag Tag) (*Tag, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameTag, "UpdateTag")
	logger.Trace().Interface("tag", tag).Msg("Updating tag")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewUnauthorizedError("user claims not found in context")
	}

	tag.Normalize()
	if err := tag.Validate(); err != nil {
		logger.Error().Err(err).Str("id", tag.ID.String()).Msg("Invalid tag")
		return nil, apperr.NewBadRequestError("invalid tag: %w", err)
	}

	updatedTag, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*Tag, *audit.AuditLogParams, error) {
		existingTag, err := NewTagDAO(q).GetTagByID(ctx, tag.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", tag.ID.String()).Msg("Tag not found")
				return nil, nil, apperr.NewNotFoundError("tag with ID %s not found", tag.ID)
			}
			logger.Error().Err(err).Str("id", tag.ID.String()).Msg("Failed to retrieve existing tag for update")
			return nil, nil, fmt.Errorf("failed to retrieve existing tag for update: %w", err)
		}

		t, updateErr := NewTagDAO(q).UpdateTag(ctx, tag)
		if updateErr != nil {
			logger.Error().Err(updateErr).Str("id", tag.ID.String()).Msg("Failed to update tag in database")
			return nil, nil, updateErr
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &t.ID,
			Action:     audit.ActionUpdateTag,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Tag with ID %s updated", t.ID),
			Before:     existingTag,
			After:      t,
		}
		return t, params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return nil, serviceErr
		}
		logger.Error().Err(txErr).Str("id", tag.ID.String()).Msg("Failed to update tag within transaction")
		return nil, apperr.NewInternalError("failed to update tag within transaction: %w", txErr)
	}

	logger.Debug().Str("id", updatedTag.ID.String()).Msg("Tag updated")
	return updatedTag, nil
}

// DeleteTag deletes a tag by its ID.
func (s *TagService) DeleteTag(ctx context.Context, id string) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameTag, "DeleteTag")
	logger.Trace().Str("id", id).Msg("Deleting tag")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return apperr.NewUnauthorizedError("user claims not found in context")
	}

	tagUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error().Str("id", id).Msg("Invalid tag ID format")
		return apperr.NewBadRequestError("invalid tag ID format: %w", err)
	}

	txErr := audit.WithAuditedVoidTx(ctx, s.db, func(q db.Querier) (*audit.AuditLogParams, error) {
		tagBeforeDelete, err := NewTagDAO(q).GetTagByID(ctx, tagUUID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", id).Msg("Tag not found in database")
				return nil, apperr.NewNotFoundError("tag with ID %s not found", id)
			}
			logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve tag before deletion")
			return nil, fmt.Errorf("failed to retrieve tag before deletion: %w", err)
		}

		if _, delErr := NewTagDAO(q).DeleteTagByID(ctx, tagUUID); delErr != nil {
			logger.Error().Err(delErr).Str("id", id).Msg("Failed to delete tag in database")
			return nil, delErr
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &tagUUID,
			Action:     audit.ActionDeleteTag,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Tag with ID %s deleted", tagUUID),
			Before:     tagBeforeDelete,
		}
		return params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return serviceErr
		}
		if errors.Is(txErr, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Tag not found or already deleted")
			return apperr.NewNotFoundError("tag not found or already deleted")
		}
		logger.Error().Err(txErr).Str("id", id).Msg("Failed to delete tag")
		return apperr.NewInternalError("failed to delete tag: %w", txErr)
	}

	logger.Debug().Str("id", id).Msg("Tag deleted")
	return nil
}
