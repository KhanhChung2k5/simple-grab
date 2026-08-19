package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/repository"
)

// error constants
var (
	ErrRideNotFound        = errors.New("ride not found")
	ErrRideAlreadyTaken    = errors.New("ride already taken")
	ErrRideStatusInvalid   = errors.New("ride status invalid")
	ErrForbidden           = errors.New("forbidden")
	ErrDriverOffline       = errors.New("driver must be online to accept rides")
	ErrRiderHasActiveRide  = errors.New("rider already has an active ride")
	ErrDriverHasActiveRide = errors.New("driver already has an active ride")
	ErrRideCannotCancel    = errors.New("ride cannot be cancelled")
)

// allowed transitions for a ride status
var allowedTransitions = map[string][]string{
	model.RideAccepted:   {model.RideInProgress},
	model.RideInProgress: {model.RideCompleted},
}

// RideService is the service for the ride
type RideService struct {
	rides   *repository.RideRepository
	drivers *repository.DriverRepository
}

// NewRideService creates a new ride service
func NewRideService(rides *repository.RideRepository, drivers *repository.DriverRepository) *RideService {
	return &RideService{rides: rides, drivers: drivers}
}

// Create creates a new ride
func (s *RideService) Create(
	ctx context.Context,
	riderID string,
	req model.CreateRideRequest,
) (*model.Ride, error) {
	hasActiveRide, err := s.rides.HasActiveRideByRider(ctx, riderID)
	// check if the rider has an active ride
	if err != nil {
		return nil, fmt.Errorf("check rider active ride: %w", err)
	}
	if hasActiveRide {
		return nil, ErrRiderHasActiveRide
	}
	// estimate the fare
	_, fare := EstimateFare(
		*req.PickupLat,
		*req.PickupLng,
		*req.DropoffLat,
		*req.DropoffLng,
	)
	// create the ride
	ride := &model.Ride{
		RiderID:    riderID,
		PickupLat:  *req.PickupLat,
		PickupLng:  *req.PickupLng,
		DropoffLat: *req.DropoffLat,
		DropoffLng: *req.DropoffLng,
		Fare:       &fare,
	}
	createdRide, err := s.rides.Create(ctx, ride)
	if err != nil {
		return nil, mapRideError(err)
	}
	return createdRide, nil
}

// GetByID gets a ride by its ID
func (s *RideService) GetByID(ctx context.Context, rideID, userID, role string) (*model.Ride, error) {
	ride, err := s.rides.GetByID(ctx, rideID)
	if err != nil {
		return nil, mapRideError(err)
	}
	// check if the user can view the ride
	if !canViewRide(ride, userID, role) {
		return nil, ErrForbidden
	}

	return ride, nil
}

// List lists all rides for a user
func (s *RideService) List(ctx context.Context, userID, role string) ([]model.Ride, error) {
	switch role {
	case model.RoleRider:
		return s.rides.ListByRider(ctx, userID)
	case model.RoleDriver:
		return s.rides.ListByDriver(ctx, userID)
	default:
		return nil, ErrForbidden
	}
}

// ListAvailable lists all available rides
func (s *RideService) ListAvailable(ctx context.Context) ([]model.Ride, error) {
	return s.rides.ListAvailable(ctx)
}

// UpdateStatus updates the status of a ride
func (s *RideService) UpdateStatus(ctx context.Context, rideID, driverID, newStatus string) (*model.Ride, error) {
	ride, err := s.rides.GetByID(ctx, rideID)
	if err != nil {
		return nil, mapRideError(err)
	}
	// check if the driver is the driver of the ride
	if ride.DriverID == nil || *ride.DriverID != driverID {
		return nil, ErrForbidden
	}
	// check if the new status is valid
	if !isValidTransition(ride.Status, newStatus) {
		return nil, ErrRideStatusInvalid
	}

	// update the status of the ride
	updated, err := s.rides.UpdateStatus(ctx, rideID, driverID, ride.Status, newStatus)
	if err != nil {
		return nil, mapRideError(err)
	}

	return updated, nil
}

