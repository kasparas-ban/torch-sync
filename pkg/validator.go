package pkg

import (
	"regexp"

	"github.com/go-playground/validator"
)

var Validate *validator.Validate

func InitializeValidators() {
	Validate = validator.New()
	Validate.RegisterValidation("password", passwordValidation)
}

func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 || len(password) > 30 {
		return false
	}

	tests := []string{"[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}
	for _, test := range tests {
		t, _ := regexp.MatchString(test, password)
		if !t {
			return false
		}
	}
	return true
}
