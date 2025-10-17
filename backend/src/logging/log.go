// Package logging provides a customized logging setup using zerolog.
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
//
//	// Create custom service logger
//	myLogger := logger.NewServiceLogger("myservice")
//	myLogger.Info().Msg("Custom service started")
package logging

import (
	"fmt"
	"github.com/m-milek/leszmonitor/env"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	// Internal root logger
	rootLogger zerolog.Logger

	// Pre-configured loggers
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
		resetCode,
	)
}

func getInitLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = formatLogLevel

	return zerolog.New(output).With().Timestamp().Logger()
}

var consoleWriter zerolog.ConsoleWriter
var fileWriter io.Writer
var currentLevel zerolog.Level

// NewServiceLogger creates a new logger with a custom service name.
// The service name will be added as a persistent field to all log entries.
func NewServiceLogger(serviceName string) zerolog.Logger {
	if rootLogger.GetLevel() == zerolog.Disabled {
		// If root logger isn't initialized yet, create a basic logger
		return createLogger(serviceName, currentLevel)
	}

	// Create a child logger from the root logger with the service field
	return rootLogger.With().Str("service", serviceName).Caller().Logger()
}

func createLogger(component string, level zerolog.Level) zerolog.Logger {
	// Use MultiLevelWriter if both console and file writers are provided
	var writer io.Writer
	if fileWriter != nil {
		Init.Debug().Msg("Logging to file: " + fileWriter.(*os.File).Name())
		writer = zerolog.MultiLevelWriter(fileWriter, consoleWriter)
	} else {
		Init.Debug().Msg("Logging to console only")
		writer = consoleWriter
	}

	logger := zerolog.New(writer).Level(level).With().Timestamp()

	// Create and return the logger with service field
	return logger.Str("service", component).Logger()
}

// setupLoggers initializes all component loggers and the root logger.
func setupLoggers(level zerolog.Level) {
	// Initialize the root logger
	var writer io.Writer
	if fileWriter != nil {
		writer = zerolog.MultiLevelWriter(fileWriter, consoleWriter)
	} else {
		writer = consoleWriter
	}

	rootLogger = zerolog.New(writer).Level(level).With().Timestamp().Logger()
	currentLevel = level

	// Create pre-configured loggers using the new method
	Main = NewServiceLogger("Main")
	Api = NewServiceLogger("API")
	Logs = NewServiceLogger("LogCollector")
	Uptime = NewServiceLogger("UptimeMonitor")
	Db = NewServiceLogger("MongoDB")
}

// InitLogging sets up logging with the given configuration.
func InitLogging(cfg Config) error {
	consoleWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	consoleWriter.FormatLevel = formatLogLevel

	if os.Getenv(env.LogFilePath) == "" {
		Init.Warn().Msg("No log file path specified, defaulting to stdout")
	}

	shouldLogToFile := cfg.LogFilePath != ""

	if !shouldLogToFile {
		Init.Debug().Msg("Logging to stdout only")
		setupLoggers(cfg.Level)
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
	fileWriter = io.Writer(file)
	setupLoggers(cfg.Level)

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
		filePath = envFilePath
	} else {
		Init.Warn().Msg("No log file path specified, using default")
	}

	return Config{
		Level:       logLevel,
		LogFilePath: filePath,
	}
}
