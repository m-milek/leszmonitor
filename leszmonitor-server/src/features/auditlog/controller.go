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

// @Summary	Query the audit log
// @Tags		audit-log
// @Produce	json
// @Param		page		query		int		false	"Page number"		default(1)
// @Param		per_page	query		int		false	"Items per page"	default(20)
// @Param		username	query		string	false	"Filter by username"
// @Param		resourceId	query		string	false	"Filter by resource ID"	format(uuid)
// @Param		action		query		string	false	"Filter by action"
// @Param		isSuccess	query		bool	false	"Filter by outcome"
// @Param		traceId		query		string	false	"Filter by trace ID"
// @Param		startDate	query		string	false	"Entries created at or after (RFC3339)"		format(date-time)
// @Param		endDate		query		string	false	"Entries created at or before (RFC3339)"	format(date-time)
// @Success	200			{array}		audit.AuditLogEntry
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/audit-log [get]
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
