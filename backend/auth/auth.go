package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rahulkumarpahwa/traer/config"
	"github.com/rahulkumarpahwa/traer/types"
)

func GenerateJWT(user *types.User, config *config.Config) (*string, error) {
	if user == nil || config == nil {
		return nil, fmt.Errorf("User or Config missing to generate token!")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.ID,
		"username": user.Username,
		"email":    user.Email,
		"time":     time.Now(),
	})
	tokenString, err := token.SignedString(config.Jwt.Secret)
	if err != nil {
		return nil, nil
	}

	return &tokenString, nil
}
