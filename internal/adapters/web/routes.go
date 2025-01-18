package web

import (
	_ "github.com/lafetz/weavo/docs"
	"github.com/lafetz/weavo/internal/adapters/web/handlers"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (a *App) initAppRoutes() {
	a.Router.HandleFunc("/api/swagger/", httpSwagger.WrapHandler)
	a.Router.HandleFunc("GET /api/v1/locations", a.recoverPanic(a.UserContext(handlers.GetAllLocations(a.locationSvc, a.logger))))
	a.Router.HandleFunc("GET /api/v1/locations/{id}", a.recoverPanic(a.UserContext(handlers.GetLocation(a.locationSvc, a.logger))))
	a.Router.HandleFunc("POST /api/v1/locations", a.recoverPanic(a.UserContext(handlers.CreateLocation(a.locationSvc, a.logger, a.validator))))
	a.Router.HandleFunc("PUT /api/v1/locations/{id}", a.recoverPanic(a.UserContext(handlers.UpdateLocation(a.locationSvc, a.logger, a.validator))))
	a.Router.HandleFunc("DELETE /api/v1/locations/{id}", a.recoverPanic(a.UserContext(handlers.DeleteLocation(a.locationSvc, a.logger))))
	a.Router.HandleFunc("GET /api/v1/weather", a.recoverPanic(a.UserContext(handlers.GetWeather(a.weatherSvc, a.logger))))

}
