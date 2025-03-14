package validation

import (
	"errors"
	"loan-service/pkg/logger"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type SchemaRules = map[string][]string

type Valid struct {
	validator *validator.Validate
	rules     Rules
}

func NewStructValidation(schemaRules SchemaRules) *Valid {

	// Initiate validation instance
	vd := validator.New()

	// Extend handler function err.Field() to return json name from struct property
	vd.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	rules := make(Rules)

	for key, rule := range schemaRules {
		rules.Add(key, strings.Join(rule, ","))
	}

	return &Valid{
		validator: vd,
		rules:     rules,
	}
}

func (v *Valid) Validate(data any) (ValidationErrors, error) {

	v.validator.RegisterStructValidationMapRules(v.rules, data)

	err := v.validator.Struct(data)
	var inputErrors = make(ValidationErrors)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.AppLog.Error(err, "Invalid validation error")
			return inputErrors, ErrValidationInput
		}
		for _, err := range err.(validator.ValidationErrors) {
			inputErrors.Add(err.Field(), getErrorValidationMessage(err.Field(), err.Tag(), err.Param()))
		}
		return inputErrors, ErrValidationInput
	}

	return inputErrors, nil

}

func (v *Valid) ValidateStopOnFirstFailure(data any) error {

	v.validator.RegisterStructValidationMapRules(v.rules, data)

	err := v.validator.Struct(data)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.AppLog.Error(err, "Invalid validation error")
			return ErrValidationInput
		}
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(getErrorValidationMessage(err.Field(), err.Tag(), err.Param()))
		}
	}

	return nil

}
