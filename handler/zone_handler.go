package handler

import (
	"net/http"
	"strconv"

	"spotsync-api/dto"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
)

// ZoneHandler handles HTTP for the parking zones module.
type ZoneHandler struct {
	zoneService *service.ZoneService
}

func NewZoneHandler(zoneService *service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

// Create handles POST /api/v1/zones (admin only)
func (h *ZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewError("Invalid request body", nil))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError("Validation failed", utils.FormatValidationErrors(err)))
	}

	zone, err := h.zoneService.Create(req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusCreated,
		dto.NewSuccess("Parking zone created successfully", zone))
}

// GetAll handles GET /api/v1/zones (public)
func (h *ZoneHandler) GetAll(c echo.Context) error {
	zones, err := h.zoneService.GetAll()
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK,
		dto.NewSuccess("Parking zones retrieved successfully", zones))
}

// GetByID handles GET /api/v1/zones/:id (public)
func (h *ZoneHandler) GetByID(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewError("Invalid zone id", nil))
	}

	zone, err := h.zoneService.GetByID(id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK,
		dto.NewSuccess("Parking zone retrieved successfully", zone))
}

// Update handles PUT /api/v1/zones/:id (admin only)
func (h *ZoneHandler) Update(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewError("Invalid zone id", nil))
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewError("Invalid request body", nil))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError("Validation failed", utils.FormatValidationErrors(err)))
	}

	zone, err := h.zoneService.Update(id, req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK,
		dto.NewSuccess("Parking zone updated successfully", zone))
}

// Delete handles DELETE /api/v1/zones/:id (admin only)
func (h *ZoneHandler) Delete(c echo.Context) error {
	id, err := parseIDParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewError("Invalid zone id", nil))
	}

	if err := h.zoneService.Delete(id); err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(http.StatusOK,
		dto.NewSuccess("Parking zone deleted successfully", nil))
}

// parseIDParam is a small helper to read a uint64 id from the URL.
func parseIDParam(c echo.Context, name string) (uint64, error) {
	idStr := c.Param(name)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}
