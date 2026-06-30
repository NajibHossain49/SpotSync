package handler

import (
	"net/http"

	"spotsync-api/dto"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
)

// handleServiceError converts a service-layer error into the correct
// HTTP status + standard error body. This is the centralized error
// handling that keeps raw GORM errors from leaking to the client.
func handleServiceError(c echo.Context, err error) error {
	if appErr, ok := err.(*utils.AppError); ok {
		return c.JSON(appErr.Code, dto.NewError(appErr.Message, nil))
	}
	// Fallback: never expose internal details.
	return c.JSON(http.StatusInternalServerError,
		dto.NewError("Internal server error", nil))
}
