// Package log provides a customized logging setup using zerolog.
//
// This package initializes several pre-configured loggers for different
// components of the application. Each logger uses a customized console output
// with colored level indicators and component name prefixes.
//
// Usage:
//
//	log.Main.Info().Msg("Application started")
//	log.Api.Error().Err(err).Msg("Failed to process request")
//	log.Uptime.Debug().Int("status", status).Msg("Health check completed")
package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

func SetupLogger(componentName string) zerolog.Logger {
	var output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		level := strings.ToUpper(fmt.Sprintf("%s", i))

		if i == nil {
			level = "LOG"
		}

		if output.NoColor {
			return fmt.Sprintf("[%s] %-5s", componentName, level)
		}

		// Define color codes (without bold)
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

		// Bold code is separate
		boldCode := "\033[1m"
		resetCode := "\033[0m"

		// Apply bold and color to the level only
		// Format: [component] bold+color+LEVEL+reset
		return fmt.Sprintf("[%s] %s%s%-5s%s",
			componentName,
			boldCode,  // Apply bold
			colorCode, // Apply color
			level,     // The level text
			resetCode) // Reset all formatting
	}

	var logger = zerolog.New(output).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	return logger
}

// Loggers initialization and export
var (
	Main   = SetupLogger("main")
	Api    = SetupLogger("http")
	Logs   = SetupLogger("logs")
	Uptime = SetupLogger("uptm")
)
