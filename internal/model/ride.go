package model

import "time"

const (
	RidePending    = "pending"
	RideAccepted   = "accepted"
	RideInProgress = "in_progress"
	RideCompleted  = "completed"
	RideCancelled  = "cancelled"
)

type Ride struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RiderID    string    `json:"rider_id" gorm:"type:uuid;not null"`
	DriverID   *string   `json:"driver_id,omitempty" gorm:"type:uuid"`
	PickupLat  float64   `json:"pickup_lat" gorm:"not null"`
	PickupLng  float64   `json:"pickup_lng" gorm:"not null"`
	DropoffLat float64   `json:"dropoff_lat" gorm:"not null"`
	DropoffLng float64   `json:"dropoff_lng" gorm:"not null"`
	Status     string    `json:"status" gorm:"size:20;not null;default:pending"`
	Fare       *float64  `json:"fare,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Ride) TableName() string {
	return "rides"
}

type CreateRideRequest struct {
	PickupLat  *float64 `json:"pickup_lat" binding:"required,gte=-90,lte=90"`
	PickupLng  *float64 `json:"pickup_lng" binding:"required,gte=-180,lte=180"`
	DropoffLat *float64 `json:"dropoff_lat" binding:"required,gte=-90,lte=90"`
	DropoffLng *float64 `json:"dropoff_lng" binding:"required,gte=-180,lte=180"`
}

type UpdateRideStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=in_progress completed"`
}
