package probe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
)

// Probe runs a single check for a monitor and returns its result.
type Probe interface {
	Run(ctx context.Context, monitorID uuid.UUID) (results.IMonitorResult, error)
	Validate() error
}

var ErrInvalidProbeConfig = errors.New("invalid probe config")

func mapProbeType(probeType kind.ProbeType) Probe {
	switch probeType {
	case kind.HTTPConfigType:
		return &HTTPProbe{}
	case kind.TCPConfigType:
		return &TCPProbe{}
	case kind.DNSConfigType:
		return &DNSProbe{}
	case kind.PushConfigType:
		return &PushProbe{}
	default:
		return nil
	}
}

func Parse[T Probe](probeType kind.ProbeType, config string) (T, error) {
	var zero T
	decoder := json.NewDecoder(strings.NewReader(config))
	decoder.DisallowUnknownFields()

	p := mapProbeType(probeType)
	if p == nil {
		return zero, fmt.Errorf("unknown probe type: %s", probeType)
	}

	if err := decoder.Decode(p); err != nil {
		return zero, fmt.Errorf("failed to parse probe config: %w", err)
	}

	typed, ok := p.(T)
	if !ok {
		return zero, fmt.Errorf("probe type %s is %T, expected %T", probeType, p, zero)
	}

	if err := typed.Validate(); err != nil {
		return zero, fmt.Errorf("%w: %w", ErrInvalidProbeConfig, err)
	}
	return typed, nil
}
