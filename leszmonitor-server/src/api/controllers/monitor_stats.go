package controllers

import (
	"errors"
	"net/http"
	"time"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/services"
)

type MonitorStatsAPIController struct {
	service services.MonitorStatsService
}

func NewMonitorStatsAPIController(service services.MonitorStatsService) MonitorStatsAPIController {
	return MonitorStatsAPIController{
		service: service,
	}
}

type AverageLatencyResponse struct {
	AverageLatency float64 `json:"averageLatency"`
}

func (c *MonitorStatsAPIController) GetAverageLatencyByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondError(ctx, w, http.StatusBadRequest, errors.New("monitor ID is required"))
		return
	}

	fromParam := r.URL.Query().Get("from")
	toParam := r.URL.Query().Get("to")

	from, err := time.Parse(time.RFC3339, fromParam)
	if err != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid 'from' parameter format, expected RFC3339"))
		return
	}

	to, err := time.Parse(time.RFC3339, toParam)
	if err != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid 'to' parameter format, expected RFC3339"))
		return
	}

	_, ok := authorization.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	avgLatency, svcErr := c.service.GetAverageLatencyByMonitorID(ctx, monitorID, from, to)
	if svcErr != nil {
		util.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}

	response := AverageLatencyResponse{
		AverageLatency: avgLatency,
	}

	util.RespondJSON(ctx, w, http.StatusOK, response)
}
