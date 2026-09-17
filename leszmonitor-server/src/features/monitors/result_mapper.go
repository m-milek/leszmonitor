package monitors

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrEmptyDetails = errors.New("empty details")

// ParseResultDetails parses the raw JSON details based on the monitorType.
func ParseResultDetails(monitorType ProbeType, rawDetails []byte) (IMonitorResultDetails, error) {
	if len(rawDetails) == 0 || string(rawDetails) == "null" {
		return nil, ErrEmptyDetails
	}

	switch monitorType {
	case HTTPConfigType:
		var details HTTPResultDetails
		if err := json.Unmarshal(rawDetails, &details); err != nil {
			return nil, fmt.Errorf("failed to parse HTTP result details: %w", err)
		}
		return &details, nil
	case TCPConfigType:
		var details TCPResultDetails
		if err := json.Unmarshal(rawDetails, &details); err != nil {
			return nil, fmt.Errorf("failed to parse TCP result details: %w", err)
		}
		return &details, nil
	case DNSConfigType:
		var details DNSResultDetails
		if err := json.Unmarshal(rawDetails, &details); err != nil {
			return nil, fmt.Errorf("failed to parse DNS result details: %w", err)
		}
		return &details, nil
	default:
		return nil, fmt.Errorf("unknown monitor type for result details: %s", monitorType)
	}
}
