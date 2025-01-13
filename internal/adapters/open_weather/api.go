package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

const metricUnit = "metric"

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
	url    string
	key    string
	client *http.Client
}

func NewOpenWeather(url, key string, timeoutS int) *OpenWeather {
	return &OpenWeather{
		url: url,
		key: key,
		client: &http.Client{
			Timeout: time.Duration(timeoutS) * time.Second,
		},
	}
}
func (o OpenWeather) GetWeather(ctx context.Context, city string) (domain.Weather, error) {
	url := fmt.Sprintf(o.url, city, o.key)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return domain.Weather{}, weather.ErrCityNotFound
	}
	if resp.StatusCode >= 300 {
		return domain.Weather{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var apiResp WeatherAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return domain.Weather{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	dateTime := time.Unix(apiResp.Dt, 0).Format("2006-01-02 15:04:05")
	weather := domain.Weather{
		Temperature: apiResp.Main.Temp,
		Description: apiResp.Weather[0].Description,
		Condition:   apiResp.Weather[0].Main,
		Icon:        apiResp.Weather[0].Icon,
		DateTime:    dateTime,
		Location:    apiResp.Name,
		Units:       metricUnit,
		Lat:         apiResp.Coord.Lat,
		Lon:         apiResp.Coord.Lon,
	}
	return weather, nil
}
