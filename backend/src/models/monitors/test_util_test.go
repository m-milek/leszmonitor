package monitors

import (
	"github.com/google/uuid"
	shared "github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/util"
)

// TestMonitor provides a simple way to create monitors for testing
type TestMonitor struct {
	// Base fields
	ID          uuid.UUID
	Slug        string
	Name        string
	Description string
	Interval    int
	ProjectID   string
	Type        shared.MonitorConfigType

	// Config fields
	HttpConfig *HttpConfig
	TCPConfig  *TCPConfig
}

// NewTestMonitor creates a new test monitor with default values
func NewTestMonitor() *TestMonitor {
	name := "Test Monitor"
	return &TestMonitor{
		ID:          uuid.Nil,
		Slug:        util.SlugFromString(name),
		Name:        name,
		Description: "Test monitor description",
		Interval:    60,
		ProjectID:   "test_owner",
		Type:        shared.HttpConfigType,
		HttpConfig: &HttpConfig{
			Method:              "GET",
			URL:                 "https://example.com",
			ExpectedStatusCodes: []int{200},
		},
	}
}

// AsHttp configures the monitor as an HTTP monitor
func (t *TestMonitor) AsHttp() *TestMonitor {
	t.Type = shared.HttpConfigType
	if t.HttpConfig == nil {
		t.HttpConfig = &HttpConfig{
			Method:              "GET",
			URL:                 "https://example.com",
			ExpectedStatusCodes: []int{200},
		}
	}
	t.TCPConfig = nil
	return t
}

// AsTCP configures the monitor as a tcpType monitor
func (t *TestMonitor) AsTCP() *TestMonitor {
	t.Type = shared.TCPConfigType
	if t.TCPConfig == nil {
		t.TCPConfig = &TCPConfig{
			Host:       "example.com",
			Port:       80,
			Protocol:   "tcp",
			Timeout:    5000,
			RetryCount: 3,
		}
	}
	t.HttpConfig = nil
	return t
}

// Build creates the monitor instance
func (t *TestMonitor) Build() IMonitor {
	base := BaseMonitor{
		ID:          t.ID,
		Slug:        t.Slug,
		Name:        t.Name,
		Description: t.Description,
		Interval:    t.Interval,
		Type:        t.Type,
	}

	switch t.Type {
	case shared.HttpConfigType:
		return &HttpMonitor{
			BaseMonitor: base,
			Config:      *t.HttpConfig,
		}
	case shared.TCPConfigType:
		return &TCPMonitor{
			BaseMonitor: base,
			Config:      *t.TCPConfig,
		}
	}
	return nil
}
