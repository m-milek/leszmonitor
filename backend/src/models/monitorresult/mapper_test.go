package monitorresult

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/monitors"
	"testing"
)

func TestParseResultDetails(t *testing.T) {
	t.Run("HTTP details parsing", func(t *testing.T) {
		rawJSON := []byte(`{"statusCode": 200}`)
		details, err := ParseResultDetails("http", rawJSON)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		httpDetails, ok := details.(*HttpResultDetails)
		if !ok {
			t.Fatalf("expected HttpResultDetails, got %T", details)
		}

		if httpDetails.StatusCode != 200 {
			t.Errorf("expected 200, got %d", httpDetails.StatusCode)
		}
	})

	t.Run("Ping details parsing", func(t *testing.T) {
		rawJSON := []byte(`{"latencyMs": 42}`)
		details, err := ParseResultDetails("ping", rawJSON)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		pingDetails, ok := details.(*PingResultDetails)
		if !ok {
			t.Fatalf("expected PingResultDetails, got %T", details)
		}

		if pingDetails.LatencyMs != 42 {
			t.Errorf("expected 42, got %d", pingDetails.LatencyMs)
		}
	})
}

func TestMonitorResultJSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		original := newMonitorResult(
			uuid.New(),
			monitors.MonitorConfigType("http"),
			"success",
			false,
			100,
			"",
			&HttpResultDetails{StatusCode: 200},
			"2023-01-01T00:00:00Z",
		)

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("expected no error marshaling, got %v", err)
		}

		// check if the output contains the details
		if string(data) == "" {
			t.Fatalf("expected non-empty JSON")
		}
	})
}
