package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/location"
)

type InMemoryLocationRepo struct {
	mu            sync.RWMutex
	locations     map[string]domain.Location
	dataRetention time.Duration
}

func NewInMemoryLocationRepo(dataRetention time.Duration) *InMemoryLocationRepo {
	repo := &InMemoryLocationRepo{
		locations:     make(map[string]domain.Location),
		dataRetention: dataRetention,
	}
	go repo.cleanupExpiredLocations()
	return repo
}

func (repo *InMemoryLocationRepo) CreateLocation(ctx context.Context, loc domain.Location) (domain.Location, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	loc.Id = uuid.New().String()
	loc.CreatedAt = time.Now()
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
	locations := []domain.Location{}
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
	metadata := domain.CalculateMetadata(int32(totalRecords), int32(filter.Page), int32(filter.PageSize))

	return paginatedLocations, metadata, nil
}

func (repo *InMemoryLocationRepo) UpdateLocation(ctx context.Context, loc domain.Location) (domain.Location, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	el, exists := repo.locations[loc.Id]
	if !exists {
		return domain.Location{}, location.ErrLocationNotFound
	}
	el.Nickname = loc.Nickname
	el.Notes = loc.Notes
	loc.CreatedAt = el.CreatedAt
	repo.locations[loc.Id] = el
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

func (repo *InMemoryLocationRepo) cleanupExpiredLocations() {
	ticker := time.NewTicker(repo.dataRetention)
	for {
		<-ticker.C
		repo.mu.Lock()
		for id, loc := range repo.locations {
			if time.Since(loc.CreatedAt) > repo.dataRetention {
				delete(repo.locations, id)
			}
		}
		repo.mu.Unlock()
	}
}
