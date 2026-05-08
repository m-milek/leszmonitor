package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	appconfig "github.com/m-milek/leszmonitor/appconfig"
	"github.com/rs/zerolog"
)

type Config struct {
	Level       zerolog.Level
	PrettyPrint bool
}

func New() zerolog.Logger {
	level, err := zerolog.ParseLevel(appconfig.LogLevel)
	if err != nil {
		level = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(level)

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = formatLogLevel

	return zerolog.New(output).Level(level).With().Timestamp().Caller().Logger()
}

func NewServiceLogger(serviceName string) zerolog.Logger {
	return New().With().Str("service", serviceName).Logger()
}

func FromContext(ctx context.Context) *zerolog.Logger {
	logger := zerolog.Ctx(ctx)
	if logger == nil {
		fallbackLogger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		fallbackLogger.Warn().Msg("No logger found in context, using fallback logger")
		return &fallbackLogger
	}
	return logger
}

func WithContext(ctx context.Context, logger *zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}

func formatLogLevel(i interface{}) string {
	level := strings.ToUpper(fmt.Sprintf("%s", i))

	if i == nil {
		level = "LOG"
	}

	// Define color codes
	var colorCode string
	switch level {
	case "TRACE":
		colorCode = "\033[90m" // bright black/gray
	case "DEBUG":
		colorCode = "\033[36m" // cyan
	case "INFO":
		colorCode = "\033[32m" // green
	case "WARN":
		colorCode = "\033[33m" // yellow
	case "ERROR":
		colorCode = "\033[31m" // red
	case "FATAL":
		colorCode = "\033[35m" // magenta
	case "PANIC":
		colorCode = "\033[41;37m" // white on red
	default:
		colorCode = "\033[0m" // default
	}

	boldCode := "\033[1m"
	resetCode := "\033[0m"

	return fmt.Sprintf("%s%s%-5s%s",
		boldCode,
		colorCode,
		level,
		resetCode,
	)
}
