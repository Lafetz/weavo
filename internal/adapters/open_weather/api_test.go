package openweather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

func TestGetWeather_Success(t *testing.T) {
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
		Temperature: 25.0,
		Description: "Clear sky",
		Condition:   "Clear",
		Icon:        "01d",
		DateTime:    "2021-06-01 15:00:00", // Adjusted to Local time
		Location:    "London",
		Units:       metricUnit,
		Lat:         51.5085,
		Lon:         -0.1257,
	}
	if !reflect.DeepEqual(weather, expectedWeather) {
		t.Fatalf("expected %+v, got %+v", expectedWeather, weather)
	}
}

func TestGetWeather_CityNotFound(t *testing.T) {
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
