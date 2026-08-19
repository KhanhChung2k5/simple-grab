package model

import "time"

// RegisterRequest is the request body for the register endpoint
type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Role     string  `json:"role" binding:"required,oneof=rider driver"`
	Phone    *string `json:"phone,omitempty"`
}

// RegisterResponse is the response body for the register endpoint
type RegisterResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Phone     *string   `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest is the request body for the login endpoint
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response body for the login endpoint
type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	User        User      `json:"user"`
	CreatedAt   time.Time `json:"created_at"`
}
