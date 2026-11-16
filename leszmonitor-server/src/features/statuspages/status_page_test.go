package statuspages_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/statuspages"
	"github.com/stretchr/testify/assert"
)

func validStatusPage() statuspages.StatusPage {
	return statuspages.StatusPage{
		Slug:              "public-status",
		Title:             "Public Status",
		TimeWindowSeconds: 3600,
		TagIDs:            []uuid.UUID{uuid.New()},
	}
}

func TestStatusPage_Validate(t *testing.T) {
	t.Run("Accepts a valid status page", func(t *testing.T) {
		page := validStatusPage()
		assert.NoError(t, page.Validate())
	})

	tests := []struct {
		name        string
		modify      func(p *statuspages.StatusPage)
		errContains string
	}{
		{"Empty slug", func(p *statuspages.StatusPage) { p.Slug = "" }, "slug cannot be empty"},
		{"Invalid slug", func(p *statuspages.StatusPage) { p.Slug = "Bad Slug" }, "invalid status page slug"},
		{"Empty title", func(p *statuspages.StatusPage) { p.Title = "" }, "title cannot be empty"},
		{"Zero time window", func(p *statuspages.StatusPage) { p.TimeWindowSeconds = 0 }, "time window must be positive"},
		{"No tags", func(p *statuspages.StatusPage) { p.TagIDs = nil }, "at least one tag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := validStatusPage()
			tt.modify(&page)
			assert.ErrorContains(t, page.Validate(), tt.errContains)
		})
	}
}
