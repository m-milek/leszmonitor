package middleware

import (
	"fmt"
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/services"
)

type SlugSource string

const (
	SlugSourcePath  SlugSource = "path"
	SlugSourceQuery SlugSource = "query"
)

// RequireInstanceAdmin checks if the user is an instance admin.
func RequireInstanceAdmin() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userClaims, ok := auth.ExtractUserOrRespond(ctx, w, r)
			if !ok {
				return
			}

			if !userClaims.IsInstanceAdmin {
				httpx.RespondError(ctx, w, http.StatusForbidden, fmt.Errorf("requires instance admin privileges"))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireInstanceAdminHandler checks if the user is an instance admin, accepting an [http.Handler].
func RequireInstanceAdminHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userClaims, ok := auth.ExtractUserOrRespond(ctx, w, r)
		if !ok {
			return
		}

		if !userClaims.IsInstanceAdmin {
			httpx.RespondError(ctx, w, http.StatusForbidden, fmt.Errorf("requires instance admin privileges"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireSelf checks if the authenticated user matches the username parameter in the URL.
func RequireSelf(usernameParam string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userClaims, ok := auth.ExtractUserOrRespond(ctx, w, r)
			if !ok {
				return
			}

			if userClaims.IsInstanceAdmin {
				next.ServeHTTP(w, r)
				return
			}

			targetUsername := r.PathValue(usernameParam)
			if targetUsername != userClaims.Username {
				httpx.RespondError(
					ctx,
					w,
					http.StatusForbidden,
					fmt.Errorf("access denied to another user's resources"),
				)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequirePermission checks if the user has the required permission.
func RequirePermission(
	authzService services.IAuthzMiddlewareService,
	perm auth.Permission,
) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userClaims, ok := auth.ExtractUserOrRespond(ctx, w, r)
			if !ok {
				return
			}

			if userClaims.IsInstanceAdmin {
				next.ServeHTTP(w, r)
				return
			}

			hasPerm, err := authzService.CheckUserPermission(ctx, userClaims.Username, perm)
			if err != nil {
				httpx.RespondError(ctx, w, http.StatusInternalServerError, err)
				return
			}
			if !hasPerm {
				httpx.RespondError(
					ctx,
					w,
					http.StatusForbidden,
					fmt.Errorf("user does not have required permission: %s", perm.Name),
				)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
