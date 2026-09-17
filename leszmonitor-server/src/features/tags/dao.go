package tags

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/platform/db"
)

type ITagDAO interface {
	InsertTag(ctx context.Context, tag Tag) (*Tag, error)
	GetTagByID(ctx context.Context, id uuid.UUID) (*Tag, error)
	GetAllTags(ctx context.Context) ([]Tag, error)
	UpdateTag(ctx context.Context, newTag Tag) (*Tag, error)
	DeleteTagByID(ctx context.Context, tagID uuid.UUID) (*uuid.UUID, error)
}

type tagDAO struct {
	pool db.Querier
}

func NewTagDAO(pool db.Querier) ITagDAO {
	return &tagDAO{
		pool: pool,
	}
}

// InsertTag adds a new tag to the database and returns the created tag.
func (r *tagDAO) InsertTag(ctx context.Context, tag Tag) (*Tag, error) {
	return db.Wrap(ctx, "InsertTag", func() (*Tag, error) {
		id := tag.ID
		if id == uuid.Nil {
			id = uuid.New()
		}

		var createdTag Tag
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
			if db.IsUniqueViolation(err) {
				return nil, db.ErrAlreadyExists
			}
			return nil, err
		}

		return &createdTag, nil
	})
}

func (r *tagDAO) GetTagByID(ctx context.Context, id uuid.UUID) (*Tag, error) {
	return db.Wrap(ctx, "GetTagByID", func() (*Tag, error) {
		var tag Tag
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
				return nil, db.ErrNotFound
			}
			return nil, err
		}
		return &tag, nil
	})
}

func (r *tagDAO) GetAllTags(ctx context.Context) ([]Tag, error) {
	return db.Wrap(ctx, "GetAllTags", func() ([]Tag, error) {
		var tags []Tag
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
			tags = []Tag{}
		}
		return tags, nil
	})
}

func (r *tagDAO) UpdateTag(ctx context.Context, newTag Tag) (*Tag, error) {
	return db.Wrap(ctx, "UpdateTag", func() (*Tag, error) {
		var updatedTag Tag
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
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		return &updatedTag, nil
	})
}

func (r *tagDAO) DeleteTagByID(ctx context.Context, tagID uuid.UUID) (*uuid.UUID, error) {
	return db.Wrap(ctx, "DeleteTagByID", func() (*uuid.UUID, error) {
		var id uuid.UUID
		err := r.pool.QueryRowxContext(ctx, `DELETE FROM tags WHERE id = $1 RETURNING id`, tagID).
			Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		return &id, nil
	})
}
