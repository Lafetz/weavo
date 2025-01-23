package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (app *App) recoverPanic(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection:", "close")
				var errorMessage string
				if e, ok := err.(error); ok {
					errorMessage = e.Error()
					app.logger.Error(errorMessage)

				} else {
					errorMessage = fmt.Sprintf("panic: %v", err)
					app.logger.Error(errorMessage)
				}
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	}
}
func (app *App) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.store.Get(r, "user-session")
		if err != nil {
			if !strings.Contains(err.Error(), "the value is not valid") {
				app.logger.Error("error getting session", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

		}
		userId, exists := session.Values["userId"].(string)
		if !exists || userId == "" {
			userId = uuid.New().String()
			session.Values["userId"] = userId
			if err := session.Save(r, w); err != nil {
				app.logger.Error("error saving session", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		fmt.Println("userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func (app *App) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
