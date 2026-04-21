package models

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func makeUUID(a byte) pgtype.UUID {
	var b [16]byte
	b[0] = a
	return pgtype.UUID{Bytes: b, Valid: true}
}

func TestNewProject_CreatesProjectWithOwner(t *testing.T) {
	owner := makeUUID(1)

	project, err := NewProject("Alpha Project", "desc", owner)

	assert.NoError(t, err)
	if assert.NotNil(t, project) {
		assert.Equal(t, "Alpha Project", project.Name)
		assert.Equal(t, "alpha-project", project.Slug)
		assert.Equal(t, "desc", project.Description)
		if assert.Len(t, project.Members, 1) {
			assert.Equal(t, owner, project.Members[0].ID)
			assert.Equal(t, RoleOwner, project.Members[0].Role)
		}
		assert.True(t, project.IsMember(owner))
	}
}

func TestNewProject_EmptyName_ReturnsError(t *testing.T) {
	owner := makeUUID(2)

	project, err := NewProject("", "desc", owner)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project name cannot be empty")
	assert.Nil(t, project)
}

func TestProject_IsMember_TrueAndFalse(t *testing.T) {
	a := makeUUID(10)
	b := makeUUID(11)
	c := makeUUID(12)
	project, err := NewProject("X", "", a)
	assert.NoError(t, err)
	project.Members = append(project.Members, ProjectMember{ID: b})

	assert.True(t, project.IsMember(a))
	assert.True(t, project.IsMember(b))
	assert.False(t, project.IsMember(c))
}

func TestProject_GetMember_FoundAndNotFound(t *testing.T) {
	a := makeUUID(20)
	owner := makeUUID(21)
	project, err := NewProject("X", "", owner)
	assert.NoError(t, err)
	project.Members = []ProjectMember{{ID: a, Role: RoleMember}}

	m := project.GetMember(a)
	if assert.NotNil(t, m) {
		assert.Equal(t, a, m.ID)
		assert.Equal(t, RoleMember, m.Role)
	}

	notFound := project.GetMember(makeUUID(22))
	assert.Nil(t, notFound)
}

func TestProject_ChangeMemberRole_Success(t *testing.T) {
	a := makeUUID(30)
	owner := makeUUID(31)
	project, err := NewProject("X", "", owner)
	assert.NoError(t, err)
	project.Members = []ProjectMember{{ID: a, Role: RoleMember}}

	err = project.ChangeMemberRole(a, RoleAdmin)

	assert.NoError(t, err)
	assert.Equal(t, RoleAdmin, project.Members[0].Role)
}

func TestProject_ChangeMemberRole_NotMember_ReturnsError(t *testing.T) {
	a := makeUUID(40)
	owner := makeUUID(41)
	project, err := NewProject("X", "", owner)
	assert.NoError(t, err)
	project.Members = []ProjectMember{{ID: a, Role: RoleMember}}

	err = project.ChangeMemberRole(makeUUID(42), RoleAdmin)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "is not a member of the project")
	}
	assert.Equal(t, RoleMember, project.Members[0].Role)
}

func TestProject_Validate_EmptyName_Error(t *testing.T) {
	owner := makeUUID(50)
	project, err := NewProject("Proj", "", owner)
	assert.NoError(t, err)
	project.SlugFromName.Name = ""

	err = project.Validate()

	if assert.Error(t, err) {
		assert.Equal(t, "project name cannot be empty", err.Error())
	}
}

func TestProject_Validate_NoMembers_Error(t *testing.T) {
	owner := makeUUID(55)
	project, err := NewProject("Proj", "", owner)
	assert.NoError(t, err)
	project.Members = []ProjectMember{}

	err = project.Validate()

	if assert.Error(t, err) {
		assert.Equal(t, "project must have at least one member", err.Error())
	}
}

func TestProject_Validate_InvalidRole_Error(t *testing.T) {
	owner := makeUUID(60)
	project, err := NewProject("Proj", "", owner)
	assert.NoError(t, err)
	project.Members = []ProjectMember{{ID: makeUUID(61), Role: Role("invalid")}}

	err = project.Validate()

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid project role: invalid")
	}
}

func TestProject_Validate_Success(t *testing.T) {
	owner := makeUUID(70)
	project, err := NewProject("Proj", "", owner)
	assert.NoError(t, err)
	assert.NoError(t, project.Validate())
}
