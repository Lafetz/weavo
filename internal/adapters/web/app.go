package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/lafetz/weavo/internal/adapters/web/webutils"
	"github.com/lafetz/weavo/internal/core/service/location"
	"github.com/lafetz/weavo/internal/core/service/weather"
)

type App struct {
	port        int
	Router      *http.ServeMux
	logger      *slog.Logger
	validator   *webutils.CustomValidator
	locationSvc location.ServiceApi
	weatherSvc  weather.ServiceApi
	store       *sessions.CookieStore
}

func NewApp(
	port int, logger *slog.Logger,
	store *sessions.CookieStore,
	validator *webutils.CustomValidator,
	locationSvc *location.Service,
	weatherSvc weather.ServiceApi,
) *App {
	a := &App{
		Router:      http.NewServeMux(),
		logger:      logger,
		port:        port,
		validator:   validator,
		locationSvc: locationSvc,
		weatherSvc:  weatherSvc,
		store:       store,
	}
	a.initAppRoutes()
	return a
}
func (a *App) Run() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", strconv.Itoa(a.port)),
		Handler:      a.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		a.logger.Info("shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()
	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	a.logger.Info("server stopped")
	return nil
}
