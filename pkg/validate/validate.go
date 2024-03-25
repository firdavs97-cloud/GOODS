package validate

import (
	"github.com/go-playground/validator"
)

type ErrorValidateResponse struct {
	FailedField string `json:"failedField"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

type ErrorValidate []*ErrorValidateResponse

func Struct(value interface{}) ErrorValidate {
	var errors []*ErrorValidateResponse
	validate := validator.New()
	err := validate.Struct(value)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorValidateResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (e ErrorValidate) Error() string {
	err := ""
	enter := "\n"
	for i, response := range e {
		if i == len(e)-1 {
			enter = ""
		}
		err += "FailedField: " + response.FailedField + ", Tag: " + response.Tag + ", Value: " + response.Value + enter
	}
	return err
}
