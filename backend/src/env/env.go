package env

import (
	"fmt"
	"os"
)

type EnvVar string

const (
	ENV           EnvVar = "ENV" // DEV, PROD
	API_PORT             = "API_PORT"
	LOG_LEVEL            = "LOG_LEVEL" // TRACE, DEBUG, INFO, WARN, ERROR
	LOG_FILE_PATH        = "LOG_FILE_PATH"
	MONGODB_URI          = "MONGODB_URI"
)

func Validate() error {
	env := os.Getenv(string(ENV))
	if env == "" {
		return fmt.Errorf("environment variable %s is not set", ENV)
	}

	// Check if the environment variable is one of the allowed values
	if env != "DEV" && env != "PROD" {
		return fmt.Errorf("environment variable %s must be either DEV or PROD", ENV)
	}

	envApiPort := os.Getenv(API_PORT)
	if envApiPort == "" {
		return fmt.Errorf("environment variable %s is not set", API_PORT)
	}

	envMongoDBURI := os.Getenv(MONGODB_URI)
	if envMongoDBURI == "" {
		return fmt.Errorf("environment variable %s is not set", MONGODB_URI)
	}

	return nil
}
