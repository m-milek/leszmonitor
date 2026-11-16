package statuspages

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/platform/db"
)

type IStatusPageDAO interface {
	InsertStatusPage(ctx context.Context, page StatusPage) (*StatusPage, error)
	GetStatusPageByID(ctx context.Context, id uuid.UUID) (*StatusPage, error)
	GetAllStatusPages(ctx context.Context) ([]StatusPage, error)
	UpdateStatusPage(ctx context.Context, page StatusPage) (*StatusPage, error)
	DeleteStatusPageByID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)
}

type statusPageDAO struct {
	pool db.Querier
}

func NewStatusPageDAO(pool db.Querier) IStatusPageDAO {
	return &statusPageDAO{
		pool: pool,
	}
}

func (r *statusPageDAO) InsertStatusPage(ctx context.Context, page StatusPage) (*StatusPage, error) {
	return db.Wrap(ctx, "InsertStatusPage", func() (*StatusPage, error) {
		id := page.ID
		if id == uuid.Nil {
			id = uuid.New()
		}

		_, err := r.pool.ExecContext(
			ctx,
			`INSERT INTO status_pages (id, slug, title, description, source)
			VALUES ($1, $2, $3, $4, $5)`,
			id,
			page.Slug,
			page.Title,
			page.Description,
			page.Source,
		)
		if err != nil {
			if db.IsUniqueViolation(err) {
				return nil, db.ErrAlreadyExists
			}
			return nil, err
		}

		if err = r.replaceStatusPageTags(ctx, id, page.TagIDs); err != nil {
			return nil, err
		}

		return r.GetStatusPageByID(ctx, id)
	})
}

func (r *statusPageDAO) GetStatusPageByID(ctx context.Context, id uuid.UUID) (*StatusPage, error) {
	return db.Wrap(ctx, "GetStatusPageByID", func() (*StatusPage, error) {
		var page StatusPage
		err := sqlx.GetContext(
			ctx,
			r.pool,
			&page,
			`SELECT p.id, p.slug, p.title, p.description, p.source, p.created_at, p.updated_at
			 FROM status_pages p
			 WHERE p.id = $1`,
			id,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		if err = r.attachTags(ctx, &page); err != nil {
			return nil, err
		}

		return &page, nil
	})
}

func (r *statusPageDAO) GetAllStatusPages(ctx context.Context) ([]StatusPage, error) {
	return db.Wrap(ctx, "GetAllStatusPages", func() ([]StatusPage, error) {
		var pages []StatusPage
		err := sqlx.SelectContext(
			ctx,
			r.pool,
			&pages,
			`SELECT p.id, p.slug, p.title, p.description, p.source, p.created_at, p.updated_at
			 FROM status_pages p
			 ORDER BY p.title`,
		)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if pages == nil {
			pages = []StatusPage{}
		}

		if err = r.attachTagsToAll(ctx, pages); err != nil {
			return nil, err
		}

		return pages, nil
	})
}

func (r *statusPageDAO) UpdateStatusPage(ctx context.Context, page StatusPage) (*StatusPage, error) {
	return db.Wrap(ctx, "UpdateStatusPage", func() (*StatusPage, error) {
		res, err := r.pool.ExecContext(
			ctx,
			`UPDATE status_pages
			SET slug=$1, title=$2, description=$3
			WHERE id=$4`,
			page.Slug,
			page.Title,
			page.Description,
			page.ID,
		)
		if err != nil {
			if db.IsUniqueViolation(err) {
				return nil, db.ErrAlreadyExists
			}
			return nil, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rowsAffected == 0 {
			return nil, db.ErrNotFound
		}

		if err = r.replaceStatusPageTags(ctx, page.ID, page.TagIDs); err != nil {
			return nil, err
		}

		return r.GetStatusPageByID(ctx, page.ID)
	})
}

func (r *statusPageDAO) DeleteStatusPageByID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	return db.Wrap(ctx, "DeleteStatusPageByID", func() (*uuid.UUID, error) {
		var deletedID uuid.UUID
		err := r.pool.QueryRowxContext(ctx, `DELETE FROM status_pages WHERE id = $1 RETURNING id`, id).
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

func (r *statusPageDAO) attachTags(ctx context.Context, page *StatusPage) error {
	tagIDs := []uuid.UUID{}
	err := sqlx.SelectContext(
		ctx,
		r.pool,
		&tagIDs,
		`SELECT tag_id FROM status_page_tags WHERE status_page_id = $1 ORDER BY tag_id`,
		page.ID,
	)
	if err != nil {
		return err
	}

	page.TagIDs = tagIDs
	return nil
}

func (r *statusPageDAO) attachTagsToAll(ctx context.Context, pages []StatusPage) error {
	var rows []struct {
		StatusPageID uuid.UUID `db:"status_page_id"`
		TagID        uuid.UUID `db:"tag_id"`
	}
	err := sqlx.SelectContext(
		ctx,
		r.pool,
		&rows,
		`SELECT status_page_id, tag_id FROM status_page_tags ORDER BY tag_id`,
	)
	if err != nil {
		return err
	}

	byPageID := make(map[uuid.UUID][]uuid.UUID, len(pages))
	for _, row := range rows {
		byPageID[row.StatusPageID] = append(byPageID[row.StatusPageID], row.TagID)
	}

	for i := range pages {
		tagIDs := byPageID[pages[i].ID]
		if tagIDs == nil {
			tagIDs = []uuid.UUID{}
		}
		pages[i].TagIDs = tagIDs
	}

	return nil
}

func (r *statusPageDAO) replaceStatusPageTags(ctx context.Context, pageID uuid.UUID, tagIDs []uuid.UUID) error {
	_, err := r.pool.ExecContext(ctx, `DELETE FROM status_page_tags WHERE status_page_id = $1`, pageID)
	if err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		_, err = r.pool.ExecContext(
			ctx,
			`INSERT OR IGNORE INTO status_page_tags (status_page_id, tag_id) VALUES ($1, $2)`,
			pageID,
			tagID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
