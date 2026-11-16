package statuspages

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/util"
)

type StatusPage struct {
	ID          uuid.UUID           `json:"id"          db:"id"`
	Slug        string              `json:"slug"        db:"slug"`
	Title       string              `json:"title"       db:"title"`
	Description string              `json:"description" db:"description"`
	Source      util.ResourceSource `json:"source"      db:"source"`
	TagIDs      []uuid.UUID         `json:"tagIds"      db:"-"`
	util.Timestamps
}

func (p *StatusPage) Validate() error {
	if p.Slug == "" {
		return fmt.Errorf("status page slug cannot be empty")
	}
	if p.Title == "" {
		return fmt.Errorf("status page title cannot be empty")
	}
	return nil
}
