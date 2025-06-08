package env

import (
	"fmt"
	"os"
)

const (
	ApiPort        = "API_PORT"
	LogLevel       = "LOG_LEVEL" // TRACE, DEBUG, INFO, WARN, ERROR
	LogFilePath    = "LOG_FILE_PATH"
	MongoDbUri     = "MONGODB_URI"
	JwtSecret      = "JWT_SECRET"
	JwtExpiryHours = "JWT_EXPIRY_HOURS"
)

func Validate() error {
	envApiPort := os.Getenv(ApiPort)
	if envApiPort == "" {
		return fmt.Errorf("environment variable %s is not set", ApiPort)
	}

	envMongoDBURI := os.Getenv(MongoDbUri)
	if envMongoDBURI == "" {
		return fmt.Errorf("environment variable %s is not set", MongoDbUri)
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
