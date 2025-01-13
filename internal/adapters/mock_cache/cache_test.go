package mockcache

import (
	"errors"
	"testing"

	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

func TestMockCache_GetWeather(t *testing.T) {
	cache := NewMockCache()
	weatherData := domain.Weather{
		Icon:        "01d",
		Temperature: 25.0,
		Description: "Clear sky",
		Condition:   "Clear",
		DateTime:    "2023-10-10 10:00:00",
		Location:    "London",
		Units:       "metric",
	}
	cache.SetWeather("London", weatherData)

	t.Run("existing city", func(t *testing.T) {
		weather, err := cache.GetWeather("London")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if weather != weatherData {
			t.Fatalf("expected weather %+v, got %+v", weatherData, weather)
		}
	})

	t.Run("non-existing city", func(t *testing.T) {
		_, err := cache.GetWeather("Paris")
		if !errors.Is(err, weather.ErrWeatherNotFound) {
			t.Fatal("expected error, got nil")
		}
		if err != weather.ErrWeatherNotFound {
			t.Fatalf("expected error %v, got %v", weather.ErrWeatherNotFound, err)
		}
	})
}

func TestMockCache_SetWeather(t *testing.T) {
	cache := NewMockCache()
	weatherData := domain.Weather{
		Icon:        "01d",
		Temperature: 25.0,
		Description: "Clear sky",
		Condition:   "Clear",
		DateTime:    "2023-10-10 10:00:00",
		Location:    "London",
		Units:       "metric",
	}

	err := cache.SetWeather("London", weatherData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedWeather, err := cache.GetWeather("London")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedWeather != weatherData {
		t.Fatalf("expected weather %+v, got %+v", weatherData, retrievedWeather)
	}
}
