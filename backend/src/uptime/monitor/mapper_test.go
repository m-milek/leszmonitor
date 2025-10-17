package monitors

import (
	"bytes"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestMapMonitorType(t *testing.T) {
	type args struct {
		typeTag MonitorConfigType
	}
	tests := []struct {
		name string
		args args
		want IMonitor
	}{
		{
			name: "Valid Http Monitor Type",
			args: args{typeTag: Http},
			want: &HttpMonitor{},
		},
		{
			name: "Valid Ping Monitor Type",
			args: args{typeTag: Ping},
			want: &PingMonitor{},
		},
		{
			name: "Invalid Monitor Type",
			args: args{typeTag: "InvalidType"},
			want: nil,
		},
		{
			name: "Empty Monitor Type",
			args: args{typeTag: ""},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, MapMonitorType(tt.args.typeTag), "MapMonitorType(%v)", tt.args.typeTag)
		})
	}
}

func TestFromReader(t *testing.T) {
	t.Run("Successfully parse HTTP monitor", func(t *testing.T) {
		// Create a test HTTP monitor
		testMonitor := NewTestMonitor().AsHttp()
		testMonitor.Name = "Test HTTP Monitor"
		testMonitor.ID = pgtype.UUID{}
		testMonitor.Interval = 30
		testMonitor.HttpConfig.Url = "https://example.com/api"
		testMonitor.HttpConfig.ExpectedStatusCodes = []int{200, 201}

		// Convert to JSON
		jsonData, err := json.Marshal(testMonitor.Build())
		assert.NoError(t, err)

		// Create a reader
		reader := bytes.NewReader(jsonData)

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, monitor)
		assert.Equal(t, pgtype.UUID{}, monitor.GetId())
		assert.Equal(t, "Test HTTP Monitor", monitor.GetName())
		assert.Equal(t, Http, monitor.GetType())

		// Type assertion to check specific config
		httpMonitor, ok := monitor.(*HttpMonitor)
		assert.True(t, ok)
		assert.Equal(t, "https://example.com/api", httpMonitor.Config.Url)
		assert.Equal(t, []int{200, 201}, httpMonitor.Config.ExpectedStatusCodes)
	})

	t.Run("Successfully parse Ping monitor", func(t *testing.T) {
		// Create a test Ping monitor
		testMonitor := NewTestMonitor().AsPing()
		testMonitor.Name = "Test Ping Monitor"
		testMonitor.Description = "This is a test ping monitor"
		testMonitor.ID = pgtype.UUID{}
		testMonitor.Interval = 10
		testMonitor.PingConfig.Host = "ping.example.com"
		testMonitor.PingConfig.Port = "443"
		testMonitor.PingConfig.Protocol = "tcp"

		// Convert to JSON
		jsonData, err := json.Marshal(testMonitor.Build())
		assert.NoError(t, err)

		// Create a reader
		reader := bytes.NewReader(jsonData)

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, monitor)
		assert.Equal(t, "Test Ping Monitor", monitor.GetName())
		assert.Equal(t, "This is a test ping monitor", monitor.GetDescription())
		assert.Equal(t, Ping, monitor.GetType())

		// Type assertion to check specific config
		pingMonitor, ok := monitor.(*PingMonitor)
		assert.True(t, ok)
		assert.Equal(t, "ping.example.com", pingMonitor.Config.Host)
		assert.Equal(t, "443", pingMonitor.Config.Port)
		assert.Equal(t, "tcp", pingMonitor.Config.Protocol)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		// Create a reader with invalid JSON
		reader := strings.NewReader("this is not valid JSON")

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "failed to decode request body")
	})

	t.Run("Missing type field", func(t *testing.T) {
		// Create JSON without a type field
		jsonData := `{
		"id": "test123",
		"name": "Test Monitor",
		"description": "Test description",
		"interval": 30,
		"ownerId": "test_owner"
	}`

		// Create a reader
		reader := strings.NewReader(jsonData)

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "monitor type cannot be empty")
	})

	t.Run("Unknown monitor type", func(t *testing.T) {
		// Create JSON with an unknown monitor type
		jsonData := `{
		"id": "test123",
		"name": "Test Monitor",
		"description": "Test description",
		"interval": 30,
		"ownerId": "test_owner",
		"type": "unknown_type"
	}`

		// Create a reader
		reader := strings.NewReader(jsonData)

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "unknown monitor type")
	})

	t.Run("Invalid monitor configuration", func(t *testing.T) {
		// Create JSON with invalid configuration
		jsonData := `{
		"id": "test123",
		"name": "Test HTTP Monitor",
		"description": "Test description",
		"interval": 30,
		"ownerId": "test_owner",
		"type": "http",
		"config": {
			"url": 123,
			"httpMethod": "GET"
		}
	}`

		// Create a reader
		reader := strings.NewReader(jsonData)

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "failed to parse monitor config")
	})

	t.Run("Empty reader", func(t *testing.T) {
		// Create an empty reader
		reader := strings.NewReader("")

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "failed to decode request body")
	})

	t.Run("Reader error", func(t *testing.T) {
		// Create a reader that returns an error
		reader := &ErrorReader{Err: io.ErrUnexpectedEOF}

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "failed to decode request body")
	})

	t.Run("Partial JSON", func(t *testing.T) {
		// Create a reader with partial JSON
		reader := strings.NewReader(`{"type": "http", "name": "Test Monitor"`)

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "failed to decode request body")
	})
}

// ErrorReader is a mock reader that always returns an error
type ErrorReader struct {
	Err error
}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, r.Err
}
