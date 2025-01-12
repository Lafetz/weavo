package weather

import (
	"errors"
	"testing"

	"github.com/lafetz/weavo/internal/core/domain"
)

type MockWeatherProvider struct{}

func (m *MockWeatherProvider) GetWeather(city string) (domain.Weather, error) {
	if city == "London" {
		return domain.Weather{
			Temperature: 15.5,
			Description: "Clear sky",
		}, nil
	}
	return domain.Weather{}, errors.New("weather not found")
}

type MockCache struct {
	data map[string]domain.Weather
}

func (m *MockCache) GetWeather(city string) (domain.Weather, error) {
	weather, exists := m.data[city]
	if exists {
		return weather, nil
	}
	return domain.Weather{}, ErrWeatherNotFound
}

func (m *MockCache) SetWeather(city string, weather domain.Weather) error {
	m.data[city] = weather
	return nil
}

func TestGetWeather_CacheHit(t *testing.T) {
	cache := &MockCache{data: map[string]domain.Weather{
		"London": {Temperature: 15.5, Description: "Clear sky"},
	}}
	provider := &MockWeatherProvider{}
	service := NewService(provider, cache)

	weather, err := service.GetWeather("London")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if weather.Description != "Clear sky" {
		t.Fatalf("expected weather description 'Clear sky', got %v", weather.Description)
	}
}

func TestGetWeather_CacheMissAndAPICall(t *testing.T) {
	cache := &MockCache{data: make(map[string]domain.Weather)}
	provider := &MockWeatherProvider{}
	service := NewService(provider, cache)
	weather, err := service.GetWeather("London")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if weather.Description != "Clear sky" {
		t.Fatalf("expected weather description 'Clear sky', got %v", weather.Description)
	}

	if len(cache.data) == 0 {
		t.Fatalf("expected cache to have data, but it is empty")
	}
}

func TestGetWeather_CacheMissAndAPICallFailure(t *testing.T) {
	cache := &MockCache{data: make(map[string]domain.Weather)} // Empty cache
	provider := &MockWeatherProvider{}
	service := NewService(provider, cache)
	weather, err := service.GetWeather("Paris")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if weather != (domain.Weather{}) {
		t.Fatalf("expected empty weather, got %+v", weather)
	}
}
