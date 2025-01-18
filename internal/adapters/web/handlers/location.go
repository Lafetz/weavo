package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/lafetz/weavo/internal/adapters/web/dto"
	"github.com/lafetz/weavo/internal/adapters/web/webutils"
	"github.com/lafetz/weavo/internal/core/service/location"
)

func CreateLocation(locationSvc location.ServiceApi, logger *slog.Logger, validator *webutils.CustomValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LocationReq
		if err := webutils.ReadJSON(w, r, &req); err != nil {
			webutils.WriteJSON(w, http.StatusBadRequest, "Invalid input format", nil, nil)
			logger.Error(err.Error())
			return
		}
		if validator.ValidateAndRespond(w, req) {
			return
		}

		location := req.ToDomain()
		location.UserID = r.Context().Value("userId").(string)
		loc, err := locationSvc.CreateLocation(r.Context(), location)

		if err != nil {
			webutils.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil, nil)
			logger.Error("error on creating location", "error", err.Error())
			return
		}

		webutils.WriteJSON(w, http.StatusCreated, "location created successfully", dto.GetLocationRes(loc), nil)
	}
}

func GetLocation(locationSvc location.ServiceApi, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			webutils.WriteJSON(w, http.StatusBadRequest, "invalid id", nil, nil)
			return
		}

		loc, err := locationSvc.GetLocation(r.Context(), id)

		if err != nil {
			if errors.Is(err, location.ErrLocationNotFound) {
				webutils.WriteJSON(w, http.StatusNotFound, "location not found", nil, nil)
				return
			}
			webutils.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil, nil)
			logger.Error("error on getting location by id", "error", err.Error())
			return
		}

		webutils.WriteJSON(w, http.StatusOK, "location retrieved successfully", dto.GetLocationRes(loc), nil)
	}
}

func UpdateLocation(locationSvc location.ServiceApi, logger *slog.Logger, validator *webutils.CustomValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.LocationReq
		if err := webutils.ReadJSON(w, r, &req); err != nil {
			webutils.WriteJSON(w, http.StatusBadRequest, "Invalid input format", nil, nil)
			return
		}
		if validator.ValidateAndRespond(w, req) {
			return
		}
		id := r.PathValue("id")
		if id == "" {
			webutils.WriteJSON(w, http.StatusBadRequest, "invalid id", nil, nil)
			return
		}

		loc := req.ToDomain()
		loc.Id = id
		loc.UserID = r.Context().Value("userId").(string)
		loc, err := locationSvc.UpdateLocation(r.Context(), loc)
		if err != nil {
			if errors.Is(err, location.ErrLocationNotFound) {
				webutils.WriteJSON(w, http.StatusNotFound, "location not found", nil, nil)
				return
			}
			webutils.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil, nil)
			logger.Error("error on updating location", "error", err.Error())
			return
		}

		webutils.WriteJSON(w, http.StatusOK, "location updated successfully", dto.GetLocationRes(loc), nil)
	}
}

func DeleteLocation(locationSvc location.ServiceApi, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			webutils.WriteJSON(w, http.StatusBadRequest, "invalid id", nil, nil)
			return
		}
		userId := r.Context().Value("userId").(string)
		err := locationSvc.DeleteLocation(r.Context(), id, userId)
		if err != nil {
			if errors.Is(err, location.ErrLocationNotFound) {
				webutils.WriteJSON(w, http.StatusNotFound, "location not found", nil, nil)
				return
			}
			webutils.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil, nil)
			logger.Error("error on deleting location", "error", err.Error())
			return
		}

		webutils.WriteJSON(w, http.StatusOK, "location deleted successfully", nil, nil)
	}
}

func GetAllLocations(locationSvc location.ServiceApi, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("userId").(string)

		filters := location.Filter{
			Page:     webutils.GetQueryInt(r, "page", 1),
			PageSize: webutils.GetQueryInt(r, "pageSize", 5),
		}

		locations, metadata, err := locationSvc.GetLocations(r.Context(), userId, filters)
		if err != nil {
			webutils.WriteJSON(w, http.StatusInternalServerError, "internal server error", nil, nil)
			logger.Error("error on getting all locations", "error", err.Error())
			return
		}
		locationsRes := dto.GetLocationsRes(locations, metadata)
		webutils.WriteJSON(w, http.StatusOK, "locations retrieved successfully", locationsRes.Locations, locationsRes.Meta)
	}
}
