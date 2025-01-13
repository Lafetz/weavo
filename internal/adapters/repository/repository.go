package repository

import (
	"context"
	"sync"

	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/location"
)

type InMemoryLocationRepo struct {
	mu        sync.RWMutex
	locations map[string]domain.Location
}

func NewInMemoryLocationRepo() *InMemoryLocationRepo {
	return &InMemoryLocationRepo{
		locations: make(map[string]domain.Location),
	}
}

func (repo *InMemoryLocationRepo) CreateLocation(ctx context.Context, loc domain.Location) (domain.Location, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.locations[loc.Id] = loc
	return loc, nil
}

func (repo *InMemoryLocationRepo) GetLocation(ctx context.Context, id string) (domain.Location, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	loc, exists := repo.locations[id]
	if !exists {
		return domain.Location{}, location.ErrLocationNotFound
	}
	return loc, nil
}

func (repo *InMemoryLocationRepo) GetLocations(ctx context.Context, userID string, filter location.Filter) ([]domain.Location, domain.Metadata, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	var locations []domain.Location
	for _, loc := range repo.locations {
		if loc.UserID == userID {
			locations = append(locations, loc)
		}
	}

	totalRecords := len(locations)
	start := filter.PageSize * (filter.Page - 1)
	end := start + filter.PageSize
	if start > totalRecords {
		start = totalRecords
	}
	if end > totalRecords {
		end = totalRecords
	}

	paginatedLocations := locations[start:end]
	metadata := domain.CalculateMetadata(int32(totalRecords), int32(start), int32(filter.PageSize))

	return paginatedLocations, metadata, nil
}

func (repo *InMemoryLocationRepo) UpdateLocation(ctx context.Context, loc domain.Location) (domain.Location, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, exists := repo.locations[loc.Id]
	if !exists {
		return domain.Location{}, location.ErrLocationNotFound
	}
	repo.locations[loc.Id] = loc
	return loc, nil
}

func (repo *InMemoryLocationRepo) DeleteLocation(ctx context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, exists := repo.locations[id]
	if !exists {
		return location.ErrLocationNotFound
	}
	delete(repo.locations, id)
	return nil
}
