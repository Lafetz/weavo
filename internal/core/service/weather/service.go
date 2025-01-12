package weather

import (
	"errors"

	"github.com/lafetz/weavo/internal/core/domain"
)

var (
	ErrWeatherNotFound = errors.New("weather not found")
	ErrCityNotFound    = errors.New("city not found")
)

type Service struct {
	weatherProvider WeatherProvider
	cache           CachePort
}

func NewService(weatherProvider WeatherProvider, cache CachePort) *Service {
	return &Service{weatherProvider: weatherProvider, cache: cache}
}

func (s *Service) GetWeather(City string) (domain.Weather, error) {
	weather, err := s.cache.GetWeather(City)
	if err == nil {
		return weather, nil
	}
	if !errors.Is(err, ErrWeatherNotFound) { // if the error is not a cache miss, return the error
		return domain.Weather{}, err
	}
	weather, err = s.weatherProvider.GetWeather(City)
	if err != nil {
		return domain.Weather{}, err
	}
	err = s.cache.SetWeather(City, weather)
	if err != nil {
		return domain.Weather{}, err
	}

	return weather, nil
}
