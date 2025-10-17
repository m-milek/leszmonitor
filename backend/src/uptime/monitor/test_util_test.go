package monitors

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/util"
)

// TestMonitor provides a simple way to create monitors for testing
type TestMonitor struct {
	// Base fields
	ID          pgtype.UUID
	DisplayId   string
	Name        string
	Description string
	Interval    int
	GroupId     string
	Type        MonitorConfigType

	// Config fields
	HttpConfig *HttpConfig
	PingConfig *PingConfig
}

// NewTestMonitor creates a new test monitor with default values
func NewTestMonitor() *TestMonitor {
	name := "Test Monitor"
	return &TestMonitor{
		ID:          pgtype.UUID{Valid: false},
		DisplayId:   util.IdFromString(name),
		Name:        name,
		Description: "Test monitor description",
		Interval:    60,
		GroupId:     "test_owner",
		Type:        Http,
		HttpConfig: &HttpConfig{
			HttpMethod:          "GET",
			Url:                 "https://example.com",
			ExpectedStatusCodes: []int{200},
		},
	}
}

// AsHttp configures the monitor as an HTTP monitor
func (t *TestMonitor) AsHttp() *TestMonitor {
	t.Type = Http
	if t.HttpConfig == nil {
		t.HttpConfig = &HttpConfig{
			HttpMethod:          "GET",
			Url:                 "https://example.com",
			ExpectedStatusCodes: []int{200},
		}
	}
	t.PingConfig = nil
	return t
}

// AsPing configures the monitor as a Ping monitor
func (t *TestMonitor) AsPing() *TestMonitor {
	t.Type = Ping
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
		Id:          t.ID,
		DisplayId:   t.DisplayId,
		Name:        t.Name,
		Description: t.Description,
		Interval:    t.Interval,
		Type:        t.Type,
	}

	switch t.Type {
	case Http:
		return &HttpMonitor{
			BaseMonitor: base,
			Config:      *t.HttpConfig,
		}
	case Ping:
		return &PingMonitor{
			BaseMonitor: base,
			Config:      *t.PingConfig,
		}
	}
	return nil
}
