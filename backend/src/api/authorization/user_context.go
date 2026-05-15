package authorization

import (
	"context"
	"fmt"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/log"
)

// Define a context key type to avoid collisions.
type contextKey string

const userClaimsKey contextKey = "userClaims"

func SetUserInContext(ctx context.Context, claims *auth.UserClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

// GetUserFromContext retrieves user claims from the request context.
func GetUserFromContext(ctx context.Context) (*auth.UserClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*auth.UserClaims)
	return claims, ok
}

// ExtractUserOrRespond returns the user from context or writes a 401 response and returns nil, false.
func ExtractUserOrRespond(ctx context.Context, w http.ResponseWriter, r *http.Request) (*auth.UserClaims, bool) {
	logger := log.FromContext(ctx)
	user, ok := GetUserFromContext(ctx)
	if !ok {
		logger.Warn().Msg("User claims not found in context")
		util.RespondError(ctx, w, http.StatusUnauthorized, fmt.Errorf("user claims not found in context"))
		return nil, false
	}
	return user, true
}

func GetUsernameFromRequest(ctx context.Context) (*string, error) {
	userClaims, ok := GetUserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user claims not found in context")
	}
	if userClaims.Username == "" {
		return nil, fmt.Errorf("username is missing in user claims")
	}

	return &userClaims.Username, nil
}
