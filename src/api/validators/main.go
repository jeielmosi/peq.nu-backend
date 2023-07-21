package validators

import "github.com/go-playground/validator/v10"

func GetValidator() *validator.Validate {
	validate := validator.New()
	{
		str, validatorFn := AvailableGetHashValidator()
		validate.RegisterValidation(str, validatorFn)
	}

	{
		str, validatorFn := AvailablePostHashValidator()
		validate.RegisterValidation(str, validatorFn)
	}

	return validate
}
