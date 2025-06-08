package monitors

var monitorTypeMap = map[MonitorConfigType]IMonitor{
	Http: &HttpMonitor{},
	Ping: &PingMonitor{},
}

func MapMonitorType(typeTag MonitorConfigType) IMonitor {
	if typeTag == "" {
		return nil
	}
	if monitorInstance, ok := monitorTypeMap[typeTag]; ok {
		return monitorInstance
	}
	return nil
}
