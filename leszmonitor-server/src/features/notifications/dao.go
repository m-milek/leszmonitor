package notifications

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/util"
)

type INotificationChannelDAO interface {
	InsertChannel(ctx context.Context, channel NotificationChannel) (*NotificationChannel, error)
	GetChannelByID(ctx context.Context, id uuid.UUID) (*NotificationChannel, error)
	GetAllChannels(ctx context.Context) ([]NotificationChannel, error)
	UpdateChannel(ctx context.Context, channel NotificationChannel) (*NotificationChannel, error)
	DeleteChannelByID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)
	UpsertConfigBasedChannel(ctx context.Context, channel *NotificationChannel) error
}

type notificationChannelDAO struct {
	pool db.Querier
}

func NewNotificationChannelDAO(pool db.Querier) INotificationChannelDAO {
	return &notificationChannelDAO{
		pool: pool,
	}
}

func (r *notificationChannelDAO) InsertChannel(ctx context.Context, channel NotificationChannel) (*NotificationChannel, error) {
	return db.Wrap(ctx, "InsertChannel", func() (*NotificationChannel, error) {
		var created NotificationChannel
		err := r.pool.QueryRowxContext(
			ctx,
			`INSERT INTO notification_channels (id, name, kind, enabled, config, source)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, name, kind, enabled, config, source, created_at, updated_at`,
			channel.ID,
			channel.Name,
			channel.Kind,
			channel.Enabled,
			channel.Config,
			channel.Source,
		).StructScan(&created)
		if err != nil {
			if db.IsUniqueViolation(err) {
				return nil, db.ErrAlreadyExists
			}
			return nil, err
		}

		return &created, nil
	})
}

func (r *notificationChannelDAO) GetChannelByID(ctx context.Context, id uuid.UUID) (*NotificationChannel, error) {
	return db.Wrap(ctx, "GetChannelByID", func() (*NotificationChannel, error) {
		var channel NotificationChannel
		err := sqlx.GetContext(
			ctx,
			r.pool,
			&channel,
			`SELECT id, name, kind, enabled, config, source, created_at, updated_at
			 FROM notification_channels
			 WHERE id = $1`,
			id,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}
		return &channel, nil
	})
}

func (r *notificationChannelDAO) GetAllChannels(ctx context.Context) ([]NotificationChannel, error) {
	return db.Wrap(ctx, "GetAllChannels", func() ([]NotificationChannel, error) {
		var channels []NotificationChannel
		err := sqlx.SelectContext(
			ctx,
			r.pool,
			&channels,
			`SELECT id, name, kind, enabled, config, source, created_at, updated_at
			 FROM notification_channels
			 ORDER BY name`,
		)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if channels == nil {
			channels = []NotificationChannel{}
		}
		return channels, nil
	})
}

func (r *notificationChannelDAO) UpdateChannel(ctx context.Context, channel NotificationChannel) (*NotificationChannel, error) {
	return db.Wrap(ctx, "UpdateChannel", func() (*NotificationChannel, error) {
		var updated NotificationChannel
		err := r.pool.QueryRowxContext(
			ctx,
			`UPDATE notification_channels
			SET name=$1, kind=$2, enabled=$3, config=$4, updated_at=CURRENT_TIMESTAMP
			WHERE id=$5
			RETURNING id, name, kind, enabled, config, source, created_at, updated_at`,
			channel.Name,
			channel.Kind,
			channel.Enabled,
			channel.Config,
			channel.ID,
		).StructScan(&updated)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			if db.IsUniqueViolation(err) {
				return nil, db.ErrAlreadyExists
			}
			return nil, err
		}

		return &updated, nil
	})
}

func (r *notificationChannelDAO) DeleteChannelByID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	return db.Wrap(ctx, "DeleteChannelByID", func() (*uuid.UUID, error) {
		var deletedID uuid.UUID
		err := r.pool.QueryRowxContext(ctx, `DELETE FROM notification_channels WHERE id = $1 RETURNING id`, id).
			Scan(&deletedID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		return &deletedID, nil
	})
}

func (r *notificationChannelDAO) UpsertConfigBasedChannel(ctx context.Context, channel *NotificationChannel) error {
	_, err := db.Wrap(ctx, "UpsertConfigBasedChannel", func() (any, error) {
		_, err := r.pool.ExecContext(
			ctx,
			`INSERT INTO notification_channels (id, name, kind, enabled, config, source)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE SET
				kind = excluded.kind,
				enabled = excluded.enabled,
				config = excluded.config`,
			channel.ID,
			channel.Name,
			channel.Kind,
			channel.Enabled,
			channel.Config,
			util.ResourceSourceConfig,
		)
		if err != nil {
			if db.IsUniqueViolation(err) {
				return nil, db.ErrAlreadyExists
			}
			return nil, err
		}

		return nil, nil
	})
	return err
}
