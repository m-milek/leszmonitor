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
	GroupsRepo   IGroupRepository
	TeamsRepo    ITeamRepository
	CloseFn      func()
}

type MockUserRepository struct {
	mock.Mock
}

func (r *MockUserRepository) InsertUser(ctx context.Context, user *models.User) (*struct{}, error) {
	args := r.Called(ctx, user)
	return args.Get(0).(*struct{}), args.Error(1)
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

type MockGroupRepository struct {
	mock.Mock
}

func (r *MockGroupRepository) InsertGroup(ctx context.Context, group *models.MonitorGroup) error {
	args := r.Called(ctx, group)
	return args.Error(0)
}

func (r *MockGroupRepository) GetGroupByDisplayID(ctx context.Context, displayID string) (*models.MonitorGroup, error) {
	args := r.Called(ctx, displayID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MonitorGroup), args.Error(1)
}

func (r *MockGroupRepository) GetGroupsByTeamID(ctx context.Context, team *models.Team) ([]models.MonitorGroup, error) {
	args := r.Called(ctx, team)
	return args.Get(0).([]models.MonitorGroup), args.Error(1)
}

func (r *MockGroupRepository) UpdateGroup(ctx context.Context, team *models.Team, oldGroup, newGroup *models.MonitorGroup) (bool, error) {
	args := r.Called(ctx, team, oldGroup, newGroup)
	return args.Bool(0), args.Error(1)
}

func (r *MockGroupRepository) DeleteGroup(ctx context.Context, team *models.Team, groupID string) (bool, error) {
	args := r.Called(ctx, team, groupID)
	return args.Bool(0), args.Error(1)
}

type MockTeamRepository struct {
	mock.Mock
}

func (r *MockTeamRepository) InsertTeam(ctx context.Context, team *models.Team) (*struct{}, error) {
	args := r.Called(ctx, team)
	return args.Get(0).(*struct{}), args.Error(1)
}

func (r *MockTeamRepository) GetTeamByDisplayID(ctx context.Context, displayID string) (*models.Team, error) {
	args := r.Called(ctx, displayID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (r *MockTeamRepository) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	args := r.Called(ctx)
	return args.Get(0).([]models.Team), args.Error(1)
}

func (r *MockTeamRepository) DeleteTeamByID(ctx context.Context, displayID string) (bool, error) {
	args := r.Called(ctx, displayID)
	return args.Bool(0), args.Error(1)
}

func (r *MockTeamRepository) UpdateTeam(ctx context.Context, team *models.Team) (bool, error) {
	args := r.Called(ctx, team)
	return args.Bool(0), args.Error(1)
}

func (r *MockTeamRepository) AddMemberToTeam(ctx context.Context, teamID string, member *models.TeamMember) (bool, error) {
	args := r.Called(ctx, teamID, member)
	return args.Bool(0), args.Error(1)
}

func (r *MockTeamRepository) RemoveMemberFromTeam(ctx context.Context, teamID string, userID pgtype.UUID) (bool, error) {
	args := r.Called(ctx, teamID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) Users() IUserRepository       { return m.UsersRepo }
func (m *MockDB) Monitors() IMonitorRepository { return m.MonitorsRepo }
func (m *MockDB) Groups() IGroupRepository     { return m.GroupsRepo }
func (m *MockDB) Teams() ITeamRepository       { return m.TeamsRepo }
func (m *MockDB) Close() {
	if m.CloseFn != nil {
		m.CloseFn()
	}
}
