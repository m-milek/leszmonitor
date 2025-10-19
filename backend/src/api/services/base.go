package services

import (
	"github.com/rs/zerolog"
)

type baseService struct {
	authService   *authorizationServiceT
	serviceLogger zerolog.Logger
	methodLoggers map[string]zerolog.Logger
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
