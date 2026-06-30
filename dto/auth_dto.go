package dto

import "time"

// RegisterRequest is the body for POST /api/v1/auth/register
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"omitempty,oneof=driver admin"`
}

// LoginRequest is the body for POST /api/v1/auth/login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse is the safe (password-free) representation of a user.
type UserResponse struct {
	ID        uint64      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginUser is the minimal user info returned on login.
type LoginUser struct {
	ID    uint64   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// LoginResponse wraps the JWT token and basic user info.
type LoginResponse struct {
	Token string    `json:"token"`
	User  LoginUser `json:"user"`
}
