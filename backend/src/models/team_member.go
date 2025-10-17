package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type TeamMember struct {
	Id   pgtype.UUID `json:"id"`
	Role TeamRole    `json:"role"`
	Timestamps
}

func (tm *TeamMember) Validate() error {
	return tm.Role.Validate()
}
