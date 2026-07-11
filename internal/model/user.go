package model

import "time"

const (
	RoleRider  = "rider"
	RoleDriver = "driver"
	RoleAdmin  = "admin"
)

type User struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Phone        *string   `json:"phone,omitempty" gorm:"size:20"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;size:255;not null"`
	Role         string    `json:"role" gorm:"size:20;not null"`
	CreatedAt    time.Time `json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

type Driver struct {
	ID             string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         string     `json:"user_id" gorm:"type:uuid;uniqueIndex;not null"`
	VehicleType    string     `json:"vehicle_type" gorm:"size:50;not null;default:car"`
	IsOnline       bool       `json:"is_online" gorm:"not null;default:false"`
	Lat            *float64   `json:"lat,omitempty"`
	Lng            *float64   `json:"lng,omitempty"`
	LastLocationAt *time.Time `json:"last_location_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (Driver) TableName() string {
	return "drivers"
}
