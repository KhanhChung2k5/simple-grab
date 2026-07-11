package model

import "time"

type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Role     string  `json:"role" binding:"required,oneof=rider driver"`
	Phone    *string `json:"phone,omitempty"`
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Phone     *string   `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string     `json:"access_token"`
	ExpiresIn   int        `json:"expires_in"`
	User        User       `json:"user"`
}
