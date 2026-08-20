package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

const (
	APIPort               = "API_PORT"
	LogLevel              = "LOG_LEVEL"  // TRACE, DEBUG, INFO, WARN, ERROR
	LogFormat             = "LOG_FORMAT" // JSON or CONSOLE
	SqliteDBPath          = "SQLITE_DB_PATH"
	JwtSecret             = "JWT_SECRET"
	JwtExpiryHours        = "JWT_EXPIRY_HOURS"
	InstanceAdminUsername = "INSTANCE_ADMIN_USERNAME"
	InstanceAdminPassword = "INSTANCE_ADMIN_PASSWORD"
	EnableLogdy           = "ENABLE_LOGDY"
	LogdyConfigFilePath   = "LOGDY_CONFIG_FILE_PATH"
)

func Validate() error {
	log.Info().Msg("Validating environment variables...")
	var missingVars []string

	if os.Getenv(APIPort) == "" {
		missingVars = append(missingVars, APIPort)
	}

	if os.Getenv(SqliteDBPath) == "" {
		missingVars = append(missingVars, SqliteDBPath)
	}

	if os.Getenv(JwtSecret) == "" {
		missingVars = append(missingVars, JwtSecret)
	}

	if os.Getenv(JwtExpiryHours) == "" {
		missingVars = append(missingVars, JwtExpiryHours)
	}

	if os.Getenv(InstanceAdminUsername) == "" {
		missingVars = append(missingVars, InstanceAdminUsername)
	}

	if os.Getenv(InstanceAdminPassword) == "" {
		missingVars = append(missingVars, InstanceAdminPassword)
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing environment variables: %v", missingVars)
	}

	return nil
}
