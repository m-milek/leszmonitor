// Package logger provides a customized logging setup using zerolog.
//
// This package initializes several pre-configured loggers for different
// components of the application. Each logger uses a customized console output
// with colored level indicators and component name prefixes.
//
// Usage:
//
//	logger.Main.Info().Msg("Application started")
//	logger.Api.Error().Err(err).Msg("Failed to process request")
//	logger.Uptime.Debug().Int("status", status).Msg("Health check completed")
package logger

import (
	"fmt"
	"github.com/m-milek/leszmonitor/env"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type LogDestination int

const (
	Stdout LogDestination = iota
	File
	Both
)

type Config struct {
	Level       zerolog.Level
	Destination LogDestination
	FilePath    string
}

// Global variables
var (
	// Configuration with defaults
	defaultConfig = Config{
		Level:       zerolog.InfoLevel,
		FilePath:    "", // Empty means stdout only
		Destination: Stdout,
	}

	// The shared writer for all loggers
	logWriter io.Writer = os.Stdout

	// Mutex for thread-safe operations
	mu sync.Mutex

	tmpLogger = getTemporaryLogger()

	Main   zerolog.Logger
	Api    zerolog.Logger
	Logs   zerolog.Logger
	Uptime zerolog.Logger
)

// setupLoggers configures all loggers with the given writer
func setupLoggers(writer io.Writer) {
	Main = createLogger("main", writer)
	Api = createLogger("http", writer)
	Logs = createLogger("logs", writer)
	Uptime = createLogger("uptm", writer)
}

func getTemporaryLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		level := strings.ToUpper(fmt.Sprintf("%s", i))
		if i == nil {
			level = "LOG"
		}
		return fmt.Sprintf("[%s] %-5s", "init", level)
	}
	return zerolog.New(output).With().Timestamp().Logger()
}

// createLogger creates a new zerolog logger with the given component name
func createLogger(component string, writer io.Writer) zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: writer, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		level := strings.ToUpper(fmt.Sprintf("%s", i))

		if i == nil {
			level = "LOG"
		}

		if output.NoColor {
			return fmt.Sprintf("[%s] %-5s", component, level)
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

		return fmt.Sprintf("[%s] %s%s%-5s%s",
			component,
			boldCode,
			colorCode,
			level,
			resetCode)
	}

	return zerolog.New(output).
		Level(defaultConfig.Level).
		With().
		Timestamp().
		Logger()
}

// threadSafeWriter provides synchronized access to the underlying writer
type threadSafeWriter struct {
	mu     sync.Mutex
	writer io.Writer
}

// Write implements io.Writer with thread safety
func (t *threadSafeWriter) Write(p []byte) (n int, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.writer.Write(p)
}

// InitLogging sets up logging with the given configuration
func InitLogging(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	// Update the configuration
	defaultConfig = cfg

	// If no file path is specified, just use stdout
	if cfg.FilePath == "" {
		tmpLogger.Warn().Msg("No log file path specified, using stdout only")
		logWriter = os.Stdout
		setupLoggers(logWriter)
		return nil
	}

	// Ensure the directory exists
	dir := filepath.Dir(cfg.FilePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			tmpLogger.Error().Msg("Failed to create log directory")
			return err
		}
	}

	// Create the appropriate writer
	var writer io.Writer

	if cfg.Destination == Stdout {
		writer = os.Stdout
		tmpLogger.Debug().Msg("Logging to stdout only")
	} else {
		file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			tmpLogger.Error().Msg("Failed to open log file")
			return err
		}

		if cfg.Destination == Both {
			writer = io.MultiWriter(file, os.Stdout)
			tmpLogger.Debug().Msg("Logging to both stdout and file")
		} else if cfg.Destination == File {
			writer = file
			tmpLogger.Debug().Msg("Logging to file only")
		}
	}

	// Wrap with thread-safe writer
	logWriter = &threadSafeWriter{writer: writer}

	// Set up the loggers with the new writer
	setupLoggers(logWriter)

	return nil
}

func GetLoggerConfig() Config {
	envLevel := os.Getenv(env.LOG_LEVEL)
	var logLevel zerolog.Level

	switch envLevel {
	case "TRACE":
		logLevel = zerolog.TraceLevel
	case "DEBUG":
		logLevel = zerolog.DebugLevel
	case "INFO":
		logLevel = zerolog.InfoLevel
	case "WARN":
		logLevel = zerolog.WarnLevel
	case "ERROR":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
		tmpLogger.Warn().Msg("Invalid log level, defaulting to 'INFO'")
	}

	// Check if the environment variable is set

	return Config{
		Level:    logLevel,
		FilePath: defaultConfig.FilePath,
	}
}
