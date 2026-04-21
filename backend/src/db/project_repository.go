package db

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models"
)

type IProjectRepository interface {
	InsertProject(ctx context.Context, project *models.Project) error
	GetProjectByID(ctx context.Context, id pgtype.UUID) (*models.Project, error)
	GetProjectByDisplayID(ctx context.Context, displayID string) (*models.Project, error)
	GetProjectsByUserID(ctx context.Context, userID pgtype.UUID) ([]models.Project, error)
	UpdateProject(ctx context.Context, oldProject, newProject *models.Project) (bool, error)
	DeleteProject(ctx context.Context, projectDisplayID string) (bool, error)
	AddMemberToProject(ctx context.Context, projectDisplayID string, member *models.ProjectMember) (bool, error)
	RemoveMemberFromProject(ctx context.Context, projectDisplayID string, userID pgtype.UUID) (bool, error)
}

type projectRepository struct {
	baseRepository
}

// projectMemberFromCollectableRow maps a pgx.CollectableRow to a models.ProjectMember.
func projectMemberFromCollectableRow(row pgx.CollectableRow) (models.ProjectMember, error) {
	member := models.ProjectMember{}
	err := row.Scan(&member.ID, &member.Username, &member.Role, &member.CreatedAt, &member.UpdatedAt)
	return member, err
}

// projectFromCollectableRow maps a pgx.CollectableRow to a models.Project struct (without members).
func projectFromCollectableRow(row pgx.CollectableRow) (models.Project, error) {
	project := models.Project{}
	err := row.Scan(&project.ID, &project.DisplayID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)
	return project, err
}

func newProjectRepository(repository baseRepository) IProjectRepository {
	return &projectRepository{
		baseRepository: repository,
	}
}

// loadMembers fetches members for a single project and attaches them.
func (r *projectRepository) loadMembers(ctx context.Context, project *models.Project) error {
	memberRows, err := r.pool.Query(ctx,
		`SELECT up.user_id, u.username, up.role, up.created_at, up.updated_at
		 FROM user_projects up JOIN users u ON u.id = up.user_id
		 WHERE up.project_id = $1`,
		project.ID)
	if err != nil {
		return err
	}
	members, err := pgx.CollectRows(memberRows, projectMemberFromCollectableRow)
	if err != nil {
		return err
	}
	project.Members = members
	return nil
}

func (r *projectRepository) InsertProject(ctx context.Context, project *models.Project) error {
	_, err := dbWrap(ctx, "InsertProject", func() (*any, error) {
		tx, err := r.pool.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer func() { _ = tx.Rollback(ctx) }()

		var projectID pgtype.UUID
		row := tx.QueryRow(ctx,
			`INSERT INTO projects (display_id, name, description) VALUES ($1, $2, $3) RETURNING id`,
			project.DisplayID, project.Name, project.Description)
		if err = row.Scan(&projectID); err != nil {
			if pgErrIs(err, pgerrcode.UniqueViolation) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		// Insert the owner member
		_, err = tx.Exec(ctx,
			`INSERT INTO user_projects (project_id, user_id, role) VALUES ($1, $2, $3)`,
			projectID, project.Members[0].ID, project.Members[0].Role)
		if err != nil {
			return nil, err
		}

		return nil, tx.Commit(ctx)
	})
	if pgErrIs(err, pgerrcode.UniqueViolation) {
		return ErrAlreadyExists
	}
	return err
}

func (r *projectRepository) GetProjectByID(ctx context.Context, id pgtype.UUID) (*models.Project, error) {
	return dbWrap(ctx, "GetProjectByID", func() (*models.Project, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT id, display_id, name, description, created_at, updated_at FROM projects WHERE id = $1`,
			id)
		if err != nil {
			return nil, err
		}

		project, err := pgx.CollectExactlyOneRow(rows, projectFromCollectableRow)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		if err = r.loadMembers(ctx, &project); err != nil {
			return nil, err
		}

		return &project, nil
	})
}

func (r *projectRepository) GetProjectByDisplayID(ctx context.Context, displayID string) (*models.Project, error) {
	return dbWrap(ctx, "GetProjectByDisplayID", func() (*models.Project, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT id, display_id, name, description, created_at, updated_at FROM projects WHERE display_id = $1`,
			displayID)
		if err != nil {
			return nil, err
		}

		project, err := pgx.CollectExactlyOneRow(rows, projectFromCollectableRow)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		if err = r.loadMembers(ctx, &project); err != nil {
			return nil, err
		}

		return &project, nil
	})
}

func (r *projectRepository) GetProjectsByUserID(ctx context.Context, userID pgtype.UUID) ([]models.Project, error) {
	return dbWrap(ctx, "GetProjectsByUserID", func() ([]models.Project, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT p.id, p.display_id, p.name, p.description, p.created_at, p.updated_at
			 FROM projects p
			 JOIN user_projects up ON up.project_id = p.id
			 WHERE up.user_id = $1`,
			userID)
		if err != nil {
			return nil, err
		}

		projects, err := pgx.CollectRows(rows, projectFromCollectableRow)
		if err != nil {
			return nil, err
		}

		for i := range projects {
			if err = r.loadMembers(ctx, &projects[i]); err != nil {
				return nil, err
			}
		}

		return projects, nil
	})
}

func (r *projectRepository) UpdateProject(ctx context.Context, oldProject, newProject *models.Project) (bool, error) {
	return dbWrap(ctx, "UpdateProject", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`UPDATE projects SET display_id = $1, name = $2, description = $3 WHERE id = $4`,
			newProject.DisplayID, newProject.Name, newProject.Description, oldProject.ID)
		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}

func (r *projectRepository) DeleteProject(ctx context.Context, projectDisplayID string) (bool, error) {
	return dbWrap(ctx, "DeleteProject", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`DELETE FROM projects WHERE display_id = $1`,
			projectDisplayID)
		if err != nil {
			return false, err
		}
		return result.RowsAffected() > 0, nil
	})
}

func (r *projectRepository) AddMemberToProject(ctx context.Context, projectDisplayID string, member *models.ProjectMember) (bool, error) {
	return dbWrap(ctx, "AddMemberToProject", func() (bool, error) {
		var projectID pgtype.UUID
		err := r.pool.QueryRow(ctx,
			`SELECT id FROM projects WHERE display_id = $1`, projectDisplayID).Scan(&projectID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := r.pool.Exec(ctx,
			`INSERT INTO user_projects (project_id, user_id, role) VALUES ($1, $2, $3)`,
			projectID, member.ID, member.Role)
		if err != nil {
			if pgErrIs(err, pgerrcode.UniqueViolation) {
				return false, ErrAlreadyExists
			}
			return false, err
		}
		return result.RowsAffected() > 0, nil
	})
}

func (r *projectRepository) RemoveMemberFromProject(ctx context.Context, projectDisplayID string, userID pgtype.UUID) (bool, error) {
	return dbWrap(ctx, "RemoveMemberFromProject", func() (bool, error) {
		var projectID pgtype.UUID
		err := r.pool.QueryRow(ctx,
			`SELECT id FROM projects WHERE display_id = $1`, projectDisplayID).Scan(&projectID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := r.pool.Exec(ctx,
			`DELETE FROM user_projects WHERE project_id = $1 AND user_id = $2`,
			projectID, userID)
		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}
