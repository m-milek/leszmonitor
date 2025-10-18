package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeUUID(a byte) pgtype.UUID {
	var b [16]byte
	b[0] = a
	return pgtype.UUID{Bytes: b, Valid: true}
}

func TestNewTeam_CreatesTeamWithOwner(t *testing.T) {
	owner := makeUUID(1)

	team, err := NewTeam("Alpha Team", "desc", owner)

	assert.NoError(t, err)
	if assert.NotNil(t, team) {
		assert.Equal(t, "Alpha Team", team.Name)
		assert.Equal(t, "alpha-team", team.DisplayID)
		assert.Equal(t, "desc", team.Description)
		if assert.Len(t, team.Members, 1) {
			assert.Equal(t, owner, team.Members[0].ID)
			assert.Equal(t, RoleOwner, team.Members[0].Role)
		}
		assert.True(t, team.IsMember(owner))
	}
}

func TestNewTeam_EmptyName_ReturnsError(t *testing.T) {
	owner := makeUUID(2)

	team, err := NewTeam("", "desc", owner)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "team name cannot be empty")
	assert.NotNil(t, team)
}

func TestTeam_IsMember_TrueAndFalse(t *testing.T) {
	a := makeUUID(10)
	b := makeUUID(11)
	c := makeUUID(12)
	team, err := NewTeam("X", "", a)
	assert.NoError(t, err)
	// add another member
	team.Members = append(team.Members, TeamMember{ID: b})

	assert.True(t, team.IsMember(a))
	assert.True(t, team.IsMember(b))
	assert.False(t, team.IsMember(c))
}

func TestTeam_GetMember_FoundAndNotFound(t *testing.T) {
	a := makeUUID(20)
	owner := makeUUID(21)
	team, err := NewTeam("X", "", owner)
	assert.NoError(t, err)
	// replace members to have a predictable single member
	team.Members = []TeamMember{{ID: a, Role: RoleMember}}

	m := team.GetMember(a)
	if assert.NotNil(t, m) {
		assert.Equal(t, a, m.ID)
		assert.Equal(t, RoleMember, m.Role)
	}

	notFound := team.GetMember(makeUUID(22))
	assert.Nil(t, notFound)
}

func TestTeam_ChangeMemberRole_Success(t *testing.T) {
	a := makeUUID(30)
	owner := makeUUID(31)
	team, err := NewTeam("X", "", owner)
	assert.NoError(t, err)
	team.Members = []TeamMember{{ID: a, Role: RoleMember}}

	err = team.ChangeMemberRole(a, RoleAdmin)

	assert.NoError(t, err)
	assert.Equal(t, RoleAdmin, team.Members[0].Role)
}

func TestTeam_ChangeMemberRole_NotMember_ReturnsError(t *testing.T) {
	a := makeUUID(40)
	owner := makeUUID(41)
	team, err := NewTeam("X", "", owner)
	assert.NoError(t, err)
	team.Members = []TeamMember{{ID: a, Role: RoleMember}}

	err = team.ChangeMemberRole(makeUUID(42), RoleAdmin)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "is not a member of the team")
	}
	assert.Equal(t, RoleMember, team.Members[0].Role)
}

func TestTeam_Validate_EmptyName_Error(t *testing.T) {
	owner := makeUUID(50)
	team, err := NewTeam("Team", "", owner)
	assert.NoError(t, err)
	// make name empty after construction
	team.DisplayIDFromName.Name = ""

	err = team.Validate()

	if assert.Error(t, err) {
		assert.Equal(t, "team name cannot be empty", err.Error())
	}
}

func TestTeam_Validate_NoMembers_Error(t *testing.T) {
	owner := makeUUID(55)
	team, err := NewTeam("Team", "", owner)
	assert.NoError(t, err)
	team.Members = []TeamMember{}

	err = team.Validate()

	if assert.Error(t, err) {
		assert.Equal(t, "team must have at least one member", err.Error())
	}
}

func TestTeam_Validate_InvalidRole_Error(t *testing.T) {
	owner := makeUUID(60)
	team, err := NewTeam("Team", "", owner)
	assert.NoError(t, err)
	team.Members = []TeamMember{{ID: makeUUID(61), Role: Role("invalid")}}

	err = team.Validate()

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid role for user 0: invalid team role: invalid")
	}
}

func TestTeam_Validate_Success(t *testing.T) {
	owner := makeUUID(70)
	team, err := NewTeam("Team", "", owner)
	assert.NoError(t, err)
	// team created by NewTeam already valid
	assert.NoError(t, team.Validate())
}
