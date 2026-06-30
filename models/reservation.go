package models

import "time"

// Reservation represents a driver booking a spot in a parking zone.
type Reservation struct {
	ID           uint64      `gorm:"primaryKey" json:"id"`
	UserID       uint64      `gorm:"not null;index" json:"user_id"`
	ZoneID       uint64      `gorm:"not null;index" json:"zone_id"`
	LicensePlate string    `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       string    `gorm:"type:varchar(20);default:active;not null" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Associations (used with Preload)
	User User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone ParkingZone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
}
