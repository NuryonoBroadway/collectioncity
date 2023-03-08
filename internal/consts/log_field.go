// Package consts
package consts

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.privy.id/privypass/privypass-package-core/response/errbank"
)

const (
	// LogEventNameServiceTerminated const
	LogEventNameServiceTerminated = "ServiceTerminated"
	//LogEventNameServiceStarting const
	LogEventNameServiceStarting = "ServiceStarting"
)

func ToValidationError(err error, prefix ...string) errbank.ValidationError {
	if err == nil {
		return nil
	}

	ve := make(errbank.ValidationError, 0)

	switch _typed := err.(type) {
	case validation.Errors:
		for k, _err := range _typed {
			switch _err.(type) {
			case validation.Errors:
				ve = append(ve, ToValidationError(_err, fieldName(k, prefix...))...)
			default:
				ve = append(ve, errbank.FieldError{
					Field: fieldName(k, prefix...),
					Error: _err.Error(),
				})
			}
		}

	default:
		ve = append(ve, errbank.FieldError{
			Error: _typed.Error(),
		})
	}

	if len(ve) == 0 {
		return nil
	}

	return ve
}

func fieldName(key string, prefix ...string) string {
	return strings.Join(
		append(prefix, key),
		".",
	)
}
