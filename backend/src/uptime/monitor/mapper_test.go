package monitors

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func TestFromRawBsonDoc(t *testing.T) {
	t.Run("Successfully parse HTTP monitor", func(t *testing.T) {
		// Create a test HTTP monitor
		testMonitor := NewTestMonitor().AsHttp()
		testMonitor.Name = "Test HTTP Monitor"
		testMonitor.ID = "test123"
		testMonitor.Interval = 30
		testMonitor.HttpConfig.Url = "https://example.com/api"
		testMonitor.HttpConfig.ExpectedStatusCodes = []int{200, 201}

		rawDoc := testMonitor.ToBSON()

		// Execute
		monitor, err := FromRawBsonDoc(rawDoc)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, monitor)
		assert.Equal(t, "test123", monitor.GetId())
		assert.Equal(t, "Test HTTP Monitor", monitor.GetName())
		assert.Equal(t, Http, monitor.GetType())
		assert.Equal(t, 30, monitor.(*HttpMonitor).Interval)

		// Type assertion to check specific config
		httpMonitor, ok := monitor.(*HttpMonitor)
		assert.True(t, ok)
		assert.Equal(t, "https://example.com/api", httpMonitor.Config.Url)
		assert.Equal(t, []int{200, 201}, httpMonitor.Config.ExpectedStatusCodes)
	})

	t.Run("Successfully parse Ping monitor", func(t *testing.T) {
		// Create a test Ping monitor
		testMonitor := NewTestMonitor().AsPing()
		testMonitor.Name = "Test Ping Monitor 19"
		testMonitor.Description = "This is a test ping monitor after updates"
		testMonitor.ID = "2C07IrLHg"
		testMonitor.Interval = 10
		testMonitor.PingConfig.Host = "gnu.org"
		testMonitor.PingConfig.Port = "443"
		testMonitor.PingConfig.Protocol = "tcp"
		testMonitor.PingConfig.PingTimeout = 5
		testMonitor.PingConfig.RetryCount = 3

		rawDoc := testMonitor.ToBSON()

		// Execute
		monitor, err := FromRawBsonDoc(rawDoc)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, monitor)
		assert.Equal(t, "2C07IrLHg", monitor.GetId())
		assert.Equal(t, "Test Ping Monitor 19", monitor.GetName())
		assert.Equal(t, "This is a test ping monitor after updates", monitor.GetDescription())
		assert.Equal(t, Ping, monitor.GetType())
		assert.Equal(t, 10, monitor.(*PingMonitor).Interval)

		// Type assertion to check specific config
		pingMonitor, ok := monitor.(*PingMonitor)
		assert.True(t, ok)
		assert.Equal(t, "gnu.org", pingMonitor.Config.Host)
		assert.Equal(t, "443", pingMonitor.Config.Port)
		assert.Equal(t, "tcp", pingMonitor.Config.Protocol)
		assert.Equal(t, 5, pingMonitor.Config.PingTimeout)
		assert.Equal(t, 3, pingMonitor.Config.RetryCount)
	})

	t.Run("Missing type field", func(t *testing.T) {
		// Create a raw document without a type field
		rawDoc := bson.M{
			"_id":         bson.NewObjectID(),
			"id":          "test123",
			"name":        "Test Monitor",
			"description": "Test description",
			"interval":    30,
			"ownerId":     "test_owner",
			// No type field
		}

		// Execute
		monitor, err := FromRawBsonDoc(rawDoc)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "missing 'type' field")
	})

	t.Run("Type field not a string", func(t *testing.T) {
		// Create a raw document with a non-string type field
		rawDoc := bson.M{
			"_id":         bson.NewObjectID(),
			"id":          "test123",
			"name":        "Test Monitor",
			"description": "Test description",
			"interval":    30,
			"ownerId":     "test_owner",
			"type":        123, // Not a string
		}

		// Execute
		monitor, err := FromRawBsonDoc(rawDoc)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "not a string")
	})

	t.Run("Unknown monitor type", func(t *testing.T) {
		// Create a raw document with an unknown monitor type
		rawDoc := bson.M{
			"_id":         bson.NewObjectID(),
			"id":          "test123",
			"name":        "Test Monitor",
			"description": "Test description",
			"interval":    30,
			"ownerId":     "test_owner",
			"type":        "unknown_type",
		}

		// Execute
		monitor, err := FromRawBsonDoc(rawDoc)

		// Assert
		// The behavior depends on MapMonitorType implementation
		// If it returns nil for unknown types, this should fail
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "failed to unmarshal")
	})

	t.Run("Unmarshal error with invalid config", func(t *testing.T) {
		// Create a monitor with invalid configuration
		testMonitor := NewTestMonitor().AsHttp()
		rawDoc := testMonitor.ToBSON()

		// Replace the config with a new map containing invalid types
		rawDoc["config"] = bson.M{
			"url":        123, // URL should be a string, not a number
			"httpMethod": "GET",
		}

		// Execute
		monitor, err := FromRawBsonDoc(rawDoc)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "failed to unmarshal")
	})

	t.Run("Unmarshal error with invalid structure", func(t *testing.T) {
		rawDoc := bson.M{
			"_id":  bson.NewObjectID(),
			"id":   "test123",
			"name": "Test Monitor",
			"type": "http",
			// Use a complex nested structure that doesn't match the expected monitor structure
			"config": bson.M{
				"url": bson.M{
					"protocol": 123,         // Should be a string
					"host":     []int{1, 2}, // Should be a string
				},
			},
		}

		// Execute
		monitor, err := FromRawBsonDoc(rawDoc)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "failed to unmarshal")
	})
}

