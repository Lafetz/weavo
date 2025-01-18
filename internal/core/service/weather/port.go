package weather

import (
	"context"

	"github.com/lafetz/weavo/internal/core/domain"
)

type ServiceApi interface {
	GetWeather(ctx context.Context, City string) (domain.Weather, error)
}
type WeatherProvider interface {
	GetWeather(ctx context.Context, city string) (domain.Weather, error)
}
type CachePort interface {
	GetWeather(city string) (domain.Weather, error)
	SetWeather(city string, weather domain.Weather) error
}
