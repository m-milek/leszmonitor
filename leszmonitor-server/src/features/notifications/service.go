package notifications

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
)

type INotificationChannelService interface {
	CreateChannel(ctx context.Context, channel NotificationChannel) (*NotificationChannel, *apperr.ServiceError)
	GetChannelByID(ctx context.Context, id uuid.UUID) (*NotificationChannel, *apperr.ServiceError)
	GetAllChannels(ctx context.Context) ([]NotificationChannel, *apperr.ServiceError)
	UpdateChannel(ctx context.Context, channel NotificationChannel) (*NotificationChannel, *apperr.ServiceError)
	DeleteChannel(ctx context.Context, id uuid.UUID) *apperr.ServiceError
	ExportToYAML(ctx context.Context, id uuid.UUID, w io.Writer) *apperr.ServiceError
}

type NotificationChannelService struct {
	db db.DB
}

type NotificationChannelServiceDeps struct {
	DB db.DB
}

func NewNotificationChannelService(deps NotificationChannelServiceDeps) *NotificationChannelService {
	return &NotificationChannelService{
		db: deps.DB,
	}
}

func (s *NotificationChannelService) CreateChannel(
	ctx context.Context,
	channel NotificationChannel,
) (*NotificationChannel, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameNotificationChannel, "CreateChannel")
	logger.Trace().Str("name", channel.Name).Str("kind", string(channel.Kind)).Msg("Creating notification channel")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewInternalError("user claims not found in context")
	}

	newChannel := NotificationChannel{
		ID:      uuid.New(),
		Name:    channel.Name,
		Kind:    channel.Kind,
		Enabled: channel.Enabled,
		Config:  channel.Config,
		Source:  util.ResourceSourceUI,
	}
	if err := newChannel.prepare(nil); err != nil {
		logger.Error().Err(err).Msg("Invalid notification channel")
		return nil, apperr.NewBadRequestError("invalid notification channel: %w", err)
	}

	created, txErr := audit.WithAuditedTx(
		ctx,
		s.db,
		func(q db.Querier) (*NotificationChannel, *audit.AuditLogParams, error) {
			c, createErr := NewNotificationChannelDAO(q).InsertChannel(ctx, newChannel)
			if createErr != nil {
				return nil, nil, createErr
			}

			redacted, redactErr := c.redacted()
			if redactErr != nil {
				return nil, nil, redactErr
			}

			params := &audit.AuditLogParams{
				Username:   &userClaims.Username,
				ResourceID: &c.ID,
				Action:     audit.ActionCreateNotificationChannel,
				IsSuccess:  true,
				Summary:    fmt.Sprintf("Notification channel with ID %s created", c.ID),
				After:      redacted,
			}
			return &redacted, params, nil
		},
	)
	if txErr != nil {
		if errors.Is(txErr, db.ErrAlreadyExists) {
			return nil, apperr.NewConflictError("notification channel with name %s already exists", newChannel.Name)
		}
		logger.Error().Err(txErr).Msg("Failed to create notification channel within transaction")
		return nil, apperr.NewInternalError("failed to create notification channel within transaction: %w", txErr)
	}

	logger.Debug().Str("id", created.ID.String()).Msg("Notification channel created")
	return created, nil
}

func (s *NotificationChannelService) GetChannelByID(
	ctx context.Context,
	id uuid.UUID,
) (*NotificationChannel, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameNotificationChannel, "GetChannelByID")
	logger.Trace().Str("id", id.String()).Msg("Retrieving notification channel by ID")

	channel, err := NewNotificationChannelDAO(s.db.Querier()).GetChannelByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", id.String()).Msg("Notification channel not found in database")
			return nil, apperr.NewNotFoundError("notification channel with ID %s not found", id)
		}
		logger.Error().Err(err).Str("id", id.String()).Msg("Failed to retrieve notification channel from database")
		return nil, apperr.NewInternalError("failed to retrieve notification channel: %w", err)
	}

	redacted, err := channel.redacted()
	if err != nil {
		logger.Error().Err(err).Str("id", id.String()).Msg("Failed to redact notification channel")
		return nil, apperr.NewInternalError("failed to redact notification channel: %w", err)
	}

	logger.Debug().Str("id", id.String()).Msg("Notification channel retrieved successfully")
	return &redacted, nil
}

func (s *NotificationChannelService) GetAllChannels(ctx context.Context) ([]NotificationChannel, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameNotificationChannel, "GetAllChannels")
	logger.Trace().Msg("Retrieving all notification channels")

	channels, err := NewNotificationChannelDAO(s.db.Querier()).GetAllChannels(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve notification channels from database")
		return nil, apperr.NewInternalError("failed to retrieve notification channels: %w", err)
	}

	for i, channel := range channels {
		redacted, redactErr := channel.redacted()
		if redactErr != nil {
			logger.Error().Err(redactErr).Str("id", channel.ID.String()).Msg("Failed to redact notification channel")
			return nil, apperr.NewInternalError("failed to redact notification channel: %w", redactErr)
		}
		channels[i] = redacted
	}

	logger.Debug().Int("count", len(channels)).Msg("Notification channels retrieved successfully")
	return channels, nil
}

