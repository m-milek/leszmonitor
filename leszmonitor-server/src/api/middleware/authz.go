package middleware

import (
	"fmt"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/services"
)

type SlugSource string

const (
	SlugSourcePath  SlugSource = "path"
	SlugSourceQuery SlugSource = "query"
)

// RequireProjectPermission checks if the user has the required permission for the project.
// The slugSource determines whether the project slug is extracted from the URL path or query parameters.
func RequireProjectPermission(
	authzService services.IAuthzMiddlewareService,
	perm models.Permission,
	slugSource SlugSource,
) func(http.HandlerFunc) http.HandlerFunc {
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

			var projectSlug string
			switch slugSource {
			case SlugSourcePath:
				projectSlug = r.PathValue("projectSlug")
			case SlugSourceQuery:
				projectSlug = r.URL.Query().Get("projectSlug")
			}

			if projectSlug == "" {
				util.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("project slug not found in %s", slugSource))
				return
			}

			hasPerm, err := authzService.CheckProjectPermissionBySlug(ctx, userClaims.Username, projectSlug, perm)
			if err != nil {
				util.RespondError(ctx, w, http.StatusInternalServerError, err)
				return
			}
			if !hasPerm {
				util.RespondError(
					ctx,
					w,
					http.StatusForbidden,
					fmt.Errorf("user does not have required project permission: %s", perm.Name),
				)
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

// RequireInstanceAdminHandler checks if the user is an instance admin, accepting an [http.Handler].
func RequireInstanceAdminHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
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

// RequireMonitorPermission checks if the user has the required permission for the monitor.
func RequireMonitorPermission(
	authzService services.IAuthzMiddlewareService,
	perm models.Permission,
) func(http.HandlerFunc) http.HandlerFunc {
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

			monitorID := r.PathValue("monitorId")
			if monitorID == "" {
				util.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("monitor ID not found in URL"))
				return
			}

			hasPerm, err := authzService.CheckMonitorPermissionByID(ctx, userClaims.Username, monitorID, perm)
			if err != nil {
				util.RespondError(ctx, w, http.StatusInternalServerError, err)
				return
			}
			if !hasPerm {
				util.RespondError(
					ctx,
					w,
					http.StatusForbidden,
					fmt.Errorf("user does not have required monitor permission: %s", perm.Name),
				)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireProjectPermissionByIDQuery checks if the user has the required permission for the project.
// The project ID is expected to be a query parameter named "projectId".
// If "projectId" is missing, it will pass through (for instance admins to query globally), but if present, it validates permission.
func RequireProjectPermissionByIDQuery(
	authzService services.IAuthzMiddlewareService,
	perm models.Permission,
) func(http.HandlerFunc) http.HandlerFunc {
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

			projectID := r.URL.Query().Get("projectId")
			if projectID == "" {
				// Non-admins must provide projectId (validated by filter logic later),
				// but here we can just reject if missing since they aren't admins.
				util.RespondError(
					ctx,
					w,
					http.StatusBadRequest,
					fmt.Errorf("projectId query parameter is required for non-admins"),
				)
				return
			}

			hasPermission, err := authzService.CheckProjectPermissionByID(ctx, userClaims.Username, projectID, perm)
			if err != nil {
				util.RespondError(ctx, w, http.StatusInternalServerError, err)
				return
			}
			if !hasPermission {
				util.RespondError(
					ctx,
					w,
					http.StatusForbidden,
					fmt.Errorf("user does not have required project permission: %s", perm.Name),
				)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
