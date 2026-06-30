package dto

import "time"

// CreateReservationRequest is the body for POST /api/v1/reservations
type CreateReservationRequest struct {
	ZoneID       uint64   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

// ReservationResponse is the full reservation representation.
type ReservationResponse struct {
	ID           uint64      `json:"id"`
	UserID       uint64      `json:"user_id"`
	ZoneID       uint64      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ZoneBrief is a small zone summary embedded inside reservation responses.
type ZoneBrief struct {
	ID   uint64   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// UserBrief is a small user summary embedded inside admin reservation responses.
type UserBrief struct {
	ID    uint64   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// MyReservationResponse is used for GET /reservations/my-reservations (driver view).
type MyReservationResponse struct {
	ID           uint64      `json:"id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	Zone         ZoneBrief `json:"zone"`
	CreatedAt    time.Time `json:"created_at"`
}

// AdminReservationResponse is used for GET /reservations (admin view).
type AdminReservationResponse struct {
	ID           uint64      `json:"id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	User         UserBrief `json:"user"`
	Zone         ZoneBrief `json:"zone"`
	CreatedAt    time.Time `json:"created_at"`
}
