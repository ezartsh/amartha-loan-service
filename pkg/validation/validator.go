package validation

import (
	"errors"
	"fmt"
	"loan-service/utils"
	"reflect"
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
	firstValue, secondValue := unPackParams(tagValue)
	switch tag {
	case "required":
		message = fmt.Sprintf("%s is required", name)
	case "numeric":
		message = fmt.Sprintf("%s must be a numeric value", name)
	case "eq":
		message = fmt.Sprintf("%s must be equal to %s", name, tagValue)
	case "gt":
		message = fmt.Sprintf("%s must be greater than %s", name, tagValue)
	case "oneof":
		values := strings.Split(tagValue, " ")
		message = fmt.Sprintf("%s is must be one of %s", name, strings.Join(values, ", "))
	case "min":
		message = fmt.Sprintf("%s must at minimum %v character(s)", name, tagValue)
	case "max":
		message = fmt.Sprintf("%s must at maximum %v character(s)", name, tagValue)
	case "numprec":
		message = fmt.Sprintf("%s not in a valid format. The numeric value must be at precision of %v and scale of %v", name, firstValue, secondValue)
	case "exists":
		message = fmt.Sprintf("%s not exists", name)
	case "exists_in_lowercase":
		message = fmt.Sprintf("%s not exists", name)
	case "tlte_tfield": // local time less than local time to other field ... ( compare Date Time with type LocalTime )
		message = fmt.Sprintf("%s must less than or equal to field %s", name, secondValue)
	case "tgte_tfield": // local time greater than local time to other field ... ( compare Date Time with type LocalTime )
		message = fmt.Sprintf("%s must greater than or equal to field %s", name, secondValue)
	case "dlt_ltfield": // local date less than local time to other field ... ( compare Date with type LocalDate and LocalTime )
		message = fmt.Sprintf("%s must less than field %s", name, secondValue)
	case "dlte_ltfield": // local date less than or equal to local time to other field ... ( compare Date with type LocalDate and LocalTime )
		message = fmt.Sprintf("%s must less than or equal to field %s", name, secondValue)
	case "dgt_ltfield": // local date greater than local time to other field ... ( compare Date with type LocalDate and LocalTime )
		message = fmt.Sprintf("%s must greater than field %s", name, secondValue)
	case "dgte_ltfield": // local date greater or equal than local time to other field ... ( compare Date with type LocalDate and LocalTime )
		message = fmt.Sprintf("%s must greater than or equal to field %s", name, secondValue)
	case "dlt_dfield": // local date less than other field ... ( compare Date with both type LocalDate )
		message = fmt.Sprintf("%s must less than field %s", name, secondValue)
	case "dgt_dfield": // local date greater than other field ... ( compare Date with both type LocalDate )
		message = fmt.Sprintf("%s must greater than field %s", name, secondValue)
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

func unPackTagExistsParams(params string) (string, string, map[string]string) {
	var tableName string
	var columnName string
	var additionalQueries = make(map[string]string)

	paramValues := strings.SplitN(params, ";", 2)
	if len(paramValues) > 0 {
		paramTableAndColumns := strings.Split(paramValues[0], ".")
		if len(paramTableAndColumns) > 0 {
			tableName = paramTableAndColumns[0]
			if len(paramTableAndColumns) > 1 {
				columnName = paramTableAndColumns[1]
			}
		}
		if len(paramValues) > 1 {
			queries := strings.Split(paramValues[1], ";")
			for _, query := range queries {
				q := strings.Split(query, ".")
				if len(q) == 2 {
					additionalQueries[q[0]] = q[1]
				}
			}
		}
	}
	return tableName, columnName, additionalQueries
}

func getDateInstance(fieldType string) interface{} {
	if fieldType == LocalDateType {
		return localDate
	} else if fieldType == LocalTimeType {
		return localTime
	}
	return defaultTime
}

func getDateStringFromReflectValue(fieldValue reflect.Value, fieldType string, dateFormat string) string {
	switch fieldType {
	case LocalDateType:
		if fieldValue.Kind() == reflect.Ptr {
			return fieldValue.Interface().(*utils.LocalDate).Format(dateFormat)
		}
		return fieldValue.Interface().(utils.LocalDate).Format(dateFormat)
	case LocalTimeType:
		if fieldValue.Kind() == reflect.Ptr {
			return fieldValue.Interface().(*utils.LocalTime).Format(dateFormat)
		}
		return fieldValue.Interface().(utils.LocalTime).Format(dateFormat)
	}
	if fieldValue.Kind() == reflect.Ptr {
		return fieldValue.Interface().(*time.Time).Format(dateFormat)
	}
	return fieldValue.Interface().(time.Time).Format(dateFormat)
}
