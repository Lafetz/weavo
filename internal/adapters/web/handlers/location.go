package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/lafetz/weavo/internal/adapters/web/dto"
	"github.com/lafetz/weavo/internal/adapters/web/webutils"
	"github.com/lafetz/weavo/internal/core/service/location"
)

const hello = ""

// CreateLocation handles the creation of a new location.
//
// @Summary Create a new location
// @Description Create a new location with the provided details
// @Tags locations
// @Accept json
// @Produce json
// @Param location body dto.LocationReq true "Location request body"
// @Success 201 {object} dto.LocationRes "Location created successfully"
// @Failure 400 {string} string "Invalid input format"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/locations [post]
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

// GetLocation handles the HTTP request to retrieve a location by its ID.
//
// @Summary Retrieve a location by ID
// @Description Retrieves a location from the service using the provided ID.
// @Tags locations
// @Accept json
// @Produce json
// @Param id path string true "Location ID"
// @Success 200 {object} dto.LocationRes "location retrieved successfully"
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "location not found"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/locations/{id} [get]
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

// UpdateLocation handles the HTTP request for updating a location.
// @Summary Update a location
// @Description Update an existing location with the provided details
// @Tags locations
// @Accept json
// @Produce json
// @Param id path string true "Location ID"
// @Param LocationReq body dto.LocationReq true "Location request body"
// @Success 200 {object} dto.LocationRes "location updated successfully"
// @Failure 400 {string} string "Invalid input format or invalid id"
// @Failure 404 {string} string "location not found"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/locations/{id} [put]
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

// DeleteLocation handles the HTTP request for deleting a location.
//
// @Summary Delete a location
// @Description Deletes a location by its ID
// @Tags locations
// @Accept json
// @Produce json
// @Param id path string true "Location ID"
// @Success 200 {string} string "location deleted successfully"
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "location not found"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/locations/{id} [delete]
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

// GetAllLocations handles the HTTP request to retrieve all locations for a user.
//
// @Summary Retrieve all locations
// @Description Retrieves a list of locations for the authenticated user with pagination support.
// @Tags locations
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(5)
// @Success 200 {object} dto.LocationsRes "locations retrieved successfully"
// @Failure 500 {object} webutils.ErrorResponse "internal server error"
// @Router /api/v1/locations [get]

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
