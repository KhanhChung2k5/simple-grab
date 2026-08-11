package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"gorm.io/gorm"
)

var (
	ErrDriverNotFound     = errors.New("driver not found")
	ErrDriverAlreadyTaken = errors.New("driver already has a ride")
)

// DriverRepository is a repository for the driver model
type DriverRepository struct {
	db *gorm.DB
}

// NewDriverRepository creates a new driver repository
func NewDriverRepository(db *gorm.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

func (r *DriverRepository) GetByUserId(userId string) (*model.Driver, error) {
	var driver model.Driver
	err := r.db.Where("user_id = ?", userId).First(&driver).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("driver not found")
		}
		return nil, err
	}
	return &driver, nil
}

//Update status offline/online
// UpdateOnline updates the driver's online/offline status.
func (r *DriverRepository) UpdateOnline(userID string, isOnline bool) error {
	result := r.db.
		Model(&model.Driver{}).
		Where("user_id = ?", userID).
		Update("is_online", isOnline)

	if result.Error != nil {
		return fmt.Errorf(
			"failed to update online status: %w",
			result.Error,
		)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("driver not found")
	}

	return nil
}

//Update location
func (r *DriverRepository) UpdateLocation(userID string, latitude float64, longitude float64) error {
	result := r.db.
		Model(&model.Driver{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"lat": latitude,
			"lng": longitude,
			"last_location_at": time.Now(),
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update location: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("driver not found")
	}
	return nil
}
