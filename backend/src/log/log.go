package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/logdyhq/logdy-core/logdy"
	appconfig "github.com/m-milek/leszmonitor/appconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var LogdyGlobalWriter *logdyWriter

type logdyWriter struct {
	Logger logdy.Logdy
}

func (w *logdyWriter) Write(p []byte) (int, error) {
	if w.Logger != nil {
		var fields logdy.Fields
		if err := json.Unmarshal(p, &fields); err == nil {
			w.Logger.Log(fields)
		} else {
			w.Logger.LogString(string(p))
		}
	}
	return len(p), nil
}

func init() {
	LogdyGlobalWriter = &logdyWriter{}
}

type Config struct {
	Level       zerolog.Level
	PrettyPrint bool
}

func New() zerolog.Logger {
	level, err := zerolog.ParseLevel(strings.ToLower(os.Getenv(appconfig.LogLevel)))
	if err != nil {
		level = zerolog.TraceLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var baseOutput io.Writer
	if strings.ToLower(os.Getenv(appconfig.LogFormat)) == "json" {
		baseOutput = os.Stdout
	} else {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = formatLogLevel
		baseOutput = output
	}

	multi := zerolog.MultiLevelWriter(baseOutput, LogdyGlobalWriter)
	return zerolog.New(multi).Level(level).With().Timestamp().Caller().Logger()
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
