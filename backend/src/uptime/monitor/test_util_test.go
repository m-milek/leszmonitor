package monitors

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/util"
)

// TestMonitor provides a simple way to create monitors for testing
type TestMonitor struct {
	// Base fields
	ID          pgtype.UUID
	DisplayID   string
	Name        string
	Description string
	Interval    int
	GroupID     string
	Type        MonitorConfigType

	// Config fields
	HttpConfig *httpConfig
	PingConfig *PingConfig
}

// NewTestMonitor creates a new test monitor with default values
func NewTestMonitor() *TestMonitor {
	name := "Test Monitor"
	return &TestMonitor{
		ID:          pgtype.UUID{Valid: false},
		DisplayID:   util.IDFromString(name),
		Name:        name,
		Description: "Test monitor description",
		Interval:    60,
		GroupID:     "test_owner",
		Type:        httpType,
		HttpConfig: &httpConfig{
			Method:              "GET",
			URL:                 "https://example.com",
			ExpectedStatusCodes: []int{200},
		},
	}
}

// AsHttp configures the monitor as an HTTP monitor
func (t *TestMonitor) AsHttp() *TestMonitor {
	t.Type = httpType
	if t.HttpConfig == nil {
		t.HttpConfig = &httpConfig{
			Method:              "GET",
			URL:                 "https://example.com",
			ExpectedStatusCodes: []int{200},
		}
	}
	t.PingConfig = nil
	return t
}

// AsPing configures the monitor as a pingType monitor
func (t *TestMonitor) AsPing() *TestMonitor {
	t.Type = pingType
	if t.PingConfig == nil {
		t.PingConfig = &PingConfig{
			Host:        "example.com",
			Port:        "80",
			Protocol:    "tcp",
			PingTimeout: 5,
			RetryCount:  3,
		}
	}
	t.HttpConfig = nil
	return t
}

// Build creates the monitor instance
func (t *TestMonitor) Build() IMonitor {
	base := BaseMonitor{
		ID:          t.ID,
		DisplayID:   t.DisplayID,
		Name:        t.Name,
		Description: t.Description,
		Interval:    t.Interval,
		Type:        t.Type,
	}

	switch t.Type {
	case httpType:
		return &httpMonitor{
			BaseMonitor: base,
			Config:      *t.HttpConfig,
		}
	case pingType:
		return &PingMonitor{
			BaseMonitor: base,
			Config:      *t.PingConfig,
		}
	}
	return nil
}
