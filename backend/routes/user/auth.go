package user

import (
	"context"
	"fmt"
	"net/http"
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

		// Parse and Valid Token:
		auth := auth.Auth{Config: u.Config}
		claims, err := auth.ParseAndValidateToken(cookie)
		if err != nil {
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