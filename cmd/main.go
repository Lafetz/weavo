package main

import (
	"context"
	"fmt"
	"log"
	"os"

	mockcache "github.com/lafetz/weavo/internal/adapters/mock_cache"
	openweather "github.com/lafetz/weavo/internal/adapters/open_weather"
	"github.com/lafetz/weavo/internal/adapters/repository"
	"github.com/lafetz/weavo/internal/config"
	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/location"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Printf("error creating config: %v", err)
		os.Exit(1)
	}
	ow := openweather.NewOpenWeather(config.Open_URL, config.Open_Key, 2)
	store := repository.NewInMemoryLocationRepo()
	locationSvc := location.NewService(store)
	mc := mockcache.NewMockCache()
	locationSvc.CreateLocation(context.Background(), domain.Location{
		UserID: "1",
		Notes:  "test",
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
		Nickname: "test",
		City:     "test",
	})
	weatherSvc := weather.NewService(ow, mc)
	fmt.Println(weatherSvc)
}
