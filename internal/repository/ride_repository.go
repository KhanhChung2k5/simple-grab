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
	//
	// ErrRideNotFound is returned when a ride is not found
	ErrRideNotFound     = errors.New("ride not found")
	// ErrRideAlreadyTaken is returned when a ride already has a driver
	ErrRideAlreadyTaken = errors.New("ride already has a driver")
)

type RideRepository struct {
	db *gorm.DB
}

func NewRideRepository(db *gorm.DB) *RideRepository {
	return &RideRepository{db: db}
}

// Create creates a new ride
func (r *RideRepository) Create(ctx context.Context, ride *model.Ride) (*model.Ride, error) {
	ride.Status = model.RidePending

	err := r.db.WithContext(ctx).Create(ride).Error
	if err != nil {
		return nil, fmt.Errorf("create ride: %w", err)
	}

	return ride, nil
}

// GetByID gets a ride by its ID
func (r *RideRepository) GetByID(ctx context.Context, id string) (*model.Ride, error) {
	var ride model.Ride
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ride).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRideNotFound
		}
		return nil, fmt.Errorf("get ride by id: %w", err)
	}

	return &ride, nil
}

// ListByRider lists all rides by a rider
func (r *RideRepository) ListByRider(ctx context.Context, riderID string) ([]model.Ride, error) {
	var rides []model.Ride
	err := r.db.WithContext(ctx).
		Where("rider_id = ?", riderID).
		Order("created_at DESC").
		Find(&rides).Error
	if err != nil {
		return nil, fmt.Errorf("list rides by rider: %w", err)
	}

	return rides, nil
}

// ListByDriver lists all rides by a driver
func (r *RideRepository) ListByDriver(ctx context.Context, driverID string) ([]model.Ride, error) {
	var rides []model.Ride
	err := r.db.WithContext(ctx).
		Where("driver_id = ?", driverID).
		Order("created_at DESC").
		Find(&rides).Error
	if err != nil {
		return nil, fmt.Errorf("list rides by driver: %w", err)
	}

	return rides, nil
}

// ListAvailable lists all available rides
func (r *RideRepository) ListAvailable(ctx context.Context) ([]model.Ride, error) {
	var rides []model.Ride
	err := r.db.WithContext(ctx).
		Where("status = ? AND driver_id IS NULL", model.RidePending).
		Order("created_at ASC").
		Find(&rides).Error
	if err != nil {
		return nil, fmt.Errorf("list available rides: %w", err)
	}

	return rides, nil
}

// Accept accepts a ride by a driver
func (r *RideRepository) Accept(ctx context.Context, rideID, driverID string) (*model.Ride, error) {
	result := r.db.WithContext(ctx).
		Model(&model.Ride{}).
		Where("id = ? AND status = ? AND driver_id IS NULL", rideID, model.RidePending).
		Updates(map[string]interface{}{
			"driver_id":  driverID,
			"status":     model.RideAccepted,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return nil, fmt.Errorf("accept ride: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ErrRideAlreadyTaken
	}

	return r.GetByID(ctx, rideID)
}

// UpdateStatus updates the status of a ride
func (r *RideRepository) UpdateStatus(ctx context.Context, rideID, driverID, newStatus string) (*model.Ride, error) {
	result := r.db.WithContext(ctx).
		Model(&model.Ride{}).
		Where("id = ? AND driver_id = ?", rideID, driverID).
		Updates(map[string]interface{}{
			"status":     newStatus,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return nil, fmt.Errorf("update ride status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ErrRideNotFound
	}

	return r.GetByID(ctx, rideID)
}
