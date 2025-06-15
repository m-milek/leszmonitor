package monitors

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"io"
)

var monitorTypeMap = map[MonitorConfigType]func() IMonitor{
	Http: func() IMonitor {
		return &HttpMonitor{}
	},
	Ping: func() IMonitor {
		return &PingMonitor{}
	},
}

func MapMonitorType(typeTag MonitorConfigType) IMonitor {
	if typeTag == "" {
		return nil
	}
	if monitorInstanceCreatorFn, ok := monitorTypeMap[typeTag]; ok {
		return monitorInstanceCreatorFn()
	}
	return nil
}

func FromReader(reader io.Reader) (IMonitor, error) {
	var rawData json.RawMessage
	if err := json.NewDecoder(reader).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %w", err)
	}

	var monitorTypeExtractor MonitorTypeExtractor
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

func FromRawBsonDoc(rawDoc bson.M) (IMonitor, error) {
	// Extract the monitor type
	monitorType, ok := rawDoc["type"].(string)
	if !ok {
		return nil, fmt.Errorf("monitor document missing 'type' field or not a string")
	}
	// Create the appropriate monitor instance
	monitor := MapMonitorType(MonitorConfigType(monitorType))

	// Convert the document to BSON and unmarshal it into the monitor
	data, err := bson.Marshal(rawDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal document: %w", err)
	}
	if err := bson.Unmarshal(data, monitor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal document into monitor: %w", err)
	}

	return monitor, nil
}
