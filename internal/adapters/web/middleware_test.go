package web

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func TestRecoverPanic(t *testing.T) {
	app := &App{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	handler := app.recoverPanic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
	if !contains(w.Body.String(), "internal server error") {
		t.Errorf("Expected response body to contain 'internal server error', got %s", w.Body.String())
	}
}

func TestUserContext(t *testing.T) {
	store := sessions.NewCookieStore([]byte("secret"))
	app := &App{
		store:  store,
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	handler := app.UserContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)
		if userId == "" {
			t.Errorf("Expected userId to be set in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("new user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		session, _ := store.Get(req, "user-session")
		userId, exists := session.Values["userId"].(string)
		if !exists || userId == "" {
			t.Errorf("Expected userId to be set in session")
		}
	})

	t.Run("existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		session, _ := store.Get(req, "user-session")
		userId := uuid.New().String()
		session.Values["userId"] = userId
		session.Save(req, w)

		req = httptest.NewRequest(http.MethodGet, "/", nil)
		for _, cookie := range w.Result().Cookies() {
			req.AddCookie(cookie)
		}
		w = httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		session, _ = store.Get(req, "user-session")
		if session.Values["userId"] != userId {
			t.Errorf("Expected userId %s, got %s", userId, session.Values["userId"])
		}
	})
}
func TestEnableCORS(t *testing.T) {
	app := &App{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	handler := app.enableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("CORS headers set", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.Header.Get("Access-Control-Allow-Origin") != "http://example.com" {
			t.Errorf("Expected Access-Control-Allow-Origin header to be set")
		}
		if resp.Header.Get("Vary") != "Origin" {
			t.Errorf("Expected Vary header to be set to 'Origin'")
		}
		if resp.Header.Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
			t.Errorf("Expected Access-Control-Allow-Methods header to be set")
		}
		if resp.Header.Get("Access-Control-Allow-Headers") != "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With" {
			t.Errorf("Expected Access-Control-Allow-Headers header to be set")
		}
		if resp.Header.Get("Access-Control-Max-Age") != "3600" {
			t.Errorf("Expected Access-Control-Max-Age header to be set")
		}
	})

	t.Run("OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", "http://example.com")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
