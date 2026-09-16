package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/stretchr/testify/mock"
)

// MockDB is a simple implementation of DB for tests.
type MockDB struct {
	UsersDAO                IUserDAO
	MonitorsDAO             IMonitorDAO
	MonitorResultsDAO       IMonitorResultDAO
	MonitorStatusChangesDAO IMonitorStatusChangeDAO
	MonitorStatsDao         IMonitorStatsDAO
	AuditLogDAO             IAuditLogDAO
	TagsDAO                 ITagDAO
	CloseFn                 func()
}

type MockUserDAO struct {
	mock.Mock
}

func (r *MockUserDAO) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := r.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (r *MockUserDAO) UpdateUserRole(
	ctx context.Context,
	userID uuid.UUID,
	role auth.Role,
) (*models.User, error) {
	args := r.Called(ctx, userID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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

type MockAuditLogDAO struct {
	mock.Mock
}

func (r *MockAuditLogDAO) InsertAuditLogEntry(ctx context.Context, entry audit.AuditLogEntry) (any, error) {
	args := r.Called(ctx, entry)
	return args.Get(0), args.Error(1)
}

func (r *MockAuditLogDAO) GetAuditLogEntries(
	ctx context.Context,
	filter audit.AuditLogFilter,
	pagination util.Pagination,
) ([]audit.AuditLogEntry, error) {
	args := r.Called(ctx, filter, pagination)
	return args.Get(0).([]audit.AuditLogEntry), args.Error(1)
}

func (r *MockAuditLogDAO) Record(ctx context.Context, params audit.AuditLogParams) error {
	args := r.Called(ctx, params)
	return args.Error(0)
}

type MockTagDAO struct {
	mock.Mock
}

func (r *MockTagDAO) InsertTag(ctx context.Context, tag models.Tag) (*models.Tag, error) {
	args := r.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (r *MockTagDAO) GetTagByID(ctx context.Context, id uuid.UUID) (*models.Tag, error) {
	args := r.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (r *MockTagDAO) GetAllTags(ctx context.Context) ([]models.Tag, error) {
	args := r.Called(ctx)
	return args.Get(0).([]models.Tag), args.Error(1)
}

func (r *MockTagDAO) UpdateTag(ctx context.Context, newTag models.Tag) (*models.Tag, error) {
	args := r.Called(ctx, newTag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (r *MockTagDAO) DeleteTagByID(ctx context.Context, tagID uuid.UUID) (*uuid.UUID, error) {
	args := r.Called(ctx, tagID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*uuid.UUID), args.Error(1)
}

func (m *MockDB) Users() IUserDAO                               { return m.UsersDAO }
func (m *MockDB) Monitors() IMonitorDAO                         { return m.MonitorsDAO }
func (m *MockDB) MonitorResults() IMonitorResultDAO             { return m.MonitorResultsDAO }
func (m *MockDB) MonitorStatusChanges() IMonitorStatusChangeDAO { return m.MonitorStatusChangesDAO }
func (m *MockDB) MonitorStats() IMonitorStatsDAO                { return m.MonitorStatsDao }
func (m *MockDB) AuditLog() IAuditLogDAO                        { return m.AuditLogDAO }
func (m *MockDB) Tags() ITagDAO                                 { return m.TagsDAO }
func (m *MockDB) WithTx(_ context.Context, fn func(tx DB) error) error {
	// In tests, execute the function directly without a real transaction.
	return fn(m)
}
func (m *MockDB) Close() {
	if m.CloseFn != nil {
		m.CloseFn()
	}
}
