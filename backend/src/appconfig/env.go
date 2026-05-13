package config

import (
	"fmt"
	"os"
)

const (
	ApiPort        = "API_PORT"
	LogLevel       = "LOG_LEVEL" // TRACE, DEBUG, INFO, WARN, ERROR
	SqliteDbPath   = "SQLITE_DB_PATH"
	JwtSecret      = "JWT_SECRET"
	JwtExpiryHours = "JWT_EXPIRY_HOURS"
)

func Validate() error {
	var missingVars []string

	if os.Getenv(ApiPort) == "" {
		missingVars = append(missingVars, ApiPort)
	}

	if os.Getenv(SqliteDbPath) == "" {
		missingVars = append(missingVars, SqliteDbPath)
	}

	if os.Getenv(JwtSecret) == "" {
		missingVars = append(missingVars, JwtSecret)
	}

	if os.Getenv(JwtExpiryHours) == "" {
		missingVars = append(missingVars, JwtExpiryHours)
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing environment variables: %v", missingVars)
	}

	return nil
}