// Test the entire flow from BSON to monitor and back
// Test the entire flow from BSON to monitor and back
func TestBsonToMonitorRoundTrip(t *testing.T) {
	t.Run("HTTP monitor round trip", func(t *testing.T) {
		// Create a test HTTP monitor
		testMonitor := NewTestMonitor().AsHttp()
		testMonitor.Name = "HTTP Round Trip"
		testMonitor.HttpConfig.Url = "https://example.org"
		testMonitor.HttpConfig.HttpMethod = "POST"
		testMonitor.HttpConfig.Headers = map[string]string{"Content-Type": "application/json"}

		// Convert to BSON
		rawDoc := testMonitor.ToBSON()

		// Convert back to monitor
		monitor, err := FromRawBsonDoc(rawDoc)
		assert.NoError(t, err)

		// Convert monitor back to BSON
		data, err := bson.Marshal(monitor)
		assert.NoError(t, err)

		var roundTripDoc bson.M
		err = bson.Unmarshal(data, &roundTripDoc)
		assert.NoError(t, err)

		// Compare important fields
		assert.Equal(t, testMonitor.ID, roundTripDoc["id"])
		assert.Equal(t, testMonitor.Name, roundTripDoc["name"])
		assert.Equal(t, string(testMonitor.Type), roundTripDoc["type"])

		// Check config - handle different possible types
		if config, ok := roundTripDoc["config"].(bson.M); ok {
			// If config is a bson.M, check fields directly
			assert.Equal(t, testMonitor.HttpConfig.Url, config["url"])
			assert.Equal(t, testMonitor.HttpConfig.HttpMethod, config["httpMethod"])
		} else {
			// If config is not a bson.M, marshal it to JSON and compare the values
			// This is a more flexible approach that works regardless of the exact type
			configBytes, err := bson.MarshalExtJSON(roundTripDoc["config"], true, false)
			assert.NoError(t, err)

			// Check if the JSON contains the expected values
			configStr := string(configBytes)
			assert.Contains(t, configStr, testMonitor.HttpConfig.Url)
			assert.Contains(t, configStr, testMonitor.HttpConfig.HttpMethod)
			assert.Contains(t, configStr, "application/json")
		}
	})

	t.Run("Ping monitor round trip", func(t *testing.T) {
		// Create a test Ping monitor
		testMonitor := NewTestMonitor().AsPing()
		testMonitor.Name = "Ping Round Trip"
		testMonitor.PingConfig.Host = "ping.example.com"
		testMonitor.PingConfig.Port = "8080"

		// Convert to BSON
		rawDoc := testMonitor.ToBSON()

		// Convert back to monitor
		monitor, err := FromRawBsonDoc(rawDoc)
		assert.NoError(t, err)

		// Convert monitor back to BSON
		data, err := bson.Marshal(monitor)
		assert.NoError(t, err)

		var roundTripDoc bson.M
		err = bson.Unmarshal(data, &roundTripDoc)
		assert.NoError(t, err)

		// Compare important fields
		assert.Equal(t, testMonitor.ID, roundTripDoc["id"])
		assert.Equal(t, testMonitor.Name, roundTripDoc["name"])
		assert.Equal(t, string(testMonitor.Type), roundTripDoc["type"])

		// Check config - handle different possible types
		if config, ok := roundTripDoc["config"].(bson.M); ok {
			// If config is a bson.M, check fields directly
			assert.Equal(t, testMonitor.PingConfig.Host, config["host"])
			assert.Equal(t, testMonitor.PingConfig.Port, config["port"])
		} else {
			// If config is not a bson.M, marshal it to JSON and compare the values
			configBytes, err := bson.MarshalExtJSON(roundTripDoc["config"], true, false)
			assert.NoError(t, err)

			// Check if the JSON contains the expected values
			configStr := string(configBytes)
			assert.Contains(t, configStr, testMonitor.PingConfig.Host)
			assert.Contains(t, configStr, testMonitor.PingConfig.Port)
		}
	})
}

