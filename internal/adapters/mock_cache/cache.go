package mockcache

import (
	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

type MockCache struct {
	data map[string]domain.Weather
}

func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]domain.Weather),
	}
}

func (m *MockCache) GetWeather(city string) (domain.Weather, error) {
	wh, exists := m.data[city]
	if !exists {
		return domain.Weather{}, weather.ErrWeatherNotFound
	}
	return wh, nil
}

func (m *MockCache) SetWeather(city string, weather domain.Weather) error {
	m.data[city] = weather
	return nil
}
