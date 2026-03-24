package db

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
)

type IProjectRepository interface {
	InsertProject(ctx context.Context, project *models.Project) error
	GetProjectByDisplayID(ctx context.Context, displayID string) (*models.Project, error)
	GetProjectsByOrgID(ctx context.Context, org *models.Org) ([]models.Project, error)
	UpdateProject(ctx context.Context, org *models.Org, oldProject, newProject *models.Project) (bool, error)
	DeleteProject(ctx context.Context, org *models.Org, projectID string) (bool, error)
}

type projectRepository struct {
	baseRepository
}

// projectFromCollectableRow maps a pgx.CollectableRow to a models.Project struct.
func projectFromCollectableRow(row pgx.CollectableRow) (models.Project, error) {
	project := models.Project{}
	err := row.Scan(&project.ID, &project.OrgID, &project.DisplayID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

	return project, err
}

func newProjectRepository(repository baseRepository) IProjectRepository {
	return &projectRepository{
		baseRepository: repository,
	}
}

func (r *projectRepository) InsertProject(ctx context.Context, project *models.Project) error {
	_, err := dbWrap(ctx, "InsertProject", func() (*any, error) {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO projects (org_id, display_id, name, description) VALUES ($1, $2, $3, $4)`,
			project.OrgID, project.DisplayID, project.Name, project.Description)
		return nil, err
	})
	if pgErrIs(err, pgerrcode.UniqueViolation) {
		return ErrAlreadyExists
	}
	return err
}

func (r *projectRepository) GetProjectByDisplayID(ctx context.Context, displayID string) (*models.Project, error) {
	return dbWrap(ctx, "GetProjectByDisplayID", func() (*models.Project, error) {
		row, err := r.pool.Query(ctx,
			`SELECT id, org_id, display_id, name, description, created_at, updated_at FROM projects WHERE display_id=$1`,
			displayID)

		if err != nil {
			return nil, err
		}

		project, err := pgx.CollectExactlyOneRow(row, projectFromCollectableRow)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return &project, err
	})
}

func (r *projectRepository) GetProjectsByOrgID(ctx context.Context, org *models.Org) ([]models.Project, error) {
	return dbWrap(ctx, "GetProjectsByOrgID", func() ([]models.Project, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT id, org_id, display_id, name, description, created_at, updated_at FROM projects WHERE org_id=$1`,
			org.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logging.Db.Info().Msgf("No projects found for org %s", org.DisplayID)
				return []models.Project{}, nil
			}
			return nil, err
		}

		var projects []models.Project
		projects, err = pgx.CollectRows(rows, projectFromCollectableRow)
		if err != nil {
			return nil, err
		}
		return projects, nil
	})
}

func (r *projectRepository) UpdateProject(ctx context.Context, org *models.Org, oldProject, newProject *models.Project) (bool, error) {
	return dbWrap(ctx, "UpdateProject", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`UPDATE projects SET display_id=$1, name=$2, description=$3 WHERE id=$4 AND org_id=$5 RETURNING *`,
			newProject.DisplayID, newProject.Name, newProject.Description, oldProject.ID, org.ID)

		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}

		return true, nil
	})
}

func (r *projectRepository) DeleteProject(ctx context.Context, org *models.Org, projectID string) (bool, error) {
	return dbWrap(ctx, "DeleteProject", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`DELETE FROM projects WHERE display_id=$1 AND org_id=$2`,
			projectID, org.ID)
		if err != nil {
			return false, err
		}
		return result.RowsAffected() > 0, nil
	})
}
