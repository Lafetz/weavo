package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lafetz/weavo/internal/adapters/web/dto"
	"github.com/lafetz/weavo/internal/adapters/web/webutils"
	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/location"
)

type MockLocationService struct{}

func NewMockLocationService() *MockLocationService {
	return &MockLocationService{}
}

func (m *MockLocationService) CreateLocation(ctx context.Context, loc domain.Location) (domain.Location, error) {
	loc.Id = uuid.New().String()
	return loc, nil
}

func (m *MockLocationService) GetLocation(ctx context.Context, id string) (domain.Location, error) {
	if id == "notfound" {
		return domain.Location{}, location.ErrLocationNotFound
	}
	return domain.Location{
		Id:       id,
		UserID:   "1",
		Notes:    "Test Notes",
		Nickname: "Test Nickname",
		City:     "Test City",
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}, nil
}

func (m *MockLocationService) UpdateLocation(ctx context.Context, loc domain.Location) (domain.Location, error) {
	if loc.Id == "notfound" {
		return domain.Location{}, location.ErrLocationNotFound
	}
	return loc, nil
}

func (m *MockLocationService) DeleteLocation(ctx context.Context, id string, userId string) error {
	if id == "notfound" {
		return location.ErrLocationNotFound
	}
	return nil
}

func (m *MockLocationService) GetLocations(ctx context.Context, userID string, filter location.Filter) ([]domain.Location, domain.Metadata, error) {

	locations := []domain.Location{
		{Id: uuid.New().String(), UserID: "1", Notes: "Test Notes", Nickname: "Test Nickname", City: "Test City", Coordinates: domain.Coordinates{Lat: 1.0, Lon: 1.0}},
		{Id: uuid.New().String(), UserID: "1", Notes: "Test Notes 2", Nickname: "Test Nickname 2", City: "Test City 2", Coordinates: domain.Coordinates{Lat: 2.0, Lon: 2.0}},
	}

	metadata := domain.Metadata{
		TotalRecords: int32(len(locations)),
		CurrentPage:  int32(filter.Page),
		LastPage:     int32((len(locations) + int(filter.PageSize) - 1) / int(filter.PageSize)),
	}

	return locations, metadata, nil
}

func TestCreateLocation(t *testing.T) {
	mockSvc := NewMockLocationService()
	handler := CreateLocation(mockSvc, slog.Default(), webutils.NewCustomValidator(validator.New()))
	ctx := context.WithValue(context.Background(), "userId", "1")

	t.Run("invalid input format", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/locations", reqBody).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("Invalid input format")) {
			t.Errorf("Expected response body to contain 'Invalid input format', got %s", w.Body.String())
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		createLocation := dto.LocationReq{
			Notes:    "",
			Nickname: "Test Nickname",
			City:     "Test City",
			Coordinates: dto.Coordinates{
				Lat: 1.0,
				Lon: 1.0,
			},
		}

		reqBody, _ := json.Marshal(createLocation)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/locations", bytes.NewBuffer(reqBody)).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("This field is required")) {
			t.Errorf("Expected response body to contain 'This field is required', got %s", w.Body.String())
		}
	})

	createLocation := dto.LocationReq{
		UserID:   "1",
		Notes:    "Test Notes",
		Nickname: "Test Nickname",
		City:     "Test City",
		Coordinates: dto.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}

	reqBody, _ := json.Marshal(createLocation)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/locations", bytes.NewBuffer(reqBody)).WithContext(ctx)
	w := httptest.NewRecorder()

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/locations", handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response dto.LocationRes
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
}

func TestGetLocation(t *testing.T) {
	mockSvc := NewMockLocationService()
	handler := GetLocation(mockSvc, slog.Default())
	ctx := context.WithValue(context.Background(), "userId", "1")

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/locations/", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("invalid id")) {
			t.Errorf("Expected response body to contain 'invalid id', got %s", w.Body.String())
		}
	})

	t.Run("location not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/locations/notfound", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/{id}", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("location not found")) {
			t.Errorf("Expected response body to contain 'location not found', got %s", w.Body.String())
		}
	})

	t.Run("successful retrieval", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodGet, "/api/v1/locations/o", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/{id}", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response dto.LocationRes
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

	})
}

func TestGetAllLocations(t *testing.T) {
	mockSvc := NewMockLocationService()
	handler := GetAllLocations(mockSvc, slog.Default())
	ctx := context.WithValue(context.Background(), "userId", "1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/locations?page=1&pageSize=10", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/locations", handler)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Locations []dto.LocationRes `json:"locations"`
		Meta      domain.Metadata   `json:"meta"`
	}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
}

func TestUpdateLocation(t *testing.T) {
	mockSvc := NewMockLocationService()
	handler := UpdateLocation(mockSvc, slog.Default(), webutils.NewCustomValidator(validator.New()))
	ctx := context.WithValue(context.Background(), "userId", "1")

	t.Run("invalid input format", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`invalid json`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/locations/1", reqBody).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/{id}", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("Invalid input format")) {
			t.Errorf("Expected response body to contain 'Invalid input format', got %s", w.Body.String())
		}
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"notes": "", "nickname": "test nickname", "city": "test city", "coordinates": {"lat": 1.0, "lon": 1.0}}`)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/locations/1", reqBody).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/1", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, w.Code)
		}
	})

	t.Run("location not found", func(t *testing.T) {
		updateLocation := dto.LocationReq{
			Notes:    "Updated Notes",
			Nickname: "Updated Nickname",
			City:     "Updated City",
			Coordinates: dto.Coordinates{
				Lat: 2.0,
				Lon: 2.0,
			},
		}

		reqBody, _ := json.Marshal(updateLocation)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/locations/notfound", bytes.NewBuffer(reqBody)).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/{id}", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("location not found")) {
			t.Errorf("Expected response body to contain 'location not found', got %s", w.Body.String())
		}
	})

	t.Run("successful update", func(t *testing.T) {
		updateLocation := dto.LocationReq{
			UserID:   "1",
			Notes:    "Updated Notes",
			Nickname: "Updated Nickname",
			City:     "Updated City",
			Coordinates: dto.Coordinates{
				Lat: 2.0,
				Lon: 2.0,
			},
		}

		reqBody, _ := json.Marshal(updateLocation)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/locations/12", bytes.NewBuffer(reqBody)).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/{id}", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response dto.LocationRes
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
	})
}

func TestDeleteLocation(t *testing.T) {
	mockSvc := NewMockLocationService()
	handler := DeleteLocation(mockSvc, slog.Default())
	ctx := context.WithValue(context.Background(), "userId", "1")

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/locations/", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/{id}", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}

	})

	t.Run("non-existent location", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/locations/notfound", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/{id}", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("location not found")) {
			t.Errorf("Expected response body to contain 'location not found', got %s", w.Body.String())
		}
	})

	t.Run("successful deletion", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/locations/s", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		router := http.NewServeMux()
		router.HandleFunc("/api/v1/locations/{id}", handler)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("location deleted successfully")) {
			t.Errorf("Expected response body to contain 'location deleted successfully', got %s", w.Body.String())
		}
	})
}
