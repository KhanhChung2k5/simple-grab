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
	ErrRideNotFound        = errors.New("ride not found")
	ErrRideAlreadyTaken    = errors.New("ride already has a driver")
	ErrRideStatusChanged   = errors.New("ride status has changed")
	ErrRiderHasActiveRide  = errors.New("rider already has an active ride")
	ErrDriverHasActiveRide = errors.New("driver already has an active ride")
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrRiderHasActiveRide
		}
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
	ride, err := r.GetByID(ctx, rideID)
	if err != nil {
		return nil, err
	}

	if ride.Status != model.RidePending || ride.DriverID != nil {
		return nil, ErrRideAlreadyTaken
	}

	result := r.db.WithContext(ctx).
		Model(&model.Ride{}).
		Where(
			"id = ? AND status = ? AND driver_id IS NULL",
			rideID,
			model.RidePending,
		).
		// update the ride status and driver id
		Updates(map[string]any{
			"driver_id":  driverID,
			"status":     model.RideAccepted,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		// check if the error is duplicated key
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, ErrDriverHasActiveRide
		}
		return nil, fmt.Errorf("accept ride: %w", result.Error)
	}

	// check if the ride is not found
	if result.RowsAffected == 0 {
		return nil, ErrRideAlreadyTaken
	}

	// return the ride

	return r.GetByID(ctx, rideID)
}

// UpdateStatus updates the status of a ride
func (r *RideRepository) UpdateStatus(
	ctx context.Context,
	rideID string,
	driverID string,
	currentStatus string,
	newStatus string,
) (*model.Ride, error) {
	// update the ride status
	result := r.db.WithContext(ctx).
		Model(&model.Ride{}).
		Where(
			"id = ? AND driver_id = ? AND status = ?",
			rideID,
			driverID,
			currentStatus,
		).
		// update the ride status and updated at
		Updates(map[string]any{
			"status":     newStatus,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		// check if the error is duplicated key
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, ErrRideStatusChanged
		}
		return nil, fmt.Errorf("update ride status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ErrRideStatusChanged
	}

	return r.GetByID(ctx, rideID)
}

// HasActiveRideByRider reports whether a rider has a pending or ongoing ride.
func (r *RideRepository) HasActiveRideByRider(ctx context.Context, riderID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Ride{}).
		Where("rider_id = ? AND status IN ?", riderID, []string{
			model.RidePending,
			model.RideAccepted,
			model.RideInProgress,
		}).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check rider active ride: %w", err)
	}

	return count > 0, nil
}

// HasActiveRideByDriver reports whether a driver has an accepted or ongoing ride.
func (r *RideRepository) HasActiveRideByDriver(ctx context.Context, driverID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Ride{}).
		Where("driver_id = ? AND status IN ?", driverID, []string{
			model.RideAccepted,
			model.RideInProgress,
		}).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check driver active ride: %w", err)
	}

	return count > 0, nil
}

// CancelByRider atomically cancels a rider-owned ride if its status has not changed.
func (r *RideRepository) CancelByRider(ctx context.Context, rideID, riderID, currentStatus string) (*model.Ride, error) {
	result := r.db.WithContext(ctx).
		Model(&model.Ride{}).
		Where(
			"id = ? AND rider_id = ? AND status = ?",
			rideID,
			riderID,
			currentStatus,
		).
		Updates(map[string]any{
			"status":     model.RideCancelled,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return nil, fmt.Errorf("cancel ride: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ErrRideStatusChanged
	}

	return r.GetByID(ctx, rideID)
}

// CancelByDriver atomically cancels a driver-assigned ride if its status has not changed.
func (r *RideRepository) CancelByDriver(ctx context.Context, rideID, driverID, currentStatus string) (*model.Ride, error) {
	result := r.db.WithContext(ctx).
		Model(&model.Ride{}).
		Where(
			"id = ? AND driver_id = ? AND status = ?",
			rideID,
			driverID,
			currentStatus,
		).
		Updates(map[string]any{
			"status":     model.RideCancelled,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return nil, fmt.Errorf("cancel ride by driver: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ErrRideStatusChanged
	}

	return r.GetByID(ctx, rideID)
}
