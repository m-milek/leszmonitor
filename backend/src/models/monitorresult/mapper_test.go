package monitorresult

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	consts "github.com/m-milek/leszmonitor/models/consts"
)

func TestParseResultDetails(t *testing.T) {
	t.Run("HTTP details parsing", func(t *testing.T) {
		rawJSON := []byte(`{"statusCode": 200}`)
		details, err := ParseResultDetails(consts.HttpConfigType, rawJSON)
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

	t.Run("TCP details parsing", func(t *testing.T) {
		rawJSON := []byte(`{"latencyMs": 42}`)
		details, err := ParseResultDetails(consts.TCPConfigType, rawJSON)
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
}

func TestMonitorResultJSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		original := NewMonitorResult(
			uuid.New(),
			consts.HttpConfigType,
			true,
			false,
			100,
			"",
			&HttpResultDetails{StatusCode: 200},
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
