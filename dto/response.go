package dto

// SuccessResponse is the standard wrapper for all successful responses.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is the standard wrapper for all error responses.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// NewSuccess builds a success payload.
func NewSuccess(message string, data interface{}) SuccessResponse {
	return SuccessResponse{Success: true, Message: message, Data: data}
}

// NewError builds an error payload.
func NewError(message string, errs interface{}) ErrorResponse {
	return ErrorResponse{Success: false, Message: message, Errors: errs}
}
