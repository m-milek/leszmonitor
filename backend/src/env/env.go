package env

import (
	"fmt"
	"os"
)

const (
	ApiPort        = "API_PORT"
	LogLevel       = "LOG_LEVEL" // TRACE, DEBUG, INFO, WARN, ERROR
	LogFilePath    = "LOG_FILE_PATH"
	PostgresURI    = "POSTGRES_URI"
	JwtSecret      = "JWT_SECRET"
	JwtExpiryHours = "JWT_EXPIRY_HOURS"
)

func Validate() error {
	envAPIPort := os.Getenv(ApiPort)
	if envAPIPort == "" {
		return fmt.Errorf("environment variable %s is not set", ApiPort)
	}

	envPostgresURI := os.Getenv(PostgresURI)
	if envPostgresURI == "" {
		return fmt.Errorf("environment variable %s is not set", PostgresURI)
	}

	envJwtSecret := os.Getenv(JwtSecret)
	if envJwtSecret == "" {
		return fmt.Errorf("environment variable %s is not set", JwtSecret)
	}

	envJwtValidHours := os.Getenv(JwtExpiryHours)
	if envJwtValidHours == "" {
		return fmt.Errorf("environment variable %s is not set", JwtExpiryHours)
	}

	return nil
}
