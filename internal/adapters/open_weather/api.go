package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

type WeatherAPIResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Main        string `json:"main"`
	} `json:"weather"`
	Sys struct {
		Sunrise int64 `json:"sunrise"`
		Sunset  int64 `json:"sunset"`
	} `json:"sys"`
	Dt       int64  `json:"dt"`
	Timezone int    `json:"timezone"`
	Name     string `json:"name"`
	Units    string `json:"units"`
	Coord    Coord  `json:"coord"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
type OpenWeather struct {
	Url string
	key string
}

func NewOPenWeather(url, key string) *OpenWeather {
	return &OpenWeather{Url: url, key: key}
}
func (o OpenWeather) GetWeather(city string) (domain.Weather, error) {
	url := fmt.Sprintf(o.Url, city, o.key)
	resp, err := http.Get(url)
	if err != nil {
		return domain.Weather{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Weather{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return domain.Weather{}, weather.ErrCityNotFound
	}
	if resp.StatusCode >= 300 {
		return domain.Weather{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}
	var apiResp WeatherAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return domain.Weather{}, err
	}

	dateTime := time.Unix(apiResp.Dt, 0).Format("2006-01-02 15:04:05")
	fmt.Println("timezone", apiResp.Timezone)
	weather := domain.Weather{
		Temperature: apiResp.Main.Temp,
		Description: apiResp.Weather[0].Description,
		Condition:   apiResp.Weather[0].Main,
		Icon:        apiResp.Weather[0].Icon,
		DateTime:    dateTime,
		Location:    apiResp.Name,
		Units:       "metric",
		Lat:         apiResp.Coord.Lat,
		Lon:         apiResp.Coord.Lon,
	}

	return weather, nil
}
