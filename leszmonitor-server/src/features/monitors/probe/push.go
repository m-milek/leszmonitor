package probe

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/platform/log"
)

type PushProbe struct {
	GracePeriodSeconds int `json:"gracePeriodSeconds"`
}

type PushProbePayload struct {
	Status  kind.MonitorStatus
	Latency *int64
	Body    string
}

func (p *PushProbe) Validate() error {
	if p.GracePeriodSeconds < 0 {
		return errors.New("grace period must be 0s or more")
	}

	return nil
}

func (p *PushProbe) Run(ctx context.Context, monitorID uuid.UUID) (results.IMonitorResult, error) {
	logger := log.FromContext(ctx)

	err := errors.New("push probe cannot be ran like other probes")
	logger.Error().Err(err).Msg("Invalid probe run call")

	return nil, err

}
