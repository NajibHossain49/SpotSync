package repository

import (
	"errors"

	"spotsync-api/models"
	"spotsync-api/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ReservationRepository handles all DB operations for reservations.
type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

// CreateWithCapacityCheck is the CORE concurrency-safe operation.
//
// It solves the "EV Spot Bottleneck" race condition:
// If two drivers grab the last spot at the same millisecond, a naive
// read-then-write would let both succeed, over-filling the zone.
//
// We prevent this with a database TRANSACTION + ROW-LEVEL LOCK (FOR UPDATE):
//  1. Lock the zone row. Any other transaction trying to lock the same row
//     must WAIT until we commit/rollback.
//  2. Count active reservations safely (no one else can change them now).
//  3. If there is room, insert the reservation.
//  4. Commit -> lock released -> the next waiting request now sees the
//     updated count and is correctly rejected if the zone is full.
func (r *ReservationRepository) CreateWithCapacityCheck(reservation *models.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone

		// 1. Acquire a row-level write lock on the zone (SELECT ... FOR UPDATE).
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, reservation.ZoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrZoneNotFound
			}
			return err
		}

		// 2. Count current active reservations for this zone (inside the lock).
		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zone.ID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// Optional business rule: prevent the same plate from holding two
		// active reservations in the same zone.
		var dupCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND license_plate = ? AND status = ?",
				zone.ID, reservation.LicensePlate, "active").
			Count(&dupCount).Error; err != nil {
			return err
		}
		if dupCount > 0 {
			return utils.ErrDuplicatePlate
		}

		// 3. Capacity check.
		if int(activeCount) >= zone.TotalCapacity {
			return utils.ErrZoneFull
		}

		// 4. Safe to create the reservation. Commit releases the lock.
		reservation.Status = "active"
		return tx.Create(reservation).Error
	})
}

// FindByID returns a reservation by id, or (nil, nil) if not found.
func (r *ReservationRepository) FindByID(id uint64) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

// FindByUserID returns all reservations for a user, with the zone preloaded.
func (r *ReservationRepository) FindByUserID(userID uint64) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

// FindAll returns every reservation with user and zone preloaded (admin view).
func (r *ReservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").Preload("Zone").
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

// UpdateStatus changes a reservation's status (e.g. to "cancelled").
func (r *ReservationRepository) UpdateStatus(id uint64, status string) error {
	return r.db.Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", status).Error
}
