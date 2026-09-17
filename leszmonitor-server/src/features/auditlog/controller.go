package auditlog

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/platform/util"
)

type AuditLogAPIController struct {
	service AuditLogService
}

func NewAuditLogAPIController(service AuditLogService) AuditLogAPIController {
	return AuditLogAPIController{
		service: service,
	}
}

func (c *AuditLogAPIController) GetAuditLogByQueryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pagination, paginationErr := util.PaginationFromRequest(r)
	if paginationErr != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, paginationErr)
		return
	}

	filters, filtersErr := audit.AuditLogFilterFromRequest(r)
	if filtersErr != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, filtersErr)
		return
	}

	_, ok := auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	results, err := c.service.GetEntries(ctx, *filters, *pagination)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, results)
}
