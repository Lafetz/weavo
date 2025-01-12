package main

import (
	"log"
	"os"

	openweather "github.com/lafetz/weavo/internal/adapters/open_weather"
	"github.com/lafetz/weavo/internal/config"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Printf("error creating config: %v", err)
		os.Exit(1)
	}
	_ = openweather.NewOPenWeather(config.Open_URL, config.Open_Key)

}
