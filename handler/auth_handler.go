package handler

import (
	"net/http"

	"spotsync-api/dto"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles HTTP for the auth module.
type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError("Invalid request body", nil))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError("Validation failed", utils.FormatValidationErrors(err)))
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated,
		dto.NewSuccess("User registered successfully", user))
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError("Invalid request body", nil))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewError("Validation failed", utils.FormatValidationErrors(err)))
	}

	result, err := h.authService.Login(req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK,
		dto.NewSuccess("Login successful", result))
}
