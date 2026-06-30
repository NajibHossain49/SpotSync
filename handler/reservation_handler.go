package handler

import (
	"net/http"

	"spotsync-api/dto"
	"spotsync-api/middleware"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
)

// ReservationHandler handles HTTP for the reservations module.
type ReservationHandler struct {
	reservationService *service.ReservationService
}

func NewReservationHandler(reservationService *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

// Reserve handles POST /api/v1/reservations (authenticated)
func (h *ReservationHandler) Reserve(c echo.Context) error {
	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewError("Invalid request body", nil))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError("Validation failed", utils.FormatValidationErrors(err)))
	}

	// The JWT middleware injected the user id into the context.
	userID, _ := c.Get(middleware.ContextUserID).(uint64)

	reservation, err := h.reservationService.Reserve(userID, req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusCreated,
		dto.NewSuccess("Reservation confirmed successfully", reservation))
}

// GetMyReservations handles GET /api/v1/reservations/my-reservations (authenticated)
func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, _ := c.Get(middleware.ContextUserID).(uint64)

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK,
		dto.NewSuccess("My reservations retrieved successfully", reservations))
}

// GetAll handles GET /api/v1/reservations (admin only)
func (h *ReservationHandler) GetAll(c echo.Context) error {
	reservations, err := h.reservationService.GetAll()
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK,
		dto.NewSuccess("All reservations retrieved successfully", reservations))
}

// Cancel handles DELETE /api/v1/reservations/:id (authenticated, owner or admin)
func (h *ReservationHandler) Cancel(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewError("Invalid reservation id", nil))
	}

	userID, _ := c.Get(middleware.ContextUserID).(uint64)
	role, _ := c.Get(middleware.ContextRole).(string)

	if err := h.reservationService.Cancel(id, userID, role); err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK,
		dto.NewSuccess("Reservation cancelled successfully", nil))
}
