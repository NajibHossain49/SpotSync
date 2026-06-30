package service

import (
	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"
	"spotsync-api/utils"
)

// ZoneService contains the business logic for parking zones.
type ZoneService struct {
	zoneRepo *repository.ZoneRepository
}

func NewZoneService(zoneRepo *repository.ZoneRepository) *ZoneService {
	return &ZoneService{zoneRepo: zoneRepo}
}

// Create stores a new parking zone (admin only — role checked in handler/middleware).
func (s *ZoneService) Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, utils.ErrInternal
	}

	// A brand new zone has no reservations, so all spots are available.
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}

// GetAll returns every zone with dynamically calculated availability.
func (s *ZoneService) GetAll() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAllWithAvailability()
	if err != nil {
		return nil, utils.ErrInternal
	}

	responses := make([]dto.ZoneResponse, 0, len(zones))
	for _, z := range zones {
		available := z.TotalCapacity - z.ActiveCount
		if available < 0 {
			available = 0
		}
		responses = append(responses, dto.ZoneResponse{
			ID:             z.ID,
			Name:           z.Name,
			Type:           z.Type,
			TotalCapacity:  z.TotalCapacity,
			AvailableSpots: available,
			PricePerHour:   z.PricePerHour,
			CreatedAt:      z.CreatedAt,
			UpdatedAt:      z.UpdatedAt,
		})
	}
	return responses, nil
}

// GetByID returns a single zone with availability.
func (s *ZoneService) GetByID(id uint64) (*dto.ZoneResponse, error) {
	z, err := s.zoneRepo.FindByIDWithAvailability(id)
	if err != nil {
		return nil, utils.ErrInternal
	}
	if z == nil {
		return nil, utils.ErrZoneNotFound
	}

	available := z.TotalCapacity - z.ActiveCount
	if available < 0 {
		available = 0
	}
	return &dto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt,
		UpdatedAt:      z.UpdatedAt,
	}, nil
}

// Update modifies an existing zone (admin only). Only provided fields change.
func (s *ZoneService) Update(id uint64, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return nil, utils.ErrInternal
	}
	if zone == nil {
		return nil, utils.ErrZoneNotFound
	}

	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Type != nil {
		zone.Type = *req.Type
	}
	if req.TotalCapacity != nil {
		zone.TotalCapacity = *req.TotalCapacity
	}
	if req.PricePerHour != nil {
		zone.PricePerHour = *req.PricePerHour
	}

	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, utils.ErrInternal
	}
	// Re-fetch availability after update.
	return s.GetByID(id)
}

// Delete removes a zone (admin only).
func (s *ZoneService) Delete(id uint64) error {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		return utils.ErrInternal
	}
	if zone == nil {
		return utils.ErrZoneNotFound
	}
	if err := s.zoneRepo.Delete(id); err != nil {
		return utils.ErrInternal
	}
	return nil
}
