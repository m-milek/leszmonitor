package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/models"
)

type IProjectRepository interface {
	InsertProject(ctx context.Context, project *models.Project) error
	GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
	GetProjectBySlug(ctx context.Context, slug string) (*models.Project, error)
	GetProjectsByQuery(ctx context.Context, query GetProjectsQuery) ([]models.Project, error)
	UpdateProject(ctx context.Context, oldProject, newProject *models.Project) (bool, error)
	DeleteProject(ctx context.Context, projectSlug string) (bool, error)
	AddMemberToProject(ctx context.Context, projectSlug string, member *models.ProjectMember) (bool, error)
	RemoveMemberFromProject(ctx context.Context, projectSlug string, userID uuid.UUID) (bool, error)
}

type projectRepository struct {
	baseRepository
}

func newProjectRepository(repository baseRepository) IProjectRepository {
	return &projectRepository{
		baseRepository: repository,
	}
}

// loadMembers fetches members for a single project and attaches them.
func (r *projectRepository) loadMembers(ctx context.Context, project *models.Project) error {
	var members []models.ProjectMember
	err := sqlx.SelectContext(ctx, r.pool, &members,
		`SELECT up.user_id AS id, u.username AS username, up.role AS role, up.created_at, up.updated_at
		 FROM user_projects up JOIN users u ON u.id = up.user_id
		 WHERE up.project_id = $1`,
		project.ID)
	if err != nil {
		return err
	}
	project.Members = members
	return nil
}

func (r *projectRepository) InsertProject(ctx context.Context, project *models.Project) error {
	_, err := dbWrap(ctx, "InsertProject", func() (*any, error) {
		if project.ID == uuid.Nil {
			project.ID = uuid.New()
		}

		var projectID uuid.UUID
		row := r.pool.QueryRowxContext(ctx,
			`INSERT INTO projects (id, slug, name, description) VALUES ($1, $2, $3, $4) RETURNING id`,
			project.ID, project.Slug, project.Name, project.Description)
		if err := row.Scan(&projectID); err != nil {
			if isUniqueViolation(err) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}
		project.ID = projectID

		// Insert the owner member
		_, err := r.pool.ExecContext(ctx,
			`INSERT INTO user_projects (project_id, user_id, role) VALUES ($1, $2, $3)`,
			projectID, project.Members[0].ID, project.Members[0].Role)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if isUniqueViolation(err) {
		return ErrAlreadyExists
	}
	return err
}

func (r *projectRepository) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return dbWrap(ctx, "GetProjectByID", func() (*models.Project, error) {
		var project models.Project
		err := sqlx.GetContext(ctx, r.pool, &project,
			`SELECT id, slug, name, description, created_at, updated_at FROM projects WHERE id = $1`,
			id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
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

func (r *projectRepository) GetProjectBySlug(ctx context.Context, slug string) (*models.Project, error) {
	return dbWrap(ctx, "GetProjectBySlug", func() (*models.Project, error) {
		var project models.Project
		err := sqlx.GetContext(ctx, r.pool, &project,
			`SELECT id, slug, name, description, created_at, updated_at FROM projects WHERE slug = $1`,
			slug)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
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

type GetProjectsQuery struct {
	RequestingUserID uuid.UUID
	MemberUsername   string
}

func (r *projectRepository) GetProjectsByQuery(ctx context.Context, query GetProjectsQuery) ([]models.Project, error) {
	return dbWrap(ctx, "GetProjectsByQuery", func() ([]models.Project, error) {
		var projects []models.Project

		if query.MemberUsername == "" {
			// Return all projects the requesting user belongs to
			err := sqlx.SelectContext(ctx, r.pool, &projects,
				`SELECT p.id, p.slug, p.name, p.description, p.created_at, p.updated_at
			 FROM projects p
			 JOIN user_projects up ON up.project_id = p.id
			 WHERE up.user_id = ?`,
				query.RequestingUserID)
			if err != nil {
				return nil, err
			}
		} else {
			// Return projects where BOTH the requesting user and the target user are members
			err := sqlx.SelectContext(ctx, r.pool, &projects,
				`SELECT p.id, p.slug, p.name, p.description, p.created_at, p.updated_at
			 FROM projects p
			 JOIN user_projects requester ON requester.project_id = p.id
			 JOIN user_projects member ON member.project_id = p.id
			 JOIN users u ON u.id = member.user_id
			 WHERE requester.user_id = ?
			   AND u.username = ?`,
				query.RequestingUserID, query.MemberUsername)
			if err != nil {
				return nil, err
			}
		}

		for i := range projects {
			if err := r.loadMembers(ctx, &projects[i]); err != nil {
				return nil, err
			}
		}

		if projects == nil {
			projects = []models.Project{}
		}

		return projects, nil
	})
}

func (r *projectRepository) UpdateProject(ctx context.Context, oldProject, newProject *models.Project) (bool, error) {
	return dbWrap(ctx, "UpdateProject", func() (bool, error) {
		result, err := r.pool.ExecContext(ctx,
			`UPDATE projects SET slug = $1, name = $2, description = $3 WHERE id = $4`,
			newProject.Slug, newProject.Name, newProject.Description, oldProject.ID)
		if err != nil {
			return false, err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return false, err
		}
		if rows == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}

func (r *projectRepository) DeleteProject(ctx context.Context, projectSlug string) (bool, error) {
	return dbWrap(ctx, "DeleteProject", func() (bool, error) {
		result, err := r.pool.ExecContext(ctx,
			`DELETE FROM projects WHERE slug = $1`,
			projectSlug)
		if err != nil {
			return false, err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return false, err
		}
		return rows > 0, nil
	})
}

func (r *projectRepository) AddMemberToProject(ctx context.Context, projectSlug string, member *models.ProjectMember) (bool, error) {
	return dbWrap(ctx, "AddMemberToProject", func() (bool, error) {
		var projectID uuid.UUID
		err := r.pool.QueryRowxContext(ctx,
			`SELECT id FROM projects WHERE slug = $1`, projectSlug).Scan(&projectID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := r.pool.ExecContext(ctx,
			`INSERT INTO user_projects (project_id, user_id, role) VALUES ($1, $2, $3)`,
			projectID, member.ID, member.Role)
		if err != nil {
			if isUniqueViolation(err) {
				return false, ErrAlreadyExists
			}
			return false, err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return false, err
		}
		return rows > 0, nil
	})
}

func (r *projectRepository) RemoveMemberFromProject(ctx context.Context, projectSlug string, userID uuid.UUID) (bool, error) {
	return dbWrap(ctx, "RemoveMemberFromProject", func() (bool, error) {
		var projectID uuid.UUID
		err := r.pool.QueryRowxContext(ctx,
			`SELECT id FROM projects WHERE slug = $1`, projectSlug).Scan(&projectID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := r.pool.ExecContext(ctx,
			`DELETE FROM user_projects WHERE project_id = $1 AND user_id = $2`,
			projectID, userID)
		if err != nil {
			return false, err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return false, err
		}
		if rows == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}
