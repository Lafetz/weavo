package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/lafetz/weavo/internal/adapters/web/dto"
	"github.com/lafetz/weavo/internal/adapters/web/webutils"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

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
