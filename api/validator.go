package api

import (
	"github.com/go-playground/validator/v10"
	"go_challenge/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsCurrencySupported(currency)
	}
	return false
}
