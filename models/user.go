package models

import "time"

// User represents an account in the system (driver or admin).
type User struct {
	ID        uint64      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // json:"-" => never exposed in responses
	Role      string    `gorm:"type:varchar(20);default:driver;not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// One user can have many reservations
	Reservations []Reservation `gorm:"foreignKey:UserID" json:"-"`
}
