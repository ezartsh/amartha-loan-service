package validation

import (
	"errors"
	"fmt"
	"loan-service/utils"
	"strings"
	"time"
)

var (
	ErrValidationInput = errors.New("validation input failed.")
)

const (
	LocalDateType string = "util.LocalDate"
	LocalTimeType string = "util.LocalTime"
)

var (
	DefaultPrice float64 = 100

	localDate   = utils.LocalDate{Time: time.Now()}
	localTime   = utils.LocalTime{Time: time.Now()}
	defaultTime = time.Now()
)

const (
	Required string = "required"
	File     string = "file"
	Numeric  string = "numeric"
)

type (
	// Rules with Struct Field Name as the key, and the validation rule as the value.
	Rules map[string]string

	ValidationErrors map[string]string
)

func (r *Rules) Add(key string, schema string) {
	(*r)[key] = schema
}

func (r *ValidationErrors) SetErrors(errs ValidationErrors) {
	(*r) = errs
}

func (r *ValidationErrors) Add(key string, schema string) {
	(*r)[key] = schema
}

func (r *ValidationErrors) Has(key string) bool {
	if _, ok := (*r)[key]; ok {
		return true
	}
	return false
}

func (r *ValidationErrors) Clears() {
	(*r) = make(ValidationErrors)
}

func getErrorValidationMessage(name string, tag string, tagValue string) (message string) {
	switch tag {
	case "required":
		message = fmt.Sprintf("%s is required", name)
	case "numeric":
		message = fmt.Sprintf("%s must be a numeric value", name)
	case "gte":
		message = fmt.Sprintf("%s must be greater than or equal to %s", name, tagValue)
	case "lte":
		message = fmt.Sprintf("%s must be less than or equal to %s", name, tagValue)
	}
	return message
}

func unPackParams(params string) (string, string) {
	var parameters [2]string = [2]string{"", ""}

	paramValues := strings.Split(params, ".")
	if len(paramValues) > 0 {
		parameters[0] = paramValues[0]
		if len(paramValues) > 1 {
			parameters[1] = paramValues[1]
		}
	}
	return parameters[0], parameters[1]
}
