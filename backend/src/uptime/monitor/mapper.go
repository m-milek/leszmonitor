package monitors

import (
	"encoding/json"
	"fmt"
	"io"
)

var monitorTypeMap = map[MonitorConfigType]func() IConcreteMonitor{
	httpType: func() IConcreteMonitor {
		return &httpMonitor{}
	},
	pingType: func() IConcreteMonitor {
		return &PingMonitor{}
	},
}

func MapMonitorType(typeTag MonitorConfigType) IConcreteMonitor {
	if typeTag == "" {
		return nil
	}
	if monitorInstanceCreatorFn, ok := monitorTypeMap[typeTag]; ok {
		return monitorInstanceCreatorFn()
	}
	return nil
}

func MapMonitorConfigType(kind MonitorConfigType) IMonitorConfig {
	switch kind {
	case httpType:
		return &httpConfig{}
	case pingType:
		return &PingConfig{}
	default:
		return nil
	}
}

func FromReader(reader io.Reader) (IConcreteMonitor, error) {
	var rawData json.RawMessage
	if err := json.NewDecoder(reader).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	var monitorTypeExtractor monitorTypeExtractor
	if err := json.Unmarshal(rawData, &monitorTypeExtractor); err != nil {
		return nil, fmt.Errorf("failed to parse monitor type: %w", err)
	}

	// Map the monitor type to the appropriate config type
	monitor := MapMonitorType(monitorTypeExtractor.Type)
	if monitorTypeExtractor.Type == "" {
		return nil, fmt.Errorf("monitor type cannot be empty")
	}
	if monitor == nil {
		return nil, fmt.Errorf("unknown monitor type: %s", monitorTypeExtractor.Type)
	}

	// unmarshal the raw data into a monitor instance
	if err := json.Unmarshal(rawData, &monitor); err != nil {
		return nil, fmt.Errorf("failed to parse monitor config: %w", err)
	}

	return monitor, nil
}

func UnmarshalConfigFromBytes(kind MonitorConfigType, data []byte) (IMonitorConfig, error) {
	config := MapMonitorConfigType(kind)
	if config == nil {
		return nil, fmt.Errorf("unknown monitor config type: %s", kind)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal monitor config: %w", err)
	}
	return config, nil
}
