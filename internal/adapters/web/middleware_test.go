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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
