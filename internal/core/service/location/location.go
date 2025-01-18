package location

import (
	"context"
	"errors"

	"github.com/lafetz/weavo/internal/core/domain"
)

var (
	ErrLocationNotFound = errors.New("location not found")
	ErrUnAuthorized     = errors.New("unauthorized")
)

type Service struct {
	repo LocationRepo
}

func NewService(repo LocationRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLocation(ctx context.Context, location domain.Location) (domain.Location, error) {
	return s.repo.CreateLocation(ctx, location)
}

func (s *Service) GetLocation(ctx context.Context, id string) (domain.Location, error) {
	return s.repo.GetLocation(ctx, id)
}

func (s *Service) GetLocations(ctx context.Context, userID string, filter Filter) ([]domain.Location, domain.Metadata, error) {
	return s.repo.GetLocations(ctx, userID, filter)
}

func (s *Service) UpdateLocation(ctx context.Context, location domain.Location) (domain.Location, error) {
	loc, err := s.repo.GetLocation(ctx, location.Id)
	if err != nil {
		return domain.Location{}, ErrLocationNotFound
	}
	if loc.UserID != location.UserID {
		return domain.Location{}, ErrUnAuthorized
	}
	return s.repo.UpdateLocation(ctx, location)
}

func (s *Service) DeleteLocation(ctx context.Context, id string, userID string) error {
	loc, err := s.repo.GetLocation(ctx, id)
	if err != nil {
		return ErrLocationNotFound
	}
	if loc.UserID != userID {
		return ErrUnAuthorized
	}

	return s.repo.DeleteLocation(ctx, id)
}
