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

type Config struct {
	Level       zerolog.Level
	LogFilePath string
}

// Global variables
var (
	// Configuration with defaults
	defaultConfig = Config{
		Level:       zerolog.InfoLevel,
		LogFilePath: "", // Empty means stdout only
	}

	// Mutex for thread-safe operations
	mu sync.Mutex

	Init   = getInitLogger()
	Main   zerolog.Logger
	Api    zerolog.Logger
	Logs   zerolog.Logger
	Uptime zerolog.Logger
	Db     zerolog.Logger
)

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
		resetCode)
}

func getInitLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatTimestamp = func(i interface{}) string {
		timestamp := fmt.Sprintf("%s", i)
		return fmt.Sprintf("%s [%s]", timestamp, "init")
	}
	output.FormatLevel = formatLogLevel

	return zerolog.New(output).With().Timestamp().Logger()
}

func createComponentLogger(component string, consoleWriter zerolog.ConsoleWriter, fileWriter io.Writer, level zerolog.Level) zerolog.Logger {
	// Configure console writer with component-specific timestamp
	consoleWriter.FormatTimestamp = func(i interface{}) string {
		return fmt.Sprintf("%s [%s]", i, component)
	}

	// Use MultiLevelWriter if both console and file writers are provided
	var writer io.Writer
	if fileWriter != nil {
		writer = zerolog.MultiLevelWriter(fileWriter, consoleWriter)
	} else {
		writer = consoleWriter
	}

	// Create and return the logger
	return zerolog.New(writer).Level(level).With().Timestamp().Logger()
}

// setupLoggers initializes all component loggers
func setupLoggers(consoleWriter zerolog.ConsoleWriter, fileWriter io.Writer, level zerolog.Level) {
	Main = createComponentLogger("main", consoleWriter, fileWriter, level)
	Api = createComponentLogger("http", consoleWriter, fileWriter, level)
	Logs = createComponentLogger("logc", consoleWriter, fileWriter, level)
	Uptime = createComponentLogger("uptm", consoleWriter, fileWriter, level)
	Db = createComponentLogger("mgdb", consoleWriter, fileWriter, level)
}

// InitLogging sets up logging with the given configuration
func InitLogging(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	consoleWriter.FormatLevel = formatLogLevel

	if os.Getenv(env.LogFilePath) == "" {
		Init.Warn().Msg("No log file path specified, defaulting to stdout")
	}

	shouldLogToFile := cfg.LogFilePath != ""

	if shouldLogToFile {
		Init.Debug().Msg("Log file path: " + cfg.LogFilePath)
	} else {
		Init.Debug().Msg("Logging to stdout only")
		setupLoggers(consoleWriter, nil, cfg.Level)
		return nil
	}

	// Ensure the directory exists
	dir := filepath.Dir(cfg.LogFilePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			Init.Error().Msg("Failed to create log directory")
			return err
		}
	}

	// Open the log file
	file, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Init.Error().Msg("Failed to open log file")
		return err
	}

	Init.Debug().Msg("Initialized logging to file: " + file.Name())
	setupLoggers(consoleWriter, file, cfg.Level)

	return nil
}

func GetLoggerConfig() Config {
	envLevel := os.Getenv(env.LogLevel)
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
		Init.Warn().Msg("Invalid log level, defaulting to 'INFO'")
	}

	envFilePath := os.Getenv(env.LogFilePath)
	var filePath string

	if envFilePath != "" {
		Init.Debug().Msg("Log file path from environment variable: " + envFilePath)
		filePath = envFilePath
	} else {
		Init.Warn().Msg("No log file path specified, using default")
	}

	return Config{
		Level:       logLevel,
		LogFilePath: filePath,
	}
}
