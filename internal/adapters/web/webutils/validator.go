package webutils

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidatorApi interface {
	ValidateAndRespond(w http.ResponseWriter, input interface{}) bool
}
type CustomValidator struct {
	validate *validator.Validate
}

func NewCustomValidator(validate *validator.Validate) *CustomValidator {

	return &CustomValidator{
		validate: validate,
	}
}

func (v *CustomValidator) ValidateAndRespond(w http.ResponseWriter, input interface{}) bool {
	err := v.validate.Struct(input)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := ValidateModel(validationErrors)
			if err := WriteJSON(w, http.StatusUnprocessableEntity, "validation error", errors, nil); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return true
		}
	}

	return false
}

type ValidationErrorResponse struct {
	StatusCode int         `json:"statusCode"`
	Errors     interface{} `json:"errors"`
}

func ValidateModel(err validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, err := range err {
		errors[strings.ToLower(err.Field())] = errorMsgs(err.Tag(), err.Param())

	}
	return errors

}

func errorMsgs(tag string, value string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "numeric":
		return "must be numeric " + value
	case "lte":
		return "can not be greater than " + value
	case "gte":
		return "can not be less than " + value
	case "len":
		return "length should be equal to " + value

	}
	return ""
}
