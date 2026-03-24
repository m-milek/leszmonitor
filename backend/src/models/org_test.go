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

func TestNewOrg_CreatesOrgWithOwner(t *testing.T) {
	owner := makeUUID(1)

	org, err := NewOrg("Alpha Org", "desc", owner)

	assert.NoError(t, err)
	if assert.NotNil(t, org) {
		assert.Equal(t, "Alpha Org", org.Name)
		assert.Equal(t, "alpha-org", org.DisplayID)
		assert.Equal(t, "desc", org.Description)
		if assert.Len(t, org.Members, 1) {
			assert.Equal(t, owner, org.Members[0].ID)
			assert.Equal(t, RoleOwner, org.Members[0].Role)
		}
		assert.True(t, org.IsMember(owner))
	}
}

func TestNewOrg_EmptyName_ReturnsError(t *testing.T) {
	owner := makeUUID(2)

	org, err := NewOrg("", "desc", owner)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "org name cannot be empty")
	assert.NotNil(t, org)
}

func TestOrg_IsMember_TrueAndFalse(t *testing.T) {
	a := makeUUID(10)
	b := makeUUID(11)
	c := makeUUID(12)
	org, err := NewOrg("X", "", a)
	assert.NoError(t, err)
	// add another member
	org.Members = append(org.Members, OrgMember{ID: b})

	assert.True(t, org.IsMember(a))
	assert.True(t, org.IsMember(b))
	assert.False(t, org.IsMember(c))
}

func TestOrg_GetMember_FoundAndNotFound(t *testing.T) {
	a := makeUUID(20)
	owner := makeUUID(21)
	org, err := NewOrg("X", "", owner)
	assert.NoError(t, err)
	// replace members to have a predictable single member
	org.Members = []OrgMember{{ID: a, Role: RoleMember}}

	m := org.GetMember(a)
	if assert.NotNil(t, m) {
		assert.Equal(t, a, m.ID)
		assert.Equal(t, RoleMember, m.Role)
	}

	notFound := org.GetMember(makeUUID(22))
	assert.Nil(t, notFound)
}

func TestOrg_ChangeMemberRole_Success(t *testing.T) {
	a := makeUUID(30)
	owner := makeUUID(31)
	org, err := NewOrg("X", "", owner)
	assert.NoError(t, err)
	org.Members = []OrgMember{{ID: a, Role: RoleMember}}

	err = org.ChangeMemberRole(a, RoleAdmin)

	assert.NoError(t, err)
	assert.Equal(t, RoleAdmin, org.Members[0].Role)
}

func TestOrg_ChangeMemberRole_NotMember_ReturnsError(t *testing.T) {
	a := makeUUID(40)
	owner := makeUUID(41)
	org, err := NewOrg("X", "", owner)
	assert.NoError(t, err)
	org.Members = []OrgMember{{ID: a, Role: RoleMember}}

	err = org.ChangeMemberRole(makeUUID(42), RoleAdmin)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "is not a member of the org")
	}
	assert.Equal(t, RoleMember, org.Members[0].Role)
}

func TestOrg_Validate_EmptyName_Error(t *testing.T) {
	owner := makeUUID(50)
	org, err := NewOrg("Org", "", owner)
	assert.NoError(t, err)
	// make name empty after construction
	org.DisplayIDFromName.Name = ""

	err = org.Validate()

	if assert.Error(t, err) {
		assert.Equal(t, "org name cannot be empty", err.Error())
	}
}

func TestOrg_Validate_NoMembers_Error(t *testing.T) {
	owner := makeUUID(55)
	org, err := NewOrg("Org", "", owner)
	assert.NoError(t, err)
	org.Members = []OrgMember{}

	err = org.Validate()

	if assert.Error(t, err) {
		assert.Equal(t, "org must have at least one member", err.Error())
	}
}

func TestOrg_Validate_InvalidRole_Error(t *testing.T) {
	owner := makeUUID(60)
	org, err := NewOrg("Org", "", owner)
	assert.NoError(t, err)
	org.Members = []OrgMember{{ID: makeUUID(61), Role: Role("invalid")}}

	err = org.Validate()

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid role for user 0: invalid org role: invalid")
	}
}

func TestOrg_Validate_Success(t *testing.T) {
	owner := makeUUID(70)
	org, err := NewOrg("Org", "", owner)
	assert.NoError(t, err)
	// org created by NewOrg already valid
	assert.NoError(t, org.Validate())
}
