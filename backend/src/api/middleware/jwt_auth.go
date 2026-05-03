package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/log"
)

// JwtAuth middleware validates JWT tokens from the Authorization header.
func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(rw, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userClaims, err := auth.ValidateJwt(tokenString)
		if err != nil {
			log.Api.Warn().Err(err).Msg("Unauthorized access attempt with invalid token")
			http.Error(rw, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Store the user claims in the request context
		ctx := SetUserContext(r.Context(), userClaims)

		// Call the next handler with the updated context
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// Define a context key type to avoid collisions.
type contextKey string

const userClaimsKey contextKey = "userClaims"

// SetUserContext stores user claims in the request context.
func SetUserContext(ctx context.Context, claims *auth.UserClaims) context.Context {
	log.Api.Debug().Msg("Setting user claims in context: " + claims.Username)
	return context.WithValue(ctx, userClaimsKey, claims)
}

// GetUserFromContext retrieves user claims from the request context.
func GetUserFromContext(ctx context.Context) (*auth.UserClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*auth.UserClaims)
	log.Api.Debug().Msgf("Retrieving user claims from context: %v, ok: %v", claims, ok)
	return claims, ok
}

// ProjectAuth bundles the project display ID and the authenticated username for authorization.
type ProjectAuth struct {
	ProjectID string
	Username  string
}

type AuthSourceKind string

var (
	AuthSourceProject     = AuthSourceKind("projectId")
	AuthSourceKindMonitor = AuthSourceKind("monitorId")
)

// ProjectAuthFromRequest extracts the project ID from the URL path and the username from the JWT context.
func ProjectAuthFromRequest(r *http.Request, authSource AuthSourceKind) (*ProjectAuth, error) {
	var projectSlug string
	if authSource == AuthSourceKindMonitor {
		monitorID := r.PathValue("monitorId")
		if monitorID == "" {
			return nil, fmt.Errorf("monitor ID is required")
		}
		monitorUUID, err := uuid.Parse(monitorID)
		if err != nil {
			log.Api.Error().Err(err).Msg("Invalid monitor ID format")
			return nil, fmt.Errorf("invalid monitor ID format")
		}
		monitor, err := db.Get().Monitors().GetMonitorByID(r.Context(), monitorUUID)
		if err != nil {
			log.Api.Error().Err(err).Msg("Failed to get monitor by ID")
			return nil, fmt.Errorf("failed to get monitor")
		}
		project, err := db.Get().Projects().GetProjectBySlug(r.Context(), monitor.GetProjectSlug())
		if err != nil {
			log.Api.Error().Err(err).Msg("Failed to get project by slug")
			return nil, fmt.Errorf("failed to get project")
		}
		projectSlug = project.Slug
	} else if authSource == AuthSourceProject {
		projectSlug = r.PathValue("projectId")
		if projectSlug == "" {
			return nil, fmt.Errorf("project slug is required")
		}
	}

	userClaims, ok := GetUserFromContext(r.Context())
	if !ok {
		return nil, fmt.Errorf("user claims not found in context")
	}
	if userClaims.Username == "" {
		return nil, fmt.Errorf("username is missing in user claims")
	}

	return &ProjectAuth{
		ProjectID: projectSlug,
		Username:  userClaims.Username,
	}, nil
}