// Test edge cases
func TestFromRawBsonDocEdgeCases(t *testing.T) {
	t.Run("Empty document", func(t *testing.T) {
		rawDoc := bson.M{}

		monitor, err := FromRawBsonDoc(rawDoc)

		assert.Error(t, err)
		assert.Nil(t, monitor)
		assert.Contains(t, err.Error(), "missing 'type' field")
	})
}

func TestFromReader(t *testing.T) {
	t.Run("Successfully parse HTTP monitor", func(t *testing.T) {
		// Create a test HTTP monitor
		testMonitor := NewTestMonitor().AsHttp()
		testMonitor.Name = "Test HTTP Monitor"
		testMonitor.ID = "test123"
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
		assert.Equal(t, "test123", monitor.GetId())
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
		testMonitor.ID = "ping123"
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
		assert.Equal(t, "ping123", monitor.GetId())
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

	t.Run("Valid JSON with extra fields", func(t *testing.T) {
		// Create JSON with extra fields
		jsonData := `{
		"id": "test123",
		"name": "Test HTTP Monitor",
		"description": "Test description",
		"interval": 30,
		"ownerId": "test_owner",
		"type": "http",
		"config": {
			"url": "https://example.com",
			"httpMethod": "GET"
		},
		"extraField1": "extra value",
		"extraField2": 123
	}`

		// Create a reader
		reader := strings.NewReader(jsonData)

		// Execute
		monitor, err := FromReader(reader)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, monitor)
		assert.Equal(t, "test123", monitor.GetId())
		assert.Equal(t, "Test HTTP Monitor", monitor.GetName())
		assert.Equal(t, Http, monitor.GetType())

		// Type assertion to check specific config
		httpMonitor, ok := monitor.(*HttpMonitor)
		assert.True(t, ok)
		assert.Equal(t, "https://example.com", httpMonitor.Config.Url)
		assert.Equal(t, "GET", httpMonitor.Config.HttpMethod)
	})
}

// ErrorReader is a mock reader that always returns an error
type ErrorReader struct {
	Err error
}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, r.Err
}
