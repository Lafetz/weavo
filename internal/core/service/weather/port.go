package weather

import "github.com/lafetz/weavo/internal/core/domain"

type WeatherProvider interface {
	GetWeather(City string) (domain.Weather, error)
}
type CachePort interface {
	GetWeather(City string) (domain.Weather, error)
	SetWeather(City string, weather domain.Weather) error
}
