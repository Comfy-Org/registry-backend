package mapper

import (
	"context"
	"errors"
	auth "registry-backend/server/middleware/authentication"
)

func GetUserIDFromContext(ctx context.Context) (string, error) {
	user, ok := ctx.Value(auth.UserContextKey).(*auth.UserDetails)
	if !ok || user == nil {
		return "", errors.New("no user in context")
	}

	if user.ID == "" {
		return "", errors.New("no user id in context")
	}

	return user.ID, nil
}
