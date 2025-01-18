package webutils

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestStruct struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		expectedError string
	}{
		{
			name:          "Valid JSON",
			body:          `{"name": "John", "email": "john@example.com"}`,
			expectedError: "",
		},
		{
			name:          "Malformed JSON",
			body:          `{"name": "John", "email": "john@example.com"`,
			expectedError: "body contains badly-formed JSON",
		},
		{
			name:          "Unknown Field",
			body:          `{"name": "John", "email": "john@example.com", "age": 30}`,
			expectedError: "body contains unknown key \"age\"",
		},
		{
			name:          "Incorrect JSON Type",
			body:          `{"name": "John", "email": 123}`,
			expectedError: "body contains incorrect JSON type for field \"email\"",
		},
		{
			name:          "Empty Body",
			body:          ``,
			expectedError: "body must not be empty",
		},
		{
			name:          "Multiple JSON Values",
			body:          `{"name": "John"} {"email": "john@example.com"}`,
			expectedError: "body must only contain a single JSON value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			var dst TestStruct
			err := ReadJSON(w, req, &dst)

			if tt.expectedError == "" && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else if !errors.Is(err, errors.New(tt.expectedError)) && err.Error() != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			}
		})
	}
}
