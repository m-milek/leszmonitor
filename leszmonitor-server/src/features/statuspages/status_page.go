package statuspages

import (
	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/util"
)

type StatusPage struct {
	ID          uuid.UUID           `json:"id"          db:"id"`
	Slug        string              `json:"slug"        db:"slug"`
	Title       string              `json:"title"       db:"title"`
	Description string              `json:"description" db:"description"`
	Source      util.ResourceSource `json:"source"      db:"source"`
	util.Timestamps
}
