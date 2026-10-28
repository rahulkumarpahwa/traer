package auth

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func GetUserId(r *http.Request) (uuid.UUID, error) {

	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("userid not found")
	}

	return uuid.MustParse(userId), nil
}
