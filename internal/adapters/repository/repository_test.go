package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lafetz/weavo/internal/core/domain"
	"github.com/lafetz/weavo/internal/core/service/location"
)

func TestCreateLocation(t *testing.T) {
	repo := NewInMemoryLocationRepo(24 * time.Hour)
	location := domain.Location{
		Id:       "1",
		UserID:   "user1",
		Notes:    "Test notes",
		Nickname: "Home",
		City:     "City1",
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}

	createdLocation, err := repo.CreateLocation(context.Background(), location)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdLocation.City != location.City {
		t.Fatalf("expected location %+v, got %+v", location, createdLocation)
	}
}

func TestGetLocation(t *testing.T) {
	repo := NewInMemoryLocationRepo(24 * time.Hour)
	location := domain.Location{
		Id:       "1",
		UserID:   "user1",
		Notes:    "Test notes",
		Nickname: "Home",
		City:     "City1",
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}
	loc, err := repo.CreateLocation(context.Background(), location)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedLocation, err := repo.GetLocation(context.Background(), loc.Id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedLocation.City != location.City {
		t.Fatalf("expected location %+v, got %+v", location, retrievedLocation)
	}
}

func TestGetLocations(t *testing.T) {
	repo := NewInMemoryLocationRepo(24 * time.Hour)
	location1 := domain.Location{
		Id:       "1",
		UserID:   "user1",
		Notes:    "Test notes 1",
		Nickname: "Home",
		City:     "City1",
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}
	location2 := domain.Location{
		Id:       "2",
		UserID:   "user1",
		Notes:    "Test notes 2",
		Nickname: "Work",
		City:     "City2",
		Coordinates: domain.Coordinates{
			Lat: 2.0,
			Lon: 2.0,
		},
	}
	repo.CreateLocation(context.Background(), location1)
	repo.CreateLocation(context.Background(), location2)

	filter := location.Filter{PageSize: 10, Page: 1}
	locations, metadata, err := repo.GetLocations(context.Background(), "user1", filter)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(locations) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(locations))
	}

	if metadata.TotalRecords != 2 {
		t.Fatalf("expected total records 2, got %d", metadata.TotalRecords)
	}
}

func TestUpdateLocation(t *testing.T) {
	repo := NewInMemoryLocationRepo(24 * time.Hour)
	location := domain.Location{

		UserID:   "user1",
		Notes:    "Test notes",
		Nickname: "Home",
		City:     "City1",
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}
	loc, err := repo.CreateLocation(context.Background(), location)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedLocation := domain.Location{
		Id:       loc.Id,
		UserID:   "user1",
		Notes:    "Updated notes",
		Nickname: "Home",
		City:     "City1",
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}
	_, err = repo.UpdateLocation(context.Background(), updatedLocation)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedLocation, err := repo.GetLocation(context.Background(), loc.Id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedLocation.Notes != "Updated notes" {
		t.Fatalf("expected updated notes, got %v", retrievedLocation.Notes)
	}
}

func TestDeleteLocation(t *testing.T) {
	repo := NewInMemoryLocationRepo(24 * time.Hour)
	loc := domain.Location{
		Id:       "1",
		UserID:   "user1",
		Notes:    "Test notes",
		Nickname: "Home",
		City:     "City1",
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}
	loc, err := repo.CreateLocation(context.Background(), loc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = repo.DeleteLocation(context.Background(), loc.Id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.GetLocation(context.Background(), loc.Id)
	if !errors.Is(err, location.ErrLocationNotFound) {
		t.Fatalf("expected error %v, got %v", location.ErrLocationNotFound, err)
	}
}
func TestCleanupExpiredLocations(t *testing.T) {
	repo := NewInMemoryLocationRepo(1 * time.Second)
	loc := domain.Location{
		Id:        "1",
		UserID:    "user1",
		Notes:     "Test notes",
		Nickname:  "Home",
		City:      "City1",
		CreatedAt: time.Now().Add(-25 * time.Hour), // Set CreatedAt to more than 24 hours ago
		Coordinates: domain.Coordinates{
			Lat: 1.0,
			Lon: 1.0,
		},
	}
	repo.CreateLocation(context.Background(), loc)

	dataRetention := 1 * time.Second
	time.Sleep(dataRetention + 1*time.Microsecond)

	_, err := repo.GetLocation(context.Background(), "1")
	if !errors.Is(err, location.ErrLocationNotFound) {
		t.Fatalf("expected error %v, got %v", location.ErrLocationNotFound, err)
	}
}
