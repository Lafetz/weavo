package dto

import "github.com/lafetz/weavo/internal/core/domain"

type Coordinates struct {
	Lat float64 `json:"lat" validate:"required"`
	Lon float64 `json:"lon" validate:"required"`
}

// request
type LocationReq struct {
	UserID      string
	Notes       string      `json:"notes" validate:"required,min=1"`
	Nickname    string      `json:"nickname" validate:"required,min=1"`
	City        string      `json:"city" validate:"required,min=1"`
	Coordinates Coordinates `json:"coordinates" validate:"required"`
}

func (l *LocationReq) ToDomain() domain.Location {
	return domain.Location{
		UserID:   l.UserID,
		Notes:    l.Notes,
		Nickname: l.Nickname,
		City:     l.City,
		Coordinates: domain.Coordinates{
			Lat: l.Coordinates.Lat,
			Lon: l.Coordinates.Lon,
		},
	}
}

// response
// get location
type LocationRes struct {
	Id          string      `json:"id"`
	Notes       string      `json:"notes"`
	Nickname    string      `json:"nickname"`
	City        string      `json:"city"`
	Coordinates Coordinates `json:"coordinates"`
	CreatedAt   string      `json:"created_at"`
}

func GetLocationRes(l domain.Location) LocationRes {
	return LocationRes{
		Id:       l.Id,
		Notes:    l.Notes,
		Nickname: l.Nickname,
		City:     l.City,
		Coordinates: Coordinates{
			Lat: l.Coordinates.Lat,
			Lon: l.Coordinates.Lon,
		},
		CreatedAt: l.CreatedAt.String(),
	}
}

// get locations
type JSONMetadata struct {
	CurrentPage  int32 `json:"currentPage"`
	PageSize     int32 `json:"pageSize"`
	FirstPage    int32 `json:"firstPage"`
	LastPage     int32 `json:"lastPage"`
	TotalRecords int32 `json:"totalRecords"`
}

func ConvertToJSONMetadata(meta domain.Metadata) JSONMetadata {
	return JSONMetadata{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		FirstPage:    meta.FirstPage,
		LastPage:     meta.LastPage,
		TotalRecords: meta.TotalRecords,
	}
}

type GetLocationssResponse struct {
	Meta      JSONMetadata  `json:"meta"`
	Locations []LocationRes `json:"locations"`
}

func GetLocationsRes(locations []domain.Location, meta domain.Metadata) GetLocationssResponse {
	var locationRes []LocationRes
	for _, l := range locations {
		locationRes = append(locationRes, GetLocationRes(l))
	}
	return GetLocationssResponse{
		Meta:      ConvertToJSONMetadata(meta),
		Locations: locationRes,
	}
}
