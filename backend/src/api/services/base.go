package services

import (
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/rs/zerolog"
)

type baseService struct {
	authService   IAuthorizationService
	serviceLogger zerolog.Logger
	methodLoggers map[string]zerolog.Logger
}

func newBaseService(authService IAuthorizationService, serviceName string) baseService {
	return baseService{
		authService:   authService,
		serviceLogger: logging.NewServiceLogger(serviceName),
		methodLoggers: make(map[string]zerolog.Logger),
	}
}

// getDB retrieves the database singleton safely
func (s *baseService) getDB() db.DB {
	return db.Get()
}

// Return a logger for a specific service method, creating it if it doesn't exist yet.
func (s *baseService) getMethodLogger(methodName string) zerolog.Logger {
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
