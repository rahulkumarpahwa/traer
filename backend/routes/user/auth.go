package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rahulkumarpahwa/traer/auth"
	"github.com/rahulkumarpahwa/traer/utils"
)

func (u *UserHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("authentication required"))
				return
			}
			http.Error(w, "Error reading cookie", http.StatusInternalServerError)
			return
		}

		// Parse and validate token
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return u.Config.Jwt.Secret, nil
		})

		if err != nil || !token.Valid {
			utils.ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid or expired token"))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "email", claims.Email)

		// Call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
