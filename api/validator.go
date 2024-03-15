package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/techschool/simplebank/db/util"
)

var validCurrency validator.Func = func(field validator.FieldLevel) bool {
	if currency, ok := field.Field().Interface().(string); ok {
		//check currency is valid
		return util.IsSupportedCurrency(currency)
	}
	return false

}
