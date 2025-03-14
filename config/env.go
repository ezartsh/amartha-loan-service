package config

import (
	"errors"
	"loan-service/utils"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	envVarHttpPort         = "HTTP_PORT"
	envVarLogDirectoryPath = "LOG_DIRECTORY_PATH"
)

type EnvValidationError struct {
	VariableName string `json:"variable_name"`
	Message      string `json:"message"`
	Description  string `json:"description"`
}

// EnvType is variable in .env file
type EnvConfig struct {
	HTTPPort         int
	LogDirectoryPath string
}

// Env is global var for EnvType
var Env = EnvConfig{}

func Init() (results map[string]any, errorBags []EnvValidationError, err error) {
	godotenv.Load()

	projectDirectory, _ := os.Getwd()

	httpPortConfig := getEnvAsInt(envVarHttpPort, 0)
	Env.LogDirectoryPath = getEnv(envVarLogDirectoryPath, utils.ConcatPaths(projectDirectory, "logs"))

	if httpPortConfig == 0 {
		errorBags = append(errorBags, EnvValidationError{
			VariableName: envVarHttpPort,
			Message:      "http port is required",
		})
	}

	if len(errorBags) > 0 {
		err = errors.New("config validation failed")
		return
	}

	Env.HTTPPort = httpPortConfig

	results = map[string]any{
		envVarHttpPort:         Env.HTTPPort,
		envVarLogDirectoryPath: Env.LogDirectoryPath,
	}

	return
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}
