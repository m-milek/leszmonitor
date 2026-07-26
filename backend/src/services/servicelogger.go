package services

import (
	"context"

	"github.com/m-milek/leszmonitor/log"
	"github.com/rs/zerolog"
)

func MethodLoggerFromContext(ctx context.Context, serviceName, methodName string) zerolog.Logger {
	return log.FromContext(ctx).With().Str("service", serviceName).Str("fn", methodName).Logger()
}
