package probe

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/results"
)

type PingProbe struct {
	Host      string `json:"host"`
	Count     int    `json:"count"`
	TimeoutMs int    `json:"timeoutMs"`
}

func (p *PingProbe) Run(ctx context.Context, monitorID uuid.UUID) (results.IMonitorResult, error) {
	return nil, errors.New("ping probe is not implemented")
}

func (p *PingProbe) Validate() error {
	return nil
}
