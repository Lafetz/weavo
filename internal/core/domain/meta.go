package domain

import "math"

type Metadata struct {
	CurrentPage  int32
	PageSize     int32
	FirstPage    int32
	LastPage     int32
	TotalRecords int32
}

func CalculateMetadata(totalRecords int32, offset, limit int32) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  (offset / limit) + 1,
		PageSize:     limit,
		FirstPage:    1,
		LastPage:     int32(math.Ceil(float64(totalRecords) / float64(limit))),
		TotalRecords: totalRecords,
	}
}
