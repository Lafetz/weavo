package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/lafetz/weavo/internal/adapters/web/dto"
	"github.com/lafetz/weavo/internal/adapters/web/webutils"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

// GetWeather handles the HTTP request to retrieve weather information for a given city.
//
// @Summary Get weather information
// @Description Retrieves weather information for a specified city.
// @Tags weather
// @Accept json
// @Produce json
// @Param city query string true "City name"
// @Success 200 {object} dto.WeatherRes "weather retrieved successfully"
// @Failure 400 {string} string "invalid city"
// @Failure 404 {string} string "city not found"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/weather [get]
func GetWeather(weatherSvc weather.ServiceApi, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		if city == "" {
			webutils.WriteJSON(w, http.StatusBadRequest, "invalid city", nil, nil)
			return
		}

		weatherData, err := weatherSvc.GetWeather(r.Context(), city)
		if err != nil {
			if errors.Is(err, weather.ErrCityNotFound) {
				webutils.WriteJSON(w, http.StatusNotFound, "city not found", nil, nil)
				return
			}
			webutils.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil, nil)
			logger.Error("error on getting weather by city", "error", err.Error())
			return
		}

		webutils.WriteJSON(w, http.StatusOK, "weather retrieved successfully", dto.GetWeatherRes(weatherData), nil)
	}
}
