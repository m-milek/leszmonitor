package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeMemberUUID(a byte) uuid.UUID {
	var b [16]byte
	b[0] = a
	return uuid.UUID(b)
}

func TestNewProjectMember_Success(t *testing.T) {
	t.Parallel()
	id := makeMemberUUID(1)

	member, err := NewProjectMember(id, RoleMember)

	require.NoError(t, err)
	if assert.NotNil(t, member) {
		assert.Equal(t, id, member.ID)
		assert.Equal(t, RoleMember, member.Role)
	}
}

func TestNewProjectMember_InvalidID_ReturnsError(t *testing.T) {
	t.Parallel()
	invalid := uuid.Nil // invalid ID

	member, err := NewProjectMember(invalid, RoleMember)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not valid UUID")
	assert.NotNil(t, member)
}

func TestNewProjectMember_InvalidRole_ReturnsError(t *testing.T) {
	t.Parallel()
	id := makeMemberUUID(2)

	member, err := NewProjectMember(id, "invalid")

	require.Error(t, err)
	assert.Equal(t, "invalid project role: invalid", err.Error())
	assert.NotNil(t, member)
}

func TestProjectMember_Validate_Success(t *testing.T) {
	t.Parallel()
	m := &ProjectMember{ID: makeMemberUUID(3), Role: RoleAdmin}
	assert.NoError(t, m.Validate())
}
