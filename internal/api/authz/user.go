package authz

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

// UserIDInCTX checks if the userID
// equals the authenticated user in the context.
func UserIDInCTX(ctx context.Context, userID string) error {
	if GetCtxData(ctx).UserID != userID {
		return zerrors.ThrowPermissionDenied(nil, "AUTH-Bohd2", "Errors.User.UserIDWrong")
	}
	return nil
}
