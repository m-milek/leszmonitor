package statuschange

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/platform/util"
)

type MonitorStatusChangeAPIController struct {
	service MonitorStatusChangeService
}

func NewMonitorStatusChangeAPIController(service MonitorStatusChangeService) MonitorStatusChangeAPIController {
	return MonitorStatusChangeAPIController{
		service: service,
	}
}

// @Summary	Get monitor status history
// @Tags		monitors
// @Produce	json
// @Param		monitorId	path		string	true	"Monitor ID"
// @Param		page		query		int		false	"Page number"			default(1)
// @Param		perPage		query		int		false	"Items per page"		default(20)
// @Param		from		query		string	true	"Range start (RFC3339)"	format(date-time)
// @Param		to			query		string	true	"Range end (RFC3339)"	format(date-time)
// @Success	200			{array}		MonitorStatusPeriod
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId}/status-history [get]
func (c *MonitorStatusChangeAPIController) GetStatusHistoryByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorID, err := uuid.Parse(r.PathValue("monitorId"))
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID format"))
		return
	}

	pagination, paginationErr := util.PaginationFromRequest(r)
	if paginationErr != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, paginationErr)
		return
	}

	fromParam := r.URL.Query().Get("from")
	toParam := r.URL.Query().Get("to")

	from, err := time.Parse(time.RFC3339, fromParam)
	if err != nil {
		httpx.RespondError(
			ctx,
			w,
			http.StatusBadRequest,
			errors.New("invalid 'from' parameter format, expected RFC3339"),
		)
		return
	}

	to, err := time.Parse(time.RFC3339, toParam)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid 'to' parameter format, expected RFC3339"))
		return
	}

	statusHistory, svcErr := c.service.GetStatusHistoryByMonitorID(ctx, monitorID, pagination, from, to)
	if svcErr != nil {
		httpx.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, statusHistory)
}
