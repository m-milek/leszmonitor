package monitors

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestMonitor provides a simple way to create monitors for testing
type TestMonitor struct {
	// Base fields
	ID          string
	ObjectID    bson.ObjectID
	Name        string
	Description string
	Interval    int
	OwnerID     string
	Type        MonitorConfigType

	// Config fields
	HttpConfig *HttpConfig
	PingConfig *PingConfig
}

// NewTestMonitor creates a new test monitor with default values
func NewTestMonitor() *TestMonitor {
	return &TestMonitor{
		ID:          generateId(),
		ObjectID:    bson.NewObjectID(),
		Name:        "Test Monitor",
		Description: "Test monitor description",
		Interval:    60,
		OwnerID:     "test_owner",
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
		Name:        t.Name,
		Description: t.Description,
		Interval:    t.Interval,
		OwnerId:     t.OwnerID,
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

// ToBSON converts the monitor to a BSON document
func (t *TestMonitor) ToBSON() bson.M {
	doc := bson.M{
		"_id":         t.ObjectID,
		"id":          t.ID,
		"name":        t.Name,
		"description": t.Description,
		"interval":    t.Interval,
		"ownerId":     t.OwnerID,
		"type":        string(t.Type),
	}

	switch t.Type {
	case Http:
		doc["config"] = t.HttpConfig
	case Ping:
		doc["config"] = t.PingConfig
	}

	return doc
}
