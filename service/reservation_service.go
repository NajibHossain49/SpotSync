package service

import (
	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"
	"spotsync-api/utils"
)

// ReservationService contains the business logic for reservations.
type ReservationService struct {
	reservationRepo *repository.ReservationRepository
}

func NewReservationService(reservationRepo *repository.ReservationRepository) *ReservationService {
	return &ReservationService{reservationRepo: reservationRepo}
}

// Reserve creates a reservation atomically (capacity check + row lock live in the repo).
func (s *ReservationService) Reserve(userID uint64, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	reservation := &models.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	// The repository runs the locking transaction and may return
	// ErrZoneNotFound / ErrZoneFull / ErrDuplicatePlate (all *AppError).
	if err := s.reservationRepo.CreateWithCapacityCheck(reservation); err != nil {
		if _, ok := err.(*utils.AppError); ok {
			return nil, err
		}
		return nil, utils.ErrInternal
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}, nil
}

// GetMyReservations lists the requester's own reservations.
func (s *ReservationService) GetMyReservations(userID uint64) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, utils.ErrInternal
	}

	responses := make([]dto.MyReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneBrief{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}
	return responses, nil
}

// GetAll lists every reservation in the system (admin only).
func (s *ReservationService) GetAll() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, utils.ErrInternal
	}

	responses := make([]dto.AdminReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User: dto.UserBrief{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			Zone: dto.ZoneBrief{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt,
		})
	}
	return responses, nil
}

// Cancel sets a reservation's status to "cancelled".
// Drivers may only cancel their OWN reservations; admins may cancel any.
func (s *ReservationService) Cancel(reservationID, requesterID uint64, requesterRole string) error {
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		return utils.ErrInternal
	}
	if reservation == nil {
		return utils.ErrReservationNotFound
	}

	// Ownership check: a non-admin can only cancel their own reservation.
	if requesterRole != "admin" && reservation.UserID != requesterID {
		return utils.ErrForbidden
	}

	if reservation.Status == "cancelled" {
		return utils.ErrAlreadyCancelled
	}

	if err := s.reservationRepo.UpdateStatus(reservationID, "cancelled"); err != nil {
		return utils.ErrInternal
	}
	return nil
}
