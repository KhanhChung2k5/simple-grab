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
	// bind the request body to the RegisterRequest struct
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// if the request body is invalid, return a bad request response
		response.BadRequest(c, err)
		return
	}

	// call the Register service to register the user
	result, err := h.auth.Register(c.Request.Context(), req)
	if err != nil {
		// if the email already exists, return a conflict response
		if errors.Is(err, repository.ErrEmailExists) {
			response.BadRequest(c, err)
			return
		}
		// if the role is invalid, return a bad request response
		if errors.Is(err, service.ErrInvalidRole) {
			response.BadRequest(c, err)
			return
		}
		// if the error is not nil, return an internal server error response
		response.InternalServerError(c, errors.New("failed to register"))
		return
	}

	response.Success(c, result)
}

// Login handles the login request
func (h *AuthHandler) Login(c *gin.Context) {
	// bind the request body to the LoginRequest struct
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	// call the Login service to login the user
	result, err := h.auth.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Unauthorized(c, err)
			return
		}
		// if the error is not nil, return an internal server error response
		response.InternalServerError(c, errors.New("failed to login"))
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
		response.Unauthorized(c, errors.New("unauthorized"))
		return
	}

	// cast the user from the context to the User struct
	user, ok := userVal.(*model.User)
	if !ok {
		response.InternalServerError(c, errors.New("invalid user in context"))
		return
	}

	// if the user is found, return a success response
	response.Success(c, user)
}
