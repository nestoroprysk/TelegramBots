package validator

import "gopkg.in/go-playground/validator.v8"

type Validator interface {
	Struct(s interface{}) error
}

type validate struct {
	*validator.Validate
}

func New() Validator {
	conf := &validator.Config{TagName: "validate"}
	return &validate{Validate: validator.New(conf)}
}
