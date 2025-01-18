package main

import (
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	mockcache "github.com/lafetz/weavo/internal/adapters/mock_cache"
	openweather "github.com/lafetz/weavo/internal/adapters/open_weather"
	"github.com/lafetz/weavo/internal/adapters/repository"
	"github.com/lafetz/weavo/internal/adapters/web"
	"github.com/lafetz/weavo/internal/adapters/web/webutils"
	"github.com/lafetz/weavo/internal/config"
	"github.com/lafetz/weavo/internal/core/service/location"
	"github.com/lafetz/weavo/internal/core/service/weather"
	customlogger "github.com/lafetz/weavo/internal/logger"
)

const dataRetention = 24 * time.Hour

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Printf("error creating config: %v", err)
		os.Exit(1)
	}
	ow := openweather.NewOpenWeather(config.Open_URL, config.Open_Key, 2)
	logger := customlogger.NewLogger(config.LogLevel, config.Env)
	store := repository.NewInMemoryLocationRepo(dataRetention)
	locationSvc := location.NewService(store)
	mc := mockcache.NewMockCache()
	weatherSvc := weather.NewService(ow, mc)
	val := validator.New()
	custonmVal := webutils.NewCustomValidator(val)
	cookieStore := webutils.CookieStore(dataRetention)
	web := web.NewApp(config.Port, logger, cookieStore, custonmVal, locationSvc, weatherSvc)
	logger.Info("running web server")
	err = web.Run()
	if err != nil {
		logger.Error("web server error", "error", err)
	}
}
