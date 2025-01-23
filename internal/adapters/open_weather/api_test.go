package openweather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

func TestGetWeatherSuccess(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "main": {"temp": 25.0},
            "weather": [{"description": "Clear sky", "icon": "01d", "main": "Clear"}],
            "sys": {"sunrise": 1622520000, "sunset": 1622570400},
            "dt": 1622548800,
            "timezone": 3600,
            "name": "London",
            "units": "metric",
            "coord": {"lon": -0.1257, "lat": 51.5085}
        }`))
	}))
	defer mockServer.Close()

	// Create an instance of OpenWeather with the mock server URL
	ow := NewOpenWeather(mockServer.URL+"?q=%s&appid=%s", "mock-api-key", 10)

	// Call the GetWeather function
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	weather, err := ow.GetWeather(ctx, "London")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedWeather := domain.Weather{
		Description: "Clear sky",
		Condition:   "Clear",
	}
	if weather.Condition != expectedWeather.Condition {
		t.Fatalf("expected condition %s, got %s", expectedWeather.Condition, weather.Condition)
	}
	if weather.Description != expectedWeather.Description {
		t.Fatalf("expected description %s, got %s", expectedWeather.Description, weather.Description)
	}
}

func TestGetWeatherCityNotFound(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	// Create an instance of OpenWeather with the mock server URL
	ow := NewOpenWeather(mockServer.URL+"?q=%s&appid=%s", "mock-api-key", 10)

	// Call the GetWeather function
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := ow.GetWeather(ctx, "UnknownCity")
	if err != weather.ErrCityNotFound {
		t.Fatalf("expected error %v, got %v", weather.ErrCityNotFound, err)
	}
}

func TestGetWeather_InternalServerError(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Create an instance of OpenWeather with the mock server URL
	ow := NewOpenWeather(mockServer.URL+"?q=%s&appid=%s", "mock-api-key", 10)

	// Call the GetWeather function
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := ow.GetWeather(ctx, "London")
	if err == nil || err.Error() != "unexpected status code: 500" {
		t.Fatalf("expected error 'unexpected status code: 500', got %v", err)
	}
}
