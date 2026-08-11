package service

import (
	"errors"

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
func (s *DriverService) GetByUserID(userID string) (*model.Driver, error) {
	return s.drivers.GetByUserId(userID)
}

// SetOnline updates the driver's online/offline status
func (s *DriverService) SetOnline(userID string, isOnline bool) error {
	return s.drivers.UpdateOnline(userID, isOnline)
}

// UpdateLocation updates the driver's current location
func (s *DriverService) UpdateLocation(userID string, latitude, longitude float64) error {
	return s.drivers.UpdateLocation(userID, latitude, longitude)
}
