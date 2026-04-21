package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/mock"
)

// MockDB is a simple implementation of DB for tests.
// Provide your own fake repositories as needed.
type MockDB struct {
	UsersRepo    IUserRepository
	MonitorsRepo IMonitorRepository
	ProjectsRepo IProjectRepository
	CloseFn      func()
}

type MockUserRepository struct {
	mock.Mock
}

func (r *MockUserRepository) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := r.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (r *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := r.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (r *MockUserRepository) GetUserByID(ctx context.Context, id pgtype.UUID) (*models.User, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (r *MockUserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	args := r.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

type MockProjectRepository struct {
	mock.Mock
}

func (r *MockProjectRepository) InsertProject(ctx context.Context, project *models.Project) error {
	args := r.Called(ctx, project)
	return args.Error(0)
}

func (r *MockProjectRepository) GetProjectBySlug(ctx context.Context, slug string) (*models.Project, error) {
	args := r.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (r *MockProjectRepository) GetProjectsByUserID(ctx context.Context, userID pgtype.UUID) ([]models.Project, error) {
	args := r.Called(ctx, userID)
	return args.Get(0).([]models.Project), args.Error(1)
}

func (r *MockProjectRepository) UpdateProject(ctx context.Context, oldProject, newProject *models.Project) (bool, error) {
	args := r.Called(ctx, oldProject, newProject)
	return args.Bool(0), args.Error(1)
}

func (r *MockProjectRepository) DeleteProject(ctx context.Context, projectSlug string) (bool, error) {
	args := r.Called(ctx, projectSlug)
	return args.Bool(0), args.Error(1)
}

func (r *MockProjectRepository) AddMemberToProject(ctx context.Context, projectSlug string, member *models.ProjectMember) (bool, error) {
	args := r.Called(ctx, projectSlug, member)
	return args.Bool(0), args.Error(1)
}

func (r *MockProjectRepository) RemoveMemberFromProject(ctx context.Context, projectSlug string, userID pgtype.UUID) (bool, error) {
	args := r.Called(ctx, projectSlug, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) Users() IUserRepository       { return m.UsersRepo }
func (m *MockDB) Monitors() IMonitorRepository { return m.MonitorsRepo }
func (m *MockDB) Projects() IProjectRepository { return m.ProjectsRepo }
func (m *MockDB) Close() {
	if m.CloseFn != nil {
		m.CloseFn()
	}
}
