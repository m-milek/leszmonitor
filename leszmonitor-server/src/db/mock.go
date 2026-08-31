package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/util"
	"github.com/stretchr/testify/mock"
)

// MockDB is a simple implementation of DB for tests.
type MockDB struct {
	UsersDAO                IUserDAO
	MonitorsDAO             IMonitorDAO
	ProjectsDAO             IProjectDAO
	MonitorResultsDAO       IMonitorResultDAO
	MonitorStatusChangesDAO IMonitorStatusChangeDAO
	MonitorStatsDao         IMonitorStatsDAO
	AuditLogDAO             IAuditLogDAO
	CloseFn                 func()
}

type MockUserDAO struct {
	mock.Mock
}

func (r *MockUserDAO) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := r.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (r *MockUserDAO) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := r.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (r *MockUserDAO) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (r *MockUserDAO) GetAllUsers(ctx context.Context) ([]models.User, error) {
	args := r.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

type MockProjectDAO struct {
	mock.Mock
}

func (r *MockProjectDAO) InsertProject(ctx context.Context, project *models.Project) error {
	args := r.Called(ctx, project)
	return args.Error(0)
}

func (r *MockProjectDAO) GetProjectBySlug(ctx context.Context, slug string) (*models.Project, error) {
	args := r.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (r *MockProjectDAO) GetProjectByID(ctx context.Context, slug uuid.UUID) (*models.Project, error) {
	args := r.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (r *MockProjectDAO) GetProjectsByQuery(
	ctx context.Context,
	query GetProjectsQuery,
) ([]models.Project, error) {
	args := r.Called(ctx, query)
	return args.Get(0).([]models.Project), args.Error(1)
}

func (r *MockProjectDAO) UpdateProject(
	ctx context.Context,
	oldProject, newProject *models.Project,
) (bool, error) {
	args := r.Called(ctx, oldProject, newProject)
	return args.Bool(0), args.Error(1)
}

func (r *MockProjectDAO) DeleteProject(ctx context.Context, projectSlug string) (bool, error) {
	args := r.Called(ctx, projectSlug)
	return args.Bool(0), args.Error(1)
}

func (r *MockProjectDAO) AddMemberToProject(
	ctx context.Context,
	projectSlug string,
	member *models.ProjectMember,
) (bool, error) {
	args := r.Called(ctx, projectSlug, member)
	return args.Bool(0), args.Error(1)
}

func (r *MockProjectDAO) RemoveMemberFromProject(
	ctx context.Context,
	projectSlug string,
	userID uuid.UUID,
) (bool, error) {
	args := r.Called(ctx, projectSlug, userID)
	return args.Bool(0), args.Error(1)
}

func (r *MockProjectDAO) ChangeMemberRole(
	ctx context.Context,
	projectSlug string,
	userID uuid.UUID,
	newRole models.Role,
) (bool, error) {
	args := r.Called(ctx, projectSlug, userID, newRole)
	return args.Bool(0), args.Error(1)
}

type MockAuditLogDAO struct {
	mock.Mock
}

func (r *MockAuditLogDAO) InsertAuditLogEntry(ctx context.Context, entry security.AuditLogEntry) (any, error) {
	args := r.Called(ctx, entry)
	return args.Get(0), args.Error(1)
}

func (r *MockAuditLogDAO) GetAuditLogEntries(
	ctx context.Context,
	filter security.AuditLogFilter,
	pagination util.Pagination,
) ([]security.AuditLogEntry, error) {
	args := r.Called(ctx, filter, pagination)
	return args.Get(0).([]security.AuditLogEntry), args.Error(1)
}

func (r *MockAuditLogDAO) Record(ctx context.Context, params security.AuditLogParams) error {
	args := r.Called(ctx, params)
	return args.Error(0)
}

func (m *MockDB) Users() IUserDAO                               { return m.UsersDAO }
func (m *MockDB) Monitors() IMonitorDAO                         { return m.MonitorsDAO }
func (m *MockDB) Projects() IProjectDAO                         { return m.ProjectsDAO }
func (m *MockDB) MonitorResults() IMonitorResultDAO             { return m.MonitorResultsDAO }
func (m *MockDB) MonitorStatusChanges() IMonitorStatusChangeDAO { return m.MonitorStatusChangesDAO }
func (m *MockDB) MonitorStats() IMonitorStatsDAO                { return m.MonitorStatsDao }
func (m *MockDB) AuditLog() IAuditLogDAO                        { return m.AuditLogDAO }
func (m *MockDB) WithTx(_ context.Context, fn func(tx DB) error) error {
	// In tests, execute the function directly without a real transaction.
	return fn(m)
}
func (m *MockDB) Close() {
	if m.CloseFn != nil {
		m.CloseFn()
	}
}
