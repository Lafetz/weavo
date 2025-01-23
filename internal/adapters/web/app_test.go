package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"

	mockcache "github.com/lafetz/weavo/internal/adapters/mock_cache"
	"github.com/lafetz/weavo/internal/adapters/repository"

	"github.com/lafetz/weavo/internal/adapters/web/dto"
	"github.com/lafetz/weavo/internal/adapters/web/webutils"
	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/location"
	"github.com/lafetz/weavo/internal/core/service/weather"
	customlogger "github.com/lafetz/weavo/internal/logger"
)

const dataRetention = 24 * time.Second

var locationID = "some-valid-id"

type opmock struct {
}

func (o *opmock) GetWeather(ctx context.Context, city string) (domain.Weather, error) {
	return domain.Weather{}, nil
}
func setupServer() *App {
	ow := &opmock{}
	logger := customlogger.NewLogger(slog.LevelDebug, "development")
	store := repository.NewInMemoryLocationRepo(dataRetention)
	locationID = seedDatabase(store)
	locationSvc := location.NewService(store)
	mc := mockcache.NewMockCache()
	weatherSvc := weather.NewService(ow, mc)
	val := validator.New()
	custonmVal := webutils.NewCustomValidator(val)
	cookieStore := webutils.CookieStore(dataRetention)
	app := NewApp(8080, logger, cookieStore, custonmVal, locationSvc, weatherSvc)

	return app
}

func seedDatabase(store *repository.InMemoryLocationRepo) string {

	loc, _ := store.CreateLocation(context.TODO(), domain.Location{

		UserID:      "test-user-id",
		Notes:       "Test Notes",
		Nickname:    "Test Nickname",
		City:        "Test City",
		Coordinates: domain.Coordinates{Lat: 1.0, Lon: 1.0},
		CreatedAt:   time.Now(),
	})
	return loc.Id
}

func TestCreateLocation(t *testing.T) {
	app := setupServer()
	server := httptest.NewServer(app.Router)
	defer server.Close()

	t.Run("invalid input format", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`invalid json`)
		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/locations", reqBody)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
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
		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/locations", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
		}
	})

	t.Run("successful creation", func(t *testing.T) {
		createLocation := dto.LocationReq{
			Notes:    "Test Notes",
			Nickname: "Test Nickname",
			City:     "Test City",
			Coordinates: dto.Coordinates{
				Lat: 1.0,
				Lon: 1.0,
			},
		}

		reqBody, _ := json.Marshal(createLocation)
		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/locations", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
		}
	})
}

func TestGetLocation(t *testing.T) {
	app := setupServer()
	server := httptest.NewServer(app.Router)
	defer server.Close()

	t.Run("invalid id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/locations/ ", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("location not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/locations/notfound", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("successful retrieval", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/locations/"+locationID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}

func TestGetAllLocations(t *testing.T) {
	app := setupServer()
	server := httptest.NewServer(app.Router)
	defer server.Close()

	// Create a request with a valid session cookie
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/v1/locations?page=1&pageSize=10", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	addcookie(app, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Println("resp", string(bodyBytes))
	var response struct {
		Status   string            `json:"status"`
		Message  string            `json:"message"`
		Data     []dto.LocationRes `json:"data"`
		Metadata struct {
			CurrentPage  int `json:"currentPage"`
			PageSize     int `json:"pageSize"`
			FirstPage    int `json:"firstPage"`
			LastPage     int `json:"lastPage"`
			TotalRecords int `json:"totalRecords"`
		} `json:"metadata"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expectedCount := 1
	if len(response.Data) != expectedCount {
		t.Errorf("Expected %d locations, got %d", expectedCount, len(response.Data))
	}
}

func TestUpdateLocation(t *testing.T) {
	app := setupServer()
	server := httptest.NewServer(app.Router)
	defer server.Close()

	t.Run("invalid input format", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`invalid json`)
		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/locations/1", reqBody)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := bytes.NewBufferString(`{"notes": "", "nickname": "test nickname", "city": "test city", "coordinates": {"lat": 1.0, "lon": 1.0}}`)
		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/locations/1", reqBody)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
		}
	})

	t.Run("location not found", func(t *testing.T) {
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
		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/locations/1", bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
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
		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/v1/locations/"+locationID, bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}

func TestDeleteLocation(t *testing.T) {
	app := setupServer()
	server := httptest.NewServer(app.Router)
	defer server.Close()

	t.Run("invalid id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/locations/  ", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("successful deletion", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/v1/locations/"+locationID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		addcookie(app, req)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})
}
func addcookie(app *App, req *http.Request) {
	session, _ := app.store.Get(req, "user-session")
	userId := "test-user-id"
	session.Values["userId"] = userId
	rec := httptest.NewRecorder()
	session.Save(req, rec)

	for _, cookie := range rec.Result().Cookies() {
		req.AddCookie(cookie)
	}
}
