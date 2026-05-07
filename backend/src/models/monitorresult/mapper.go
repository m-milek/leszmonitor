package monitorresult

import (
	"encoding/json"
	"fmt"

	consts "github.com/m-milek/leszmonitor/models/consts"
)

// ParseResultDetails parses the raw JSON details based on the monitorType
func ParseResultDetails(monitorType consts.ProbeType, rawDetails []byte) (IMonitorResultDetails, error) {
	if len(rawDetails) == 0 || string(rawDetails) == "null" {
		return nil, nil
	}

	switch monitorType {
	case consts.HttpConfigType:
		var details HttpResultDetails
		if err := json.Unmarshal(rawDetails, &details); err != nil {
			return nil, fmt.Errorf("failed to parse HTTP result details: %w", err)
		}
		return &details, nil
	case consts.TCPConfigType:
		var details TCPResultDetails
		if err := json.Unmarshal(rawDetails, &details); err != nil {
			return nil, fmt.Errorf("failed to parse TCP result details: %w", err)
		}
		return &details, nil
	default:
		return nil, fmt.Errorf("unknown monitor type for result details: %s", monitorType)
	}
}
