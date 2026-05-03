package monitorresult

import (
	"encoding/json"
	"fmt"

	"github.com/m-milek/leszmonitor/models/monitors"
)

// ParseResultDetails parses the raw JSON details based on the monitorType
func ParseResultDetails(monitorType monitors.MonitorConfigType, rawDetails []byte) (IMonitorResultDetails, error) {
	if len(rawDetails) == 0 || string(rawDetails) == "null" {
		return nil, nil
	}

	switch monitorType {
	case monitors.HttpConfigType:
		var details HttpResultDetails
		if err := json.Unmarshal(rawDetails, &details); err != nil {
			return nil, fmt.Errorf("failed to parse HTTP result details: %w", err)
		}
		return &details, nil
	case monitors.PingConfigType:
		var details PingResultDetails
		if err := json.Unmarshal(rawDetails, &details); err != nil {
			return nil, fmt.Errorf("failed to parse Ping result details: %w", err)
		}
		return &details, nil
	default:
		return nil, fmt.Errorf("unknown monitor type for result details: %s", monitorType)
	}
}
