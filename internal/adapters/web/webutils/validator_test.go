package webutils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

type TestStructV struct {
	Name  string `validate:"required"`
	Age   int    `validate:"numeric,gte=0,lte=130"`
	Code  string `validate:"len=5"`
	Score int    `validate:"gte=0,lte=100"`
}

func TestCustomValidator_ValidateAndRespond_Success(t *testing.T) {
	validate := validator.New()
	validator := NewCustomValidator(validate)

	rr := httptest.NewRecorder()
	input := TestStructV{
		Name:  "John Doe",
		Age:   30,
		Code:  "12345",
		Score: 90,
	}

	if validator.ValidateAndRespond(rr, input) {
		t.Fatalf("expected no validation errors, got validation errors")
	}

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, status)
	}
}

func TestCustomValidator_ValidateAndRespond_ValidationError(t *testing.T) {
	validate := validator.New()
	validator := NewCustomValidator(validate)

	rr := httptest.NewRecorder()
	input := TestStructV{
		Name:  "",
		Age:   150,
		Code:  "123",
		Score: -10,
	}

	if !validator.ValidateAndRespond(rr, input) {
		t.Fatalf("expected validation errors, got no validation errors")
	}

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Fatalf("expected status code %d, got %d", http.StatusUnprocessableEntity, status)
	}

	expectedResponse := `{
        "status": "error",
        "message": "validation error",
        "data": {
            "age": "can not be greater than 130",
            "code": "length should be equal to 5",
            "name": "This field is required",
            "score": "can not be less than 0"
        }
    }`

	actualResponse := rr.Body.String()
	if !strings.Contains(actualResponse, `"status": "error"`) ||
		!strings.Contains(actualResponse, `"message": "validation error"`) ||
		!strings.Contains(actualResponse, `"age": "can not be greater than 130"`) ||
		!strings.Contains(actualResponse, `"code": "length should be equal to 5"`) ||
		!strings.Contains(actualResponse, `"name": "This field is required"`) ||
		!strings.Contains(actualResponse, `"score": "can not be less than 0"`) {
		t.Fatalf("expected response %s, got %s", expectedResponse, actualResponse)
	}
}

func TestValidateModel(t *testing.T) {
	validate := validator.New()
	input := TestStructV{
		Name:  "",
		Age:   150,
		Code:  "123",
		Score: -10,
	}

	err := validate.Struct(input)
	if err == nil {
		t.Fatalf("expected validation errors, got no validation errors")
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("expected validator.ValidationErrors, got %T", err)
	}

	errors := ValidateModel(validationErrors)
	expectedErrors := map[string]string{
		"name":  "This field is required",
		"age":   "can not be greater than 130",
		"code":  "length should be equal to 5",
		"score": "can not be less than 0",
	}

	for field, expectedError := range expectedErrors {
		if errors[field] != expectedError {
			t.Fatalf("expected error for field %s: %s, got %s", field, expectedError, errors[field])
		}
	}
}
