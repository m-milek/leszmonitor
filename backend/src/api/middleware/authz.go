package middleware

import (
	"fmt"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/services"
)

// RequireProjectPermission checks if the user has the required permission for the project.
func RequireProjectPermission(authzService services.IAuthzMiddlewareService, perm models.Permission) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userClaims, ok := authorization.ExtractUserOrRespond(ctx, w, r)
			if !ok {
				return
			}

			if userClaims.IsInstanceAdmin {
				next.ServeHTTP(w, r)
				return
			}

			projectSlug := r.PathValue("projectSlug")
			if projectSlug == "" {
				util.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("project slug not found in URL"))
				return
			}

			hasPerm, err := authzService.CheckProjectPermissionBySlug(ctx, userClaims.Username, projectSlug, perm)
			if err != nil {
				util.RespondError(ctx, w, http.StatusInternalServerError, err)
				return
			}
			if !hasPerm {
				util.RespondError(ctx, w, http.StatusForbidden, fmt.Errorf("user does not have required project permission: %s", perm.Name))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireInstanceAdmin checks if the user is an instance admin.
func RequireInstanceAdmin() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userClaims, ok := authorization.ExtractUserOrRespond(ctx, w, r)
			if !ok {
				return
			}

			if !userClaims.IsInstanceAdmin {
				util.RespondError(ctx, w, http.StatusForbidden, fmt.Errorf("requires instance admin privileges"))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireSelf checks if the authenticated user matches the username parameter in the URL.
func RequireSelf(usernameParam string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userClaims, ok := authorization.ExtractUserOrRespond(ctx, w, r)
			if !ok {
				return
			}

			if userClaims.IsInstanceAdmin {
				next.ServeHTTP(w, r)
				return
			}

			targetUsername := r.PathValue(usernameParam)
			if targetUsername != userClaims.Username {
				util.RespondError(ctx, w, http.StatusForbidden, fmt.Errorf("access denied to another user's resources"))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
