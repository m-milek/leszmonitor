package handlers

import (
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/services"
	util2 "github.com/m-milek/leszmonitor/util"
)

func GetAuditLogByQueryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination, paginationErr := util2.PaginationFromRequest(r)
	if paginationErr != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, paginationErr)
		return
	}

	filters, filtersErr := security.AuditLogFilterFromRequest(r)
	if filtersErr != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, filtersErr)
		return
	}

	userClaims, ok := authorization.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	results, err := services.AuditLogService.GetEntries(ctx, userClaims, *filters, *pagination)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, results)
}