// canViewRide checks if a user can view a ride
func canViewRide(ride *model.Ride, userID, role string) bool {
	// check if the user can view the ride
	switch role {
	// check if the user is the rider of the ride
	case model.RoleRider:
		return ride.RiderID == userID
	case model.RoleDriver:
		// check if the driver is the driver of the ride
		if ride.Status == model.RidePending && ride.DriverID == nil {
			return true
		}
		// check if the driver is the driver of the ride
		return ride.DriverID != nil && *ride.DriverID == userID
	default:
		// check if the user is not a rider or driver
		return false
	}
}

// isValidTransition checks if a transition is valid
func isValidTransition(from, to string) bool {
	for _, status := range allowedTransitions[from] {
		if status == to {
			return true
		}
	}
	return false
}

// mapRideError maps a repository error to a service error
func mapRideError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	// check if the ride is not found
	case errors.Is(err, repository.ErrRideNotFound):
		return ErrRideNotFound
	// check if the ride is already taken
	case errors.Is(err, repository.ErrRideAlreadyTaken):
		return ErrRideAlreadyTaken
	// check if the rider has an active ride
	case errors.Is(err, repository.ErrRiderHasActiveRide):
		return ErrRiderHasActiveRide
	// check if the driver has an active ride
	case errors.Is(err, repository.ErrDriverHasActiveRide):
		return ErrDriverHasActiveRide
	// check if the ride status is changed
	case errors.Is(err, repository.ErrRideStatusChanged):
		return ErrRideStatusInvalid
	default:
		return fmt.Errorf("ride repository: %w", err)
	}
}

// Accept accepts a ride by a driver
func (s *RideService) Accept(ctx context.Context, rideID, driverID string) (*model.Ride, error) {
	driverProfile, err := s.drivers.GetByUserID(ctx, driverID)
	if err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			return nil, ErrDriverNotFound
		}
		// return the error
		return nil, fmt.Errorf("get driver profile: %w", err)
	}

	// check if the driver not online
	if !driverProfile.IsOnline {
		return nil, ErrDriverOffline
	}

	// check if the driver has an active ride
	hasActiveRide, err := s.rides.HasActiveRideByDriver(ctx, driverID)
	if err != nil {
		return nil, fmt.Errorf("check driver active ride: %w", err)
	}
	// check if the driver has an active ride
	if hasActiveRide {
		return nil, ErrDriverHasActiveRide
	}

	// accept the ride
	ride, err := s.rides.Accept(ctx, rideID, driverID)
	if err != nil {
		return nil, mapRideError(err)
	}

	return ride, nil
}

// CancelByRider cancels a pending or accepted ride owned by the rider.
func (s *RideService) CancelByRider(ctx context.Context, rideID, riderID string) (*model.Ride, error) {
	ride, err := s.rides.GetByID(ctx, rideID)
	if err != nil {
		return nil, mapRideError(err)
	}

	if ride.RiderID != riderID {
		return nil, ErrForbidden
	}
	if ride.Status != model.RidePending && ride.Status != model.RideAccepted {
		return nil, ErrRideCannotCancel
	}

	cancelledRide, err := s.rides.CancelByRider(ctx, rideID, riderID, ride.Status)
	if err != nil {
		if errors.Is(err, repository.ErrRideStatusChanged) {
			return nil, ErrRideCannotCancel
		}
		return nil, fmt.Errorf("cancel ride by rider: %w", err)
	}

	return cancelledRide, nil
}

// CancelByDriver cancels an accepted ride assigned to the driver.
func (s *RideService) CancelByDriver(ctx context.Context, rideID, driverID string) (*model.Ride, error) {
	ride, err := s.rides.GetByID(ctx, rideID)
	if err != nil {
		return nil, mapRideError(err)
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return nil, ErrForbidden
	}
	if ride.Status != model.RideAccepted {
		return nil, ErrRideCannotCancel
	}

	cancelledRide, err := s.rides.CancelByDriver(ctx, rideID, driverID, ride.Status)
	if err != nil {
		if errors.Is(err, repository.ErrRideStatusChanged) {
			return nil, ErrRideCannotCancel
		}
		return nil, fmt.Errorf("cancel ride by driver: %w", err)
	}

	return cancelledRide, nil
}
