package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rahulkumarpahwa/traer/config"
	"github.com/rahulkumarpahwa/traer/types"
)

// Claims defines the JWT payload
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *types.User, config *config.Config) (*string, error) {
	if user == nil || config == nil {
		return nil, fmt.Errorf("User or Config missing to generate token!")
	}

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1 day expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Jwt.Secret)
	if err != nil {
		return nil, nil
	}

	return &tokenString, nil
}

func ParseAndValidateToken(cookie *http.Cookie, config *config.Config) (*Claims, error) {
	// Parse and validate token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return config.Jwt.Secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
