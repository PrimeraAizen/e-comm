package domain

import (
	"github.com/go-playground/validator/v10"
)

type Example struct {
	ExampleField string `json:"example_field" validate:"required"`
}

func (e *Example) Validate() error {
	validate := validator.New()
	if err := validate.Struct(e); err != nil {
		return ErrValidation
	}
	return nil
}
