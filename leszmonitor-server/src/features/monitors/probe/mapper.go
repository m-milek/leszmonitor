package probe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
)

// Probe runs a single check for a monitor and returns its result.
type Probe interface {
	Run(ctx context.Context, monitorID uuid.UUID) (results.IMonitorResult, error)
	Validate() error
}

func mapProbeType(probeType kind.ProbeType) Probe {
	switch probeType {
	case kind.HTTPConfigType:
		return &HTTPProbe{}
	case kind.TCPConfigType:
		return &TCPProbe{}
	case kind.DNSConfigType:
		return &DNSProbe{}
	default:
		return nil
	}
}

func ProbeFromJSON(probeConfig string, probeType kind.ProbeType) (Probe, error) {
	// Map the monitor type to the appropriate config type
	probe := mapProbeType(probeType)
	if probe == nil {
		return nil, fmt.Errorf("unknown monitor type: %s", probeType)
	}

	// unmarshal the raw data into a probe instance
	if err := json.Unmarshal([]byte(probeConfig), &probe); err != nil {
		return nil, fmt.Errorf("failed to parse monitor config: %w: %s", err, probeConfig)
	}

	return probe, nil
}

func UnmarshalProbeFromBytes(probeType kind.ProbeType, data []byte) (Probe, error) {
	config := mapProbeType(probeType)
	if config == nil {
		return nil, fmt.Errorf("unknown monitor config type: %s", probeType)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal monitor config: %w", err)
	}
	return config, nil
}
