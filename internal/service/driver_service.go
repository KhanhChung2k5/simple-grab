package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/repository"
)

var ErrDriverNotFound = errors.New("driver not found")

// DriverService is the service for the driver
type DriverService struct {
	drivers *repository.DriverRepository
}

// NewDriverService creates a new driver service
func NewDriverService(drivers *repository.DriverRepository) *DriverService {
	return &DriverService{drivers: drivers}
}

// GetByUserID gets a driver by user id
func (s *DriverService) GetByUserID(ctx context.Context, userID string) (*model.Driver, error) {
	driver, err := s.drivers.GetByUserID(ctx, userID)
	if err != nil {
		return nil, mapDriverError(err)
	}
	return driver, nil
}

// SetOnline updates the driver's online/offline status
func (s *DriverService) SetOnline(ctx context.Context, userID string, isOnline bool) error {
	return mapDriverError(s.drivers.UpdateOnline(ctx, userID, isOnline))
}

// UpdateLocation updates the driver's current location
func (s *DriverService) UpdateLocation(ctx context.Context, userID string, latitude, longitude float64) error {
	return mapDriverError(s.drivers.UpdateLocation(ctx, userID, latitude, longitude))
}

func mapDriverError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrDriverNotFound) {
		return ErrDriverNotFound
	}
	return fmt.Errorf("driver repository: %w", err)
}
