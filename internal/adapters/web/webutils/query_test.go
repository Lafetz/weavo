package webutils

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetQueryInt(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		key          string
		defaultValue int
		expected     int
	}{
		{
			name:         "Key exists with valid integer",
			query:        "age=25",
			key:          "age",
			defaultValue: 0,
			expected:     25,
		},
		{
			name:         "Key exists with invalid integer",
			query:        "age=abc",
			key:          "age",
			defaultValue: 0,
			expected:     0,
		},
		{
			name:         "Key does not exist",
			query:        "name=John",
			key:          "age",
			defaultValue: 30,
			expected:     30,
		},
		{
			name:         "Empty query string",
			query:        "",
			key:          "age",
			defaultValue: 20,
			expected:     20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{
					RawQuery: tt.query,
				},
			}
			result := GetQueryInt(req, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
func TestGetQueryString(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "Key exists with value",
			query:        "name=John",
			key:          "name",
			defaultValue: "default",
			expected:     "John",
		},
		{
			name:         "Key does not exist",
			query:        "age=25",
			key:          "name",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Empty query string",
			query:        "",
			key:          "name",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Key exists with empty value",
			query:        "name=",
			key:          "name",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{
					RawQuery: tt.query,
				},
			}
			result := GetQueryString(req, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
