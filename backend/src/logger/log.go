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

	// Mutex for thread-safe operations
	mu sync.Mutex

	Init   = getTemporaryLogger()
	Main   zerolog.Logger
	Api    zerolog.Logger
	Logs   zerolog.Logger
	Uptime zerolog.Logger
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

func getTemporaryLogger() zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatTimestamp = func(i interface{}) string {
		timestamp := fmt.Sprintf("%s", i)
		return fmt.Sprintf("%s [%s]", timestamp, "init")
	}
	output.FormatLevel = formatLogLevel

	return zerolog.New(output).With().Timestamp().Logger()
}

// createComponentLogger creates a logger with the specified component name
func createComponentLogger(component string, consoleWriter *zerolog.ConsoleWriter, fileWriter io.Writer) zerolog.Logger {
	// Create a nop logger as a base
	logger := zerolog.Nop()

	// If we have a console writer, set up console logging
	if consoleWriter != nil {
		consoleCopy := *consoleWriter
		consoleCopy.FormatTimestamp = func(i interface{}) string {
			timestamp := fmt.Sprintf("%s", i)
			return fmt.Sprintf("%s [%s]", timestamp, component)
		}
		logger = zerolog.New(consoleCopy).Level(defaultConfig.Level).With().Timestamp().Logger()
	}

	// If we have a file writer, set up file logging with a hook
	if fileWriter != nil {
		// Create a direct file logger
		fileLogger := zerolog.New(fileWriter).Level(defaultConfig.Level).With().
			Timestamp().
			Str("component", component).
			Logger()

		// Add a hook to the main logger to also log to the file
		logger = logger.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
			if level >= defaultConfig.Level {
				// Create a new event with the same level
				fileEvent := fileLogger.WithLevel(level)

				// Log the message to the file
				fileEvent.Msg(message)
			}
		}))
	}

	// If no writers were provided, default to stdout
	if consoleWriter == nil && fileWriter == nil {
		logger = zerolog.New(os.Stdout).Level(defaultConfig.Level).With().Timestamp().Logger()
	}

	return logger
}

// setupLoggers initializes all component loggers
func setupLoggers(consoleWriter *zerolog.ConsoleWriter, fileWriter io.Writer) {
	Main = createComponentLogger("main", consoleWriter, fileWriter)
	Api = createComponentLogger("http", consoleWriter, fileWriter)
	Logs = createComponentLogger("logc", consoleWriter, fileWriter)
	Uptime = createComponentLogger("uptm", consoleWriter, fileWriter)
}

// InitLogging sets up logging with the given configuration
func InitLogging(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	// Update the configuration
	defaultConfig = cfg

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	consoleWriter.FormatLevel = formatLogLevel
	// The FormatTimestamp will be overridden for each component logger

	if cfg.FilePath == "" && cfg.Destination != Stdout {
		Init.Fatal().Msg("No log file path specified, defaulting to stdout")
		cfg.Destination = Stdout
	}

	var fileWriter io.Writer
	if cfg.Destination != Stdout {
		// Ensure the directory exists
		dir := filepath.Dir(cfg.FilePath)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				Init.Error().Msg("Failed to create log directory")
				return err
			}
		}

		file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			Init.Error().Msg("Failed to open log file")
			return err
		}
		Init.Debug().Msg("Logging to file: " + cfg.FilePath)
		fileWriter = file
	}

	switch cfg.Destination {
	case Stdout:
		Init.Debug().Msg("Logging to stdout only")
		setupLoggers(&consoleWriter, nil)
	case Both:
		Init.Debug().Msg("Logging to both stdout and file")
		setupLoggers(&consoleWriter, fileWriter)
	}

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
		Init.Warn().Msg("Invalid log level, defaulting to 'INFO'")
	}

	envFilePath := os.Getenv(env.LOG_FILE_PATH)
	var filePath string

	if envFilePath != "" {
		filePath = envFilePath
	} else {
		filePath = defaultConfig.FilePath
	}

	currentEnv := os.Getenv(string(env.ENV))
	var destination LogDestination

	if currentEnv == "DEV" {
		destination = Stdout
	} else {
		destination = Both
	}

	return Config{
		Level:       logLevel,
		FilePath:    filePath,
		Destination: destination,
	}
}
