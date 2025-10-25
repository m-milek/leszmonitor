package db

// MockDB is a simple implementation of DB for tests.
// Provide your own fake repositories as needed.
type MockDB struct {
	UsersRepo    IUserRepository
	MonitorsRepo IMonitorRepository
	GroupsRepo   IGroupRepository
	TeamsRepo    ITeamRepository
	CloseFn      func()
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
