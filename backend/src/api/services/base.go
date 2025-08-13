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
	if s.methodLoggers == nil {
		s.methodLoggers = make(map[string]zerolog.Logger)
	}
	if logger, exists := s.methodLoggers[methodName]; exists {
		return logger
	}
	logger := s.serviceLogger.With().Str("method", methodName).Logger()
	s.methodLoggers[methodName] = logger
	return logger
}

func InitializeServices() {
	TeamService.serviceLogger = logging.NewServiceLogger("TeamService")
	UserService.serviceLogger = logging.NewServiceLogger("UserService")
	MonitorService.serviceLogger = logging.NewServiceLogger("MonitorService")
}
