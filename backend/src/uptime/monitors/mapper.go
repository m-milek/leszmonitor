package monitors

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
