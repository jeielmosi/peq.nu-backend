package validators

import (
	"github.com/go-playground/validator/v10"

	helpers "github.com/jeielmosi/peq.nu-backend/src/core/helpers"
)

type HashPostParamDto struct {
	Hash string `json:"hash" validate:"required,availablePostHash"`
}

func AvailablePostHash(val string) bool {
	if len(val) == helpers.HASH_SIZE {
		return false
	}

	return AvailableGetHash(val)
}

func AvailablePostHashValidator() (string, func(validator.FieldLevel) bool) {
	return "availablePostHash", func(fl validator.FieldLevel) bool {
		val, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		return AvailablePostHash(val)
	}
}
