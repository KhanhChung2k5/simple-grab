package handler

import (
	"errors"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"github.com/gin-gonic/gin"
)

// DriverHandler is the handler for the driver service
type DriverHandler struct {
	drivers *service.DriverService
}

// NewDriverHandler creates a new DriverHandler
func NewDriverHandler(drivers *service.DriverService) *DriverHandler {
	return &DriverHandler{drivers: drivers}
}

// UpdateOnline updates the driver's online/offline status
func (h *DriverHandler) UpdateOnline(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.UpdateDriverOnlineRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.IsOnline == nil {
		response.BadRequest(c, errors.New("is_online is required"))
		return
	}

	if err := h.drivers.SetOnline(userID, *req.IsOnline); err != nil {
		response.InternalServerError(c, errors.New("failed to update online status"))
		return
	}

	response.Success(c, model.UpdateDriverOnlineResponse{IsOnline: req.IsOnline})
}

// UpdateLocation updates the driver's current location
func (h *DriverHandler) UpdateLocation(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.UpdateDriverLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.drivers.UpdateLocation(userID, *req.Latitude, *req.Longitude); err != nil {
		response.InternalServerError(c, errors.New("failed to update location"))
		return
	}

	response.Success(c, model.UpdateDriverLocationResponse{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
}
