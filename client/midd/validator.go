package midd

import (
	"github.com/go-playground/validator"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidate() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}
