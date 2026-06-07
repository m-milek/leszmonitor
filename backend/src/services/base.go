package services

import (
	"context"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/log"
	"github.com/rs/zerolog"
)

type baseService struct {
	authService     IAuthorizer
	auditLogService IAuditLogger
	serviceLogger   zerolog.Logger
	methodLoggers   map[string]zerolog.Logger
}

func newBaseService(authService IAuthorizer, auditLogService IAuditLogger, serviceName string) baseService {
	return baseService{
		authService:     authService,
		auditLogService: auditLogService,
		serviceLogger:   log.NewServiceLogger(serviceName),
		methodLoggers:   make(map[string]zerolog.Logger),
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

func MethodLoggerFromContext(ctx context.Context, serviceName, methodName string) zerolog.Logger {
	return log.FromContext(ctx).With().Str("service", serviceName).Str("method", methodName).Logger()
}
