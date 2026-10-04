package results

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
)

func TestParseResultDetails(t *testing.T) {
	t.Run("HTTP details parsing", func(t *testing.T) {
		rawJSON := []byte(`{"statusCode": 200}`)
		details, err := ParseResultDetails(kind.HTTPConfigType, rawJSON)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		httpDetails, ok := details.(*HTTPResultDetails)
		if !ok {
			t.Fatalf("expected HTTPResultDetails, got %T", details)
		}

		if httpDetails.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", httpDetails.StatusCode)
		}
	})

	t.Run("TCP details parsing", func(t *testing.T) {
		rawJSON := []byte(`{"latencyMs": 42}`)
		details, err := ParseResultDetails(kind.TCPConfigType, rawJSON)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		tcpDetails, ok := details.(*TCPResultDetails)
		if !ok {
			t.Fatalf("expected TCPResultDetails, got %T", details)
		}

		if tcpDetails.LatencyMs != 42 {
			t.Errorf("expected 42, got %d", tcpDetails.LatencyMs)
		}
	})

	t.Run("Push details parsing", func(t *testing.T) {
		rawJSON := []byte(`{"rawMessage": "backup done"}`)
		details, err := ParseResultDetails(kind.PushConfigType, rawJSON)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		pushDetails, ok := details.(*PushResultDetails)
		if !ok {
			t.Fatalf("expected PushResultDetails, got %T", details)
		}

		if pushDetails.RawMessage != "backup done" {
			t.Errorf("expected 'backup done', got %q", pushDetails.RawMessage)
		}
	})
}

func TestMonitorResultJSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		original := NewMonitorResult(
			uuid.New(),
			kind.HTTPConfigType,
			kind.MonitorStatusUp,
			false,
			new(int64(100)),
			&HTTPResultDetails{StatusCode: 200},
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
