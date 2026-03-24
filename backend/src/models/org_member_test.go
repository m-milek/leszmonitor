package models

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func makeMemberUUID(a byte) pgtype.UUID {
	var b [16]byte
	b[0] = a
	return pgtype.UUID{Bytes: b, Valid: true}
}

func TestNewOrgMember_Success(t *testing.T) {
	id := makeMemberUUID(1)

	member, err := NewOrgMember(id, RoleMember)

	assert.NoError(t, err)
	if assert.NotNil(t, member) {
		assert.Equal(t, id, member.ID)
		assert.Equal(t, RoleMember, member.Role)
	}
}

func TestNewOrgMember_InvalidID_ReturnsError(t *testing.T) {
	invalid := pgtype.UUID{} // Valid=false

	member, err := NewOrgMember(invalid, RoleMember)

	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "is not valid UUID")
	}
	assert.NotNil(t, member)
}

func TestNewOrgMember_InvalidRole_ReturnsError(t *testing.T) {
	id := makeMemberUUID(2)

	member, err := NewOrgMember(id, "invalid")

	if assert.Error(t, err) {
		assert.Equal(t, "invalid org role: invalid", err.Error())
	}
	assert.NotNil(t, member)
}

func TestOrgMember_Validate_Success(t *testing.T) {
	m := &OrgMember{ID: makeMemberUUID(3), Role: RoleAdmin}
	assert.NoError(t, m.Validate())
}
