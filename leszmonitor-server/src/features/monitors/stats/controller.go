package stats

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/httpx"
)

type MonitorStatsAPIController struct {
	service MonitorStatsService
}

func NewMonitorStatsAPIController(service MonitorStatsService) MonitorStatsAPIController {
	return MonitorStatsAPIController{
		service: service,
	}
}

// @Summary	Get monitor stats
// @Tags		monitors
// @Produce	json
// @Param		monitorId	path		string	true	"Monitor ID"
// @Param		from		query		string	true	"Range start (RFC3339)"	format(date-time)
// @Param		to			query		string	true	"Range end (RFC3339)"	format(date-time)
// @Success	200			{object}	MonitorStats
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId}/stats [get]
func (c *MonitorStatsAPIController) GetStatsByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorID, err := uuid.Parse(r.PathValue("monitorId"))
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID format"))
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

	_, ok := auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	stats, svcErr := c.service.GetStatsByMonitorID(ctx, monitorID, from, to)
	if svcErr != nil {
		httpx.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, stats)
}
