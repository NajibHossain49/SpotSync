package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims is the custom payload embedded inside the token.
// We store user id and role so middleware can identify the requester
// and enforce permissions without hitting the database every time.
type JWTClaims struct {
	UserID uint64   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken signs a JWT for the given user.
func GenerateToken(userID uint64, role, secret string, expiryHours int) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken validates the signature and returns the claims.
func ParseToken(tokenString, secret string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect (HMAC).
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
