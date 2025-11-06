package models

import "time"

// User represents a system user
type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never return in JSON
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Role represents a user role
type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Permission represents a system permission
type Permission struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// UserAdminArea represents a user's access to administrative areas
type UserAdminArea struct {
	ID          int  `json:"id"`
	UserID      int  `json:"userId"`
	RegionID    *int `json:"regionId,omitempty"`
	DistrictID  *int `json:"districtId,omitempty"`
	SubcountyID *int `json:"subcountyId,omitempty"`
	FacilityID  *int `json:"facilityId,omitempty"`
}
