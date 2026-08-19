package handler

import (
	"context"
	"errors"
	"regexp"

	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/gin-gonic/gin"
)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

func rideIDParam(c *gin.Context) (string, bool) {
	rideID := c.Param("id")
	if !uuidPattern.MatchString(rideID) {
		response.BadRequest(c, errors.New("invalid ride id"))
		return "", false
	}
	return rideID, true
}

func internalError(c *gin.Context, message string) {
	if errors.Is(c.Request.Context().Err(), context.DeadlineExceeded) {
		response.GatewayTimeout(c, errors.New("request timeout"))
		return
	}
	response.InternalServerError(c, errors.New(message))
}
