package handler

import (
	"errors"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"github.com/gin-gonic/gin"
)

// RideHandler is the handler for the ride service
type RideHandler struct {
	rides *service.RideService
}

// NewRideHandler creates a new RideHandler
func NewRideHandler(rides *service.RideService) *RideHandler {
	return &RideHandler{rides: rides}
}

// Create creates a new ride
func (h *RideHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.CreateRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	ride, err := h.rides.Create(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRiderHasActiveRide):
			response.Conflict(c, err)

		default:
			response.InternalServerError(c, errors.New("failed to create ride"))
			return
		}
		return
	}
	response.Success(c, ride)

}

// GetByID gets a ride by ID
func (h *RideHandler) GetByID(c *gin.Context) {
	rideID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	ride, err := h.rides.GetByID(c.Request.Context(), rideID, userID, role)
	if err != nil {
		if errors.Is(err, service.ErrRideNotFound) {
			response.NotFound(c, err)
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, err)
			return
		}
		response.InternalServerError(c, errors.New("failed to get ride by ID"))
		return
	}
	response.Success(c, ride)
}

// List returns rides for the authenticated user (filtered by role)
func (h *RideHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	rides, err := h.rides.List(c.Request.Context(), userID, role)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, err)
			return
		}
		response.InternalServerError(c, errors.New("failed to get list of rides"))
		return
	}
	response.Success(c, rides)
}

// ListAvailable lists pending rides without a driver
func (h *RideHandler) ListAvailable(c *gin.Context) {
	availableRides, err := h.rides.ListAvailable(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, errors.New("failed to get list of available rides"))
		return
	}
	response.Success(c, availableRides)
}

// Accept lets a driver accept a pending ride
func (h *RideHandler) Accept(c *gin.Context) {
	rideID := c.Param("id")
	driverID := c.GetString("user_id")

	ride, err := h.rides.Accept(c.Request.Context(), rideID, driverID)
	if err != nil {
		if errors.Is(err, service.ErrRideAlreadyTaken) {
			response.Conflict(c, err)
			return
		}
		if errors.Is(err, service.ErrRideNotFound) {
			response.NotFound(c, err)
			return
		}

		if errors.Is(err, service.ErrDriverOffline) {
			response.Forbidden(c, err)
			return
		}
		if errors.Is(err, service.ErrDriverHasActiveRide) {
			response.Conflict(c, err)
			return
		}
		response.InternalServerError(c, errors.New("failed to accept ride"))
		return
	}
	response.Success(c, ride)
}

// UpdateStatus updates ride status (accepted → in_progress → completed)
func (h *RideHandler) UpdateStatus(c *gin.Context) {
	rideID := c.Param("id")
	driverID := c.GetString("user_id")

	var req model.UpdateRideStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	ride, err := h.rides.UpdateStatus(c.Request.Context(), rideID, driverID, req.Status)
	if err != nil {
		if errors.Is(err, service.ErrRideNotFound) {
			response.NotFound(c, err)
			return
		}
		if errors.Is(err, service.ErrRideStatusInvalid) {
			response.BadRequest(c, err)
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, err)
			return
		}
		response.InternalServerError(c, errors.New("failed to update status"))
		return
	}
	response.Success(c, ride)
}
