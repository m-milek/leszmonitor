package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeMemberUUID(a byte) pgtype.UUID {
	var b [16]byte
	b[0] = a
	return pgtype.UUID{Bytes: b, Valid: true}
}

func TestNewTeamMember_Success(t *testing.T) {
	id := makeMemberUUID(1)

	member, err := NewTeamMember(id, RoleMember)

	assert.NoError(t, err)
	if assert.NotNil(t, member) {
		assert.Equal(t, id, member.ID)
		assert.Equal(t, RoleMember, member.Role)
	}
}

func TestNewTeamMember_InvalidID_ReturnsError(t *testing.T) {
	invalid := pgtype.UUID{} // Valid=false

	member, err := NewTeamMember(invalid, RoleMember)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "is not valid UUID")
	}
	assert.NotNil(t, member)
}

func TestNewTeamMember_InvalidRole_ReturnsError(t *testing.T) {
	id := makeMemberUUID(2)

	member, err := NewTeamMember(id, "invalid")

	if assert.Error(t, err) {
		assert.Equal(t, "invalid team role: invalid", err.Error())
	}
	assert.NotNil(t, member)
}

func TestTeamMember_Validate_Success(t *testing.T) {
	m := &TeamMember{ID: makeMemberUUID(3), Role: RoleAdmin}
	assert.NoError(t, m.Validate())
}
