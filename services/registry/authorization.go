package dripservices

import (
	"context"
	"fmt"
	"net/http"
	"registry-backend/ent"
	"registry-backend/ent/publisherpermission"
	"registry-backend/ent/schema"
	"registry-backend/mapper"

	"github.com/labstack/echo/v4"
)

// RequireActiveUser reloads the verified principal from the database on every
// request. Bans and role changes must not depend on cached token claims.
func RequireActiveUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	id, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	user, err := client.User.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	if err != nil {
		return nil, err
	}
	if user.Status == schema.UserStatusTypeBanned {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Permission denied")
	}
	return user, nil
}

func RequireAdmin(ctx context.Context, client *ent.Client) (*ent.User, error) {
	user, err := RequireActiveUser(ctx, client)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Permission denied")
	}
	return user, nil
}

// Share the Publisher policy with existing Registry operations. Callers choose
// the transaction client and translate permission errors for their API contract.
func assertPublisherPermissions(ctx context.Context, client *ent.Client, publisherID, userID string, permissions []schema.PublisherPermissionType) error {
	publisher, err := client.Publisher.Get(ctx, publisherID)
	if err != nil {
		return fmt.Errorf("fail to query publisher by id: %s %w", publisherID, err)
	}
	allowed, err := publisher.QueryPublisherPermissions().Where(
		publisherpermission.PermissionIn(permissions...),
		publisherpermission.UserIDEQ(userID),
	).Exist(ctx)
	if err != nil {
		return fmt.Errorf("fail to query publisher permission :%w", err)
	}
	if !allowed {
		return newErrorPermission("user '%s' doesn't have required permission on publisher '%s' ", userID, publisherID)
	}
	return nil
}
