package statuspages

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/util"
)

type StatusPage struct {
	ID                 uuid.UUID           `json:"id"                 yaml:"id"                 db:"id"`
	Slug               string              `json:"slug"               yaml:"slug"               db:"slug"`
	Title              string              `json:"title"              yaml:"title"              db:"title"`
	Description        string              `json:"description"        yaml:"description"        db:"description"`
	GroupMonitorsByTag bool                `json:"groupMonitorsByTag" yaml:"groupMonitorsByTag" db:"group_monitors_by_tag"`
	RequireAuth        bool                `json:"requireAuth"        yaml:"requireAuth"        db:"require_auth"`
	Enabled            bool                `json:"enabled"            yaml:"enabled"            db:"enabled"`
	TimeWindowSeconds  int                 `json:"timeWindowSeconds"  yaml:"timeWindowSeconds"  db:"time_window_seconds"`
	Source             util.ResourceSource `json:"source"             yaml:"source"             db:"source"`
	TagIDs             []uuid.UUID         `json:"tagIds"             yaml:"tagIds"             db:"-"`
	MonitorCount       int                 `json:"monitorCount"       yaml:"monitorCount"       db:"monitor_count"`
	util.Timestamps    `yaml:",inline"`
}

func (p *StatusPage) Validate() error {
	if p.Slug == "" {
		return fmt.Errorf("status page slug cannot be empty")
	}
	if p.Slug != util.SlugFromString(p.Slug) {
		return fmt.Errorf("invalid status page slug: %s", p.Slug)
	}
	if p.Title == "" {
		return fmt.Errorf("status page title cannot be empty")
	}
	if p.TimeWindowSeconds <= 0 {
		return fmt.Errorf("status page time window must be positive")
	}
	if len(p.TagIDs) == 0 {
		return fmt.Errorf("status page must have at least one tag")
	}
	return nil
}
