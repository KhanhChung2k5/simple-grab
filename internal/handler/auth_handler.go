package handler

import (
	"errors"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/repository"
	"github.com/KhanhChung2k5/simple-grab/internal/response"
	"github.com/KhanhChung2k5/simple-grab/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler is the handler for the auth routes
type AuthHandler struct {
	auth *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register handles the register request
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// if the request is invalid, return a bad request response
		response.BadRequest(c, err)
		return
	}
	// call the Register service to register the user
	result, err := h.auth.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEmailExists):
			response.Conflict(c, err)
		case errors.Is(err, service.ErrInvalidRole):
			response.BadRequest(c, err)
		default:
			internalError(c, "failed to register")
		}
		return
	}
	// if the registration is successful, return a success response
	response.Success(c, result)
}

// Login handles the login request
func (h *AuthHandler) Login(c *gin.Context) {
	// bind the request body to the LoginRequest struct
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// if the request is invalid, return a bad request response
		response.BadRequest(c, err)
		return
	}

	// call the Login service to login the user
	result, err := h.auth.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			// if the credentials are invalid, return an unauthorized response
			response.Unauthorized(c, err)
			return
		}
		// if the error is not nil, return an internal server error response
		internalError(c, "failed to login")
		return
	}

	// if the login is successful, return a success response
	response.Success(c, result)
}

// Me handles the me request
func (h *AuthHandler) Me(c *gin.Context) {
	// get the user from the context
	userVal, exists := c.Get("user")
	if !exists {
		// if the user is not found, return an unauthorized response
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	// cast the user from the context to the User struct
	user, ok := userVal.(*model.User)
	if !ok {
		internalError(c, "invalid user in context")
		return
	}

	// if the user is found, return a success response
	response.Success(c, user)
}
