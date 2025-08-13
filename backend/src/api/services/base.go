package services

import (
	"github.com/m-milek/leszmonitor/logging"
	"github.com/rs/zerolog"
)

type BaseService struct {
	serviceLogger zerolog.Logger
	methodLoggers map[string]zerolog.Logger
}

func (s *BaseService) getMethodLogger(methodName string) zerolog.Logger {
	if logger, exists := s.methodLoggers[methodName]; exists {
		return logger
	}
	return s.serviceLogger.With().Str("method", methodName).Logger()
}

func InitializeServices() {
	TeamService.serviceLogger = logging.NewServiceLogger("TeamService")
	UserService.serviceLogger = logging.NewServiceLogger("UserService")
	MonitorService.serviceLogger = logging.NewServiceLogger("MonitorService")
}
