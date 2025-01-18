package dto

import "github.com/lafetz/weavo/internal/core/domain"

type WeatherRes struct {
	Icon        string  `json:"icon"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Condition   string  `json:"condition"`
	DateTime    string  `json:"date_time"`
	Location    string  `json:"location"`
	Units       string  `json:"units"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

func GetWeatherRes(w domain.Weather) WeatherRes {
	return WeatherRes{
		Icon:        w.Icon,
		Temperature: w.Temperature,
		Description: w.Description,
		Condition:   w.Condition,
		DateTime:    w.DateTime,
		Location:    w.Location,
		Units:       w.Units,
		Lat:         w.Lat,
		Lon:         w.Lon,
	}
}
