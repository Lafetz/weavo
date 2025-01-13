package location

import (
	"context"

	"github.com/lafetz/weavo/internal/core/domain"
)

type Filter struct {
	PageSize int
	Page     int
}
type LocationRepo interface {
	CreateLocation(ctx context.Context, location domain.Location) (domain.Location, error)
	GetLocation(ctx context.Context, id string) (domain.Location, error)
	GetLocations(ctx context.Context, userID string, filter Filter) ([]domain.Location, domain.Metadata, error)
	UpdateLocation(ctx context.Context, location domain.Location) (domain.Location, error)
	DeleteLocation(ctx context.Context, id string) error
}
