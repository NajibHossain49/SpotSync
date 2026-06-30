package repository

import (
	"errors"

	"spotsync-api/models"

	"gorm.io/gorm"
)

// ZoneRepository handles all DB operations for parking zones.
type ZoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) *ZoneRepository {
	return &ZoneRepository{db: db}
}

// ZoneWithAvailability is a query result holding a zone plus its
// dynamically computed active reservation count.
type ZoneWithAvailability struct {
	models.ParkingZone
	ActiveCount int `gorm:"column:active_count"`
}

func (r *ZoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *ZoneRepository) Update(zone *models.ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *ZoneRepository) Delete(id uint64) error {
	return r.db.Delete(&models.ParkingZone{}, id).Error
}

// FindByID returns a single zone or (nil, nil) if not found.
func (r *ZoneRepository) FindByID(id uint64) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	err := r.db.First(&zone, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

// FindAllWithAvailability returns all zones with active_count computed via a subquery.
// available_spots = total_capacity - active_count (computed in the service).
func (r *ZoneRepository) FindAllWithAvailability() ([]ZoneWithAvailability, error) {
	var results []ZoneWithAvailability

	// Subquery counts active reservations per zone.
	subQuery := r.db.Model(&models.Reservation{}).
		Select("COUNT(*)").
		Where("reservations.zone_id = parking_zones.id AND reservations.status = ?", "active")

	err := r.db.Model(&models.ParkingZone{}).
		Select("parking_zones.*, (?) AS active_count", subQuery).
		Order("parking_zones.id ASC").
		Find(&results).Error

	return results, err
}

// FindByIDWithAvailability returns one zone with its active_count.
func (r *ZoneRepository) FindByIDWithAvailability(id uint64) (*ZoneWithAvailability, error) {
	var result ZoneWithAvailability

	subQuery := r.db.Model(&models.Reservation{}).
		Select("COUNT(*)").
		Where("reservations.zone_id = parking_zones.id AND reservations.status = ?", "active")

	err := r.db.Model(&models.ParkingZone{}).
		Select("parking_zones.*, (?) AS active_count", subQuery).
		Where("parking_zones.id = ?", id).
		First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}
