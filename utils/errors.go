package utils

import "net/http"

// AppError is a custom error type carrying an HTTP status code.
// This lets the service layer return meaningful errors without
// importing Echo or knowing about HTTP — keeping clean architecture.
type AppError struct {
	Code    int    // HTTP status code
	Message string // Client-safe message
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a custom error with a status code and message.
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Pre-defined common errors used across the app.
var (
	ErrEmailExists        = NewAppError(http.StatusBadRequest, "Email already registered")
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "Invalid email or password")
	ErrUnauthorized       = NewAppError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden          = NewAppError(http.StatusForbidden, "You do not have permission to perform this action")
	ErrZoneNotFound       = NewAppError(http.StatusNotFound, "Parking zone not found")
	ErrZoneFull           = NewAppError(http.StatusConflict, "Parking zone is full, no spots available")
	ErrReservationNotFound = NewAppError(http.StatusNotFound, "Reservation not found")
	ErrAlreadyCancelled   = NewAppError(http.StatusBadRequest, "Reservation is already cancelled")
	ErrDuplicatePlate     = NewAppError(http.StatusConflict, "This license plate already has an active reservation in this zone")
	ErrInternal           = NewAppError(http.StatusInternalServerError, "Something went wrong, please try again later")
)
