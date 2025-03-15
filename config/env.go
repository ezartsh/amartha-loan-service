package config

import (
	"errors"
	"loan-service/utils"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	envVarHttpPort               = "HTTP_PORT"
	envVarCompanyMarginInPercent = "COMPANY_MARGIN_IN_PERCENT"
	envVarMinimumLoanAmount      = "MINIMUM_LOAN_AMOUNT"
	envVarMaximumLoanAmount      = "MAXIMUM_LOAN_AMOUNT"
	envVarLogDirectoryPath       = "LOG_DIRECTORY_PATH"
)

type EnvValidationError struct {
	VariableName string `json:"variable_name"`
	Message      string `json:"message"`
	Description  string `json:"description"`
}

// EnvType is variable in .env file
type EnvConfig struct {
	HTTPPort               int
	LogDirectoryPath       string
	CompanyMarginInPercent float64
	MinimumLoanAmout       float64
	MaximumLoanAmout       float64
}

// Env is global var for EnvType
var Env = EnvConfig{}

func Init() (results map[string]any, errorBags []EnvValidationError, err error) {
	godotenv.Load()

	projectDirectory, _ := os.Getwd()

	Env.HTTPPort = getEnvAsInt(envVarHttpPort, 0)
	Env.LogDirectoryPath = getEnv(envVarLogDirectoryPath, utils.ConcatPaths(projectDirectory, "logs"))
	Env.CompanyMarginInPercent = getEnvAsFloat64(envVarCompanyMarginInPercent, 0)
	Env.MinimumLoanAmout = getEnvAsFloat64(envVarMinimumLoanAmount, 0)
	Env.MaximumLoanAmout = getEnvAsFloat64(envVarMaximumLoanAmount, 0)

	if Env.HTTPPort == 0 {
		errorBags = append(errorBags, EnvValidationError{
			VariableName: envVarHttpPort,
			Message:      "http port is required",
		})
	}

	if Env.CompanyMarginInPercent < 1 {
		errorBags = append(errorBags, EnvValidationError{
			VariableName: envVarCompanyMarginInPercent,
			Message:      "company margin in percent is required and must greater than or equal to 1",
		})
	}

	if Env.MinimumLoanAmout < 10_000 {
		errorBags = append(errorBags, EnvValidationError{
			VariableName: envVarMinimumLoanAmount,
			Message:      "minimum loan amount is required and must not be less than 10000",
		})
	}

	if Env.MaximumLoanAmout <= Env.MinimumLoanAmout {
		errorBags = append(errorBags, EnvValidationError{
			VariableName: envVarMaximumLoanAmount,
			Message:      "maximum loan amount is required and must be greater than minimum amount",
		})
	}

	if len(errorBags) > 0 {
		err = errors.New("config validation failed")
		return
	}

	results = map[string]any{
		envVarHttpPort:               Env.HTTPPort,
		envVarLogDirectoryPath:       Env.LogDirectoryPath,
		envVarCompanyMarginInPercent: Env.CompanyMarginInPercent,
		envVarMinimumLoanAmount:      Env.MinimumLoanAmout,
		envVarMaximumLoanAmount:      Env.MaximumLoanAmout,
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

func getEnvAsFloat64(name string, defaultValue float64) float64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}

	return defaultValue
}
