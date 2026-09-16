package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/models"
	platformdb "github.com/m-milek/leszmonitor/platform/db"
)

type ITagDAO interface {
	InsertTag(ctx context.Context, tag models.Tag) (*models.Tag, error)
	GetTagByID(ctx context.Context, id uuid.UUID) (*models.Tag, error)
	GetAllTags(ctx context.Context) ([]models.Tag, error)
	UpdateTag(ctx context.Context, newTag models.Tag) (*models.Tag, error)
	DeleteTagByID(ctx context.Context, tagID uuid.UUID) (*uuid.UUID, error)
}

type tagDAO struct {
	baseDAO
}

func newTagDAO(base baseDAO) ITagDAO {
	return &tagDAO{
		baseDAO: base,
	}
}

// InsertTag adds a new tag to the database and returns the created tag.
func (r *tagDAO) InsertTag(ctx context.Context, tag models.Tag) (*models.Tag, error) {
	return platformdb.Wrap(ctx, "InsertTag", func() (*models.Tag, error) {
		id := tag.ID
		if id == uuid.Nil {
			id = uuid.New()
		}

		var createdTag models.Tag
		err := r.pool.QueryRowxContext(
			ctx,
			`INSERT INTO tags (id, name, description, color_hex)
			VALUES ($1, $2, $3, $4)
			RETURNING id, name, description, color_hex, created_at, updated_at`,
			id,
			tag.Name,
			tag.Description,
			tag.ColorHex,
		).StructScan(&createdTag)
		if err != nil {
			if platformdb.IsUniqueViolation(err) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		return &createdTag, nil
	})
}

func (r *tagDAO) GetTagByID(ctx context.Context, id uuid.UUID) (*models.Tag, error) {
	return platformdb.Wrap(ctx, "GetTagByID", func() (*models.Tag, error) {
		var tag models.Tag
		err := sqlx.GetContext(
			ctx,
			r.pool,
			&tag,
			`SELECT t.id, t.name, t.description, t.color_hex, t.created_at, t.updated_at
			 FROM tags t
			 WHERE t.id = $1`,
			id,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &tag, nil
	})
}

func (r *tagDAO) GetAllTags(ctx context.Context) ([]models.Tag, error) {
	return platformdb.Wrap(ctx, "GetAllTags", func() ([]models.Tag, error) {
		var tags []models.Tag
		err := sqlx.SelectContext(
			ctx,
			r.pool,
			&tags,
			`SELECT t.id, t.name, t.description, t.color_hex, t.created_at, t.updated_at
			 FROM tags t
			 ORDER BY t.name`,
		)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if tags == nil {
			tags = []models.Tag{}
		}
		return tags, nil
	})
}

func (r *tagDAO) UpdateTag(ctx context.Context, newTag models.Tag) (*models.Tag, error) {
	return platformdb.Wrap(ctx, "UpdateTag", func() (*models.Tag, error) {
		var updatedTag models.Tag
		// updated_at is set explicitly because the AFTER UPDATE trigger fires
		// after RETURNING has already captured the row.
		err := r.pool.QueryRowxContext(
			ctx,
			`UPDATE tags
			SET name=$1, description=$2, color_hex=$3, updated_at=CURRENT_TIMESTAMP
			WHERE id=$4
			RETURNING id, name, description, color_hex, created_at, updated_at`,
			newTag.Name,
			newTag.Description,
			newTag.ColorHex,
			newTag.ID,
		).StructScan(&updatedTag)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		return &updatedTag, nil
	})
}

func (r *tagDAO) DeleteTagByID(ctx context.Context, tagID uuid.UUID) (*uuid.UUID, error) {
	return platformdb.Wrap(ctx, "DeleteTagByID", func() (*uuid.UUID, error) {
		var id uuid.UUID
		err := r.pool.QueryRowxContext(ctx, `DELETE FROM tags WHERE id = $1 RETURNING id`, tagID).
			Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		return &id, nil
	})
}