func (s *NotificationChannelService) UpdateChannel(
	ctx context.Context,
	channel NotificationChannel,
) (*NotificationChannel, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameNotificationChannel, "UpdateChannel")
	logger.Trace().Str("id", channel.ID.String()).Msg("Updating notification channel")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewInternalError("user claims not found in context")
	}

	updated, txErr := audit.WithAuditedTx(
		ctx,
		s.db,
		func(q db.Querier) (*NotificationChannel, *audit.AuditLogParams, error) {
			existing, err := NewNotificationChannelDAO(q).GetChannelByID(ctx, channel.ID)
			if err != nil {
				if errors.Is(err, db.ErrNotFound) {
					logger.Error().Str("id", channel.ID.String()).Msg("Notification channel not found")
					return nil, nil, apperr.NewNotFoundError("notification channel with ID %s not found", channel.ID)
				}
				logger.Error().Err(err).Str("id", channel.ID.String()).Msg("Failed to retrieve existing notification channel for update")
				return nil, nil, fmt.Errorf("failed to retrieve existing notification channel for update: %w", err)
			}

			if existing.Source == util.ResourceSourceConfig {
				return nil, nil, apperr.NewConflictError("notification channel %s is managed by config file", channel.ID)
			}

			if prepareErr := channel.prepare(existing); prepareErr != nil {
				logger.Error().Err(prepareErr).Str("id", channel.ID.String()).Msg("Invalid notification channel")
				return nil, nil, apperr.NewBadRequestError("invalid notification channel: %w", prepareErr)
			}

			c, updateErr := NewNotificationChannelDAO(q).UpdateChannel(ctx, channel)
			if updateErr != nil {
				logger.Error().Err(updateErr).Str("id", channel.ID.String()).Msg("Failed to update notification channel in database")
				return nil, nil, updateErr
			}

			redactedBefore, redactErr := existing.redacted()
			if redactErr != nil {
				return nil, nil, redactErr
			}
			redactedAfter, redactErr := c.redacted()
			if redactErr != nil {
				return nil, nil, redactErr
			}

			params := &audit.AuditLogParams{
				Username:   &userClaims.Username,
				ResourceID: &c.ID,
				Action:     audit.ActionUpdateNotificationChannel,
				IsSuccess:  true,
				Summary:    fmt.Sprintf("Notification channel with ID %s updated", c.ID),
				Before:     redactedBefore,
				After:      redactedAfter,
			}
			return &redactedAfter, params, nil
		},
	)
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return nil, serviceErr
		}
		if errors.Is(txErr, db.ErrAlreadyExists) {
			return nil, apperr.NewConflictError("notification channel with name %s already exists", channel.Name)
		}
		logger.Error().Err(txErr).Str("id", channel.ID.String()).Msg("Failed to update notification channel within transaction")
		return nil, apperr.NewInternalError("failed to update notification channel within transaction: %w", txErr)
	}

	logger.Debug().Str("id", updated.ID.String()).Msg("Notification channel updated")
	return updated, nil
}

func (s *NotificationChannelService) DeleteChannel(ctx context.Context, id uuid.UUID) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameNotificationChannel, "DeleteChannel")
	logger.Trace().Str("id", id.String()).Msg("Deleting notification channel")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return apperr.NewInternalError("user claims not found in context")
	}

	txErr := audit.WithAuditedVoidTx(ctx, s.db, func(q db.Querier) (*audit.AuditLogParams, error) {
		existing, err := NewNotificationChannelDAO(q).GetChannelByID(ctx, id)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", id.String()).Msg("Notification channel not found in database")
				return nil, apperr.NewNotFoundError("notification channel with ID %s not found", id)
			}
			logger.Error().Err(err).Str("id", id.String()).Msg("Failed to retrieve notification channel before deletion")
			return nil, fmt.Errorf("failed to retrieve notification channel before deletion: %w", err)
		}

		if existing.Source == util.ResourceSourceConfig {
			return nil, apperr.NewConflictError("notification channel %s is managed by config file", id)
		}

		if _, delErr := NewNotificationChannelDAO(q).DeleteChannelByID(ctx, id); delErr != nil {
			logger.Error().Err(delErr).Str("id", id.String()).Msg("Failed to delete notification channel in database")
			return nil, delErr
		}

		redacted, redactErr := existing.redacted()
		if redactErr != nil {
			return nil, redactErr
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &id,
			Action:     audit.ActionDeleteNotificationChannel,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Notification channel with ID %s deleted", id),
			Before:     redacted,
		}
		return params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return serviceErr
		}
		if errors.Is(txErr, db.ErrNotFound) {
			logger.Error().Str("id", id.String()).Msg("Notification channel not found or already deleted")
			return apperr.NewNotFoundError("notification channel not found or already deleted")
		}
		logger.Error().Err(txErr).Str("id", id.String()).Msg("Failed to delete notification channel")
		return apperr.NewInternalError("failed to delete notification channel: %w", txErr)
	}

	logger.Debug().Str("id", id.String()).Msg("Notification channel deleted")
	return nil
}

func (s *NotificationChannelService) ExportToYAML(ctx context.Context, id uuid.UUID, w io.Writer) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameNotificationChannel, "ExportToYAML")
	logger.Trace().Str("id", id.String()).Msg("Exporting to YAML")

	channel, err := NewNotificationChannelDAO(s.db.Querier()).GetChannelByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apperr.NewNotFoundError("notification channel with ID %s not found", id)
		}
		return apperr.NewInternalError("failed to retrieve notification channel: %w", err)
	}

	if err = exportToYAML(*channel, w); err != nil {
		logger.Error().Err(err).Str("id", id.String()).Msg("Failed to export notification channel to YAML")
		return apperr.NewInternalError("failed to export notification channel %s to YAML", channel.Name)
	}

	return nil
}
