package env

import (
	"fmt"
	"os"
)

type EnvVar string

const (
	API_PORT         EnvVar = "API_PORT"
	LOG_LEVEL               = "LOG_LEVEL" // TRACE, DEBUG, INFO, WARN, ERROR
	LOG_FILE_PATH           = "LOG_FILE_PATH"
	MONGODB_URI             = "MONGODB_URI"
	JWT_SECRET              = "JWT_SECRET"
	JWT_EXPIRY_HOURS        = "JWT_EXPIRY_HOURS"
)

func Validate() error {
	envApiPort := os.Getenv(string(API_PORT))
	if envApiPort == "" {
		return fmt.Errorf("environment variable %s is not set", API_PORT)
	}

	envMongoDBURI := os.Getenv(MONGODB_URI)
	if envMongoDBURI == "" {
		return fmt.Errorf("environment variable %s is not set", MONGODB_URI)
	}

	envJwtSecret := os.Getenv(JWT_SECRET)
	if envJwtSecret == "" {
		return fmt.Errorf("environment variable %s is not set", JWT_SECRET)
	}

	envJwtValidHours := os.Getenv(JWT_EXPIRY_HOURS)
	if envJwtValidHours == "" {
		return fmt.Errorf("environment variable %s is not set", JWT_EXPIRY_HOURS)
	}

	return nil
}
