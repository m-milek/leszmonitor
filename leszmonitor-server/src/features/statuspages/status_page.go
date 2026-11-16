package statuspages

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/util"
)

type StatusPage struct {
	ID                 uuid.UUID           `json:"id"                 db:"id"`
	Slug               string              `json:"slug"               db:"slug"`
	Title              string              `json:"title"              db:"title"`
	Description        string              `json:"description"        db:"description"`
	GroupMonitorsByTag bool                `json:"groupMonitorsByTag" db:"group_monitors_by_tag"`
	RequireAuth        bool                `json:"requireAuth"        db:"require_auth"`
	Enabled            bool                `json:"enabled"            db:"enabled"`
	Source             util.ResourceSource `json:"source"             db:"source"`
	TagIDs             []uuid.UUID         `json:"tagIds"             db:"-"`
	MonitorCount       int                 `json:"monitorCount"       db:"monitor_count"`
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
