package validators

import (
	"github.com/go-playground/validator/v10"

	helpers "github.com/jeielmosi/peq.nu-backend/src/core/helpers"
)

type HashGetParamDto struct {
	Hash string `json:"hash" validate:"required,availableGetHash"`
}

func AvailableGetHash(val string) bool {
	if len(val) == 0 {
		return false
	}

	alphabeticMp := make(map[rune]bool)
	for _, r := range helpers.FIRST_CHAR_ALPHABET {
		alphabeticMp[r] = true
	}

	if available, ok := alphabeticMp[rune(val[0])]; !ok || !available {
		return false
	}

	for _, r := range helpers.FIRST_CHAR_FORBIDDEN {
		alphabeticMp[r] = true
	}

	for _, r := range val {
		if available, ok := alphabeticMp[r]; !ok || !available {
			return false
		}
	}

	return true
}

func AvailableGetHashValidator() (string, func(validator.FieldLevel) bool) {
	return "availableGetHash", func(fl validator.FieldLevel) bool {
		val, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		return AvailableGetHash(val)
	}
}
