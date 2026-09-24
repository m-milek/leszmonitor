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

func (c *MonitorStatsAPIController) GetStatsByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorID, err := uuid.Parse(r.PathValue("monitorId"))
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID"))
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
