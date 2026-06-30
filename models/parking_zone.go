package models

import "time"

// ParkingZone represents a parking area (general / ev_charging / covered).
type ParkingZone struct {
	ID            uint64      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Type          string    `gorm:"type:varchar(20);not null" json:"type"`
	TotalCapacity int       `gorm:"not null" json:"total_capacity"`
	PricePerHour  float64   `gorm:"not null" json:"price_per_hour"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// One zone can have many reservations
	Reservations []Reservation `gorm:"foreignKey:ZoneID" json:"-"`
}
