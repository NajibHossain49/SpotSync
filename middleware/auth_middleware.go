package middleware

import (
	"net/http"
	"strings"

	"spotsync-api/dto"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
)

// Context keys used to store JWT claims into the Echo context.
const (
	ContextUserID = "userID"
	ContextRole   = "role"
)

// JWTAuth verifies the Bearer token and injects user id + role into the context.
func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized,
					dto.NewError("Missing authorization token", nil))
			}

			// Expect the format: "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized,
					dto.NewError("Invalid authorization header format", nil))
			}

			claims, err := utils.ParseToken(parts[1], secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized,
					dto.NewError("Invalid or expired token", nil))
			}

			// Inject the verified claims so handlers can read them.
			c.Set(ContextUserID, claims.UserID)
			c.Set(ContextRole, claims.Role)
			return next(c)
		}
	}
}

// AdminOnly rejects requests from non-admin users (403 Forbidden).
// It must run AFTER JWTAuth so the role is already in the context.
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, _ := c.Get(ContextRole).(string)
		if role != "admin" {
			return c.JSON(http.StatusForbidden,
				dto.NewError("Admin access required", nil))
		}
		return next(c)
	}
}
