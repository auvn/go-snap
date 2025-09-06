package snapx

import (
	"errors"

	validator "gopkg.in/validator.v2"
)

var _defaultValidator = func() *validator.Validator {
	return validator.NewValidator()
}()

type Validator interface {
	Validate() error
}

func validate(dest any) error {
	v, ok := dest.(Validator)
	if ok {
		if err := v.Validate(); err != nil {
			return err
		}
		return nil
	}

	if err := _defaultValidator.Validate(
		dest,
	); err != nil && !errors.Is(err, validator.ErrUnsupported) {
		return err
	}

	return nil
}
