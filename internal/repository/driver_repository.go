package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"gorm.io/gorm"
)

var (
	ErrDriverNotFound = errors.New("driver not found")
)

// DriverRepository is a repository for the driver model
type DriverRepository struct {
	db *gorm.DB
}

// NewDriverRepository creates a new driver repository
func NewDriverRepository(db *gorm.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

func (r *DriverRepository) GetByUserID(ctx context.Context, userID string) (*model.Driver, error) {
	var driver model.Driver

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&driver).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDriverNotFound
		}
		return nil, fmt.Errorf("get driver by user id: %w", err)
	}

	return &driver, nil
}

// Update status offline/online
// UpdateOnline updates the driver's online/offline status.
func (r *DriverRepository) UpdateOnline(ctx context.Context, userID string, isOnline bool) error {
	result := r.db.WithContext(ctx).
		Model(&model.Driver{}).
		Where("user_id = ?", userID).
		Update("is_online", isOnline)

	if result.Error != nil {
		// check if the error is not found
		return fmt.Errorf("update driver online: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrDriverNotFound
	}

	return nil
}

// Update location
func (r *DriverRepository) UpdateLocation(ctx context.Context, userID string, latitude float64, longitude float64) error {
	result := r.db.WithContext(ctx).
		Model(&model.Driver{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"lat":              latitude,
			"lng":              longitude,
			"last_location_at": time.Now(),
		})

	if result.Error != nil {
		// check if the error is not found
		return fmt.Errorf("update driver location: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrDriverNotFound
	}
	return nil
}
