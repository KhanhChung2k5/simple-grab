package model

type UpdateDriverOnlineRequest struct {
	IsOnline *bool `json:"is_online"`
}

type UpdateDriverOnlineResponse struct {
	IsOnline *bool `json:"is_online"`
}

type UpdateDriverOnlineError struct {
	IsOnline *bool `json:"is_online"`
}

type UpdateDriverLocationRequest struct {
	Latitude  *float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude *float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
}

type UpdateDriverLocationResponse struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type UpdateDriverLocationError struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}
