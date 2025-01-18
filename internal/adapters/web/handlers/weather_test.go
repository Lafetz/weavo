package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lafetz/weavo/internal/adapters/web/dto"
	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

type MockWeatherService struct{}

func (m *MockWeatherService) GetWeather(ctx context.Context, city string) (domain.Weather, error) {
	if city == "nonexistent" {
		return domain.Weather{}, weather.ErrCityNotFound
	}
	return domain.Weather{
		Location:    city,
		Temperature: 25.0,
		Description: "Clear",
	}, nil
}

func TestGetWeather(t *testing.T) {
	mockSvc := &MockWeatherService{}
	handler := GetWeather(mockSvc, slog.Default())

	t.Run("missing city query parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/weather", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("invalid city")) {
			t.Errorf("Expected response body to contain 'invalid city', got %s", w.Body.String())
		}
	})

	t.Run("nonexistent city", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/weather?city=nonexistent", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("city not found")) {
			t.Errorf("Expected response body to contain 'city not found', got %s", w.Body.String())
		}
	})

	t.Run("successful weather retrieval", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/weather?city=TestCity", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response struct {
			Message string         `json:"message"`
			Data    dto.WeatherRes `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		if response.Message != "weather retrieved successfully" {
			t.Errorf("Expected message 'weather retrieved successfully', got %s", response.Message)
		}
		if response.Data.Location != "TestCity" {
			t.Errorf("Expected city 'TestCity', got %s", response.Data.Location)
		}
	})
}
