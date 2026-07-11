package handler

import (
	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health is a handler for the health check endpoint.
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, "OK")
}