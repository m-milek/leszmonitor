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
	OrgsRepo     IOrgRepository
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

func (r *MockProjectRepository) GetProjectByDisplayID(ctx context.Context, displayID string) (*models.Project, error) {
	args := r.Called(ctx, displayID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (r *MockProjectRepository) GetProjectsByOrgID(ctx context.Context, org *models.Org) ([]models.Project, error) {
	args := r.Called(ctx, org)
	return args.Get(0).([]models.Project), args.Error(1)
}

func (r *MockProjectRepository) UpdateProject(ctx context.Context, org *models.Org, oldProject, newProject *models.Project) (bool, error) {
	args := r.Called(ctx, org, oldProject, newProject)
	return args.Bool(0), args.Error(1)
}

func (r *MockProjectRepository) DeleteProject(ctx context.Context, org *models.Org, projectID string) (bool, error) {
	args := r.Called(ctx, org, projectID)
	return args.Bool(0), args.Error(1)
}

type MockOrgRepository struct {
	mock.Mock
}

func (r *MockOrgRepository) InsertOrg(ctx context.Context, org *models.Org) (*struct{}, error) {
	args := r.Called(ctx, org)
	return args.Get(0).(*struct{}), args.Error(1)
}

func (r *MockOrgRepository) GetOrgByDisplayID(ctx context.Context, displayID string) (*models.Org, error) {
	args := r.Called(ctx, displayID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Org), args.Error(1)
}

func (r *MockOrgRepository) GetAllOrgs(ctx context.Context) ([]models.Org, error) {
	args := r.Called(ctx)
	return args.Get(0).([]models.Org), args.Error(1)
}

func (r *MockOrgRepository) DeleteOrgByID(ctx context.Context, displayID string) (bool, error) {
	args := r.Called(ctx, displayID)
	return args.Bool(0), args.Error(1)
}

func (r *MockOrgRepository) UpdateOrg(ctx context.Context, org *models.Org) (bool, error) {
	args := r.Called(ctx, org)
	return args.Bool(0), args.Error(1)
}

func (r *MockOrgRepository) AddMemberToOrg(ctx context.Context, orgID string, member *models.OrgMember) (bool, error) {
	args := r.Called(ctx, orgID, member)
	return args.Bool(0), args.Error(1)
}

func (r *MockOrgRepository) RemoveMemberFromOrg(ctx context.Context, orgID string, userID pgtype.UUID) (bool, error) {
	args := r.Called(ctx, orgID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) Users() IUserRepository       { return m.UsersRepo }
func (m *MockDB) Monitors() IMonitorRepository { return m.MonitorsRepo }
func (m *MockDB) Projects() IProjectRepository { return m.ProjectsRepo }
func (m *MockDB) Orgs() IOrgRepository         { return m.OrgsRepo }
func (m *MockDB) Close() {
	if m.CloseFn != nil {
		m.CloseFn()
	}
}
