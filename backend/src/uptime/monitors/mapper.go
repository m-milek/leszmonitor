package monitors

import (
	"fmt"
	"github.com/m-milek/leszmonitor/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"reflect"
)

var monitorTypeMap = map[string]reflect.Type{
	string(Http): reflect.TypeOf(HttpMonitor{}),
	string(Ping): reflect.TypeOf(PingMonitor{}),
}

func MapMonitorType(typeTag string) reflect.Type {
	monitorType := MonitorType(typeTag)
	if monitorType == "" {
		return nil
	}
	if monitorType, ok := monitorTypeMap[string(monitorType)]; ok {
		return monitorType
	}
	return nil
}

func MapFromBson(rawDoc bson.M) (IMonitor, error) {
	rawMonitorType, ok := rawDoc["type"]
	if !ok {
		return nil, fmt.Errorf("missing 'type' field in monitor document")
	}

	monitorType, ok := rawMonitorType.(string)
	if !ok {
		return nil, fmt.Errorf("invalid 'type' field in monitor document: %v", rawMonitorType)
	}

	monitorTypeReflect := MapMonitorType(monitorType)
	if monitorTypeReflect == nil {
		return nil, fmt.Errorf("unknown monitor type: %s", monitorType)
	}

	monitorInstance := reflect.New(monitorTypeReflect).Interface()

	// Check if the instance implements IMonitor
	monitor, ok := monitorInstance.(IMonitor)
	if !ok {
		logger.Uptime.Trace().Any("monitor", monitorInstance).Msg("Monitor type does not implement IMonitor interface")
		return nil, fmt.Errorf("monitor type %s does not implement IMonitor interface", monitorType)
	}

	// Now marshal and unmarshal to populate the monitor with data
	data, err := bson.Marshal(rawDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal monitor data: %v", err)
	}

	if err := bson.Unmarshal(data, monitor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal monitor data: %v", err)
	}

	return monitor, nil
}
