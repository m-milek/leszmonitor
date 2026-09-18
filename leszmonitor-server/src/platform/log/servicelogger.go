package log

import (
	"context"

	"github.com/rs/zerolog"
)

func MethodLoggerFromContext(ctx context.Context, serviceName, methodName string) zerolog.Logger {
	return FromContext(ctx).With().Str("service", serviceName).Str("fn", methodName).Logger()
}
