package results

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/platform/util"
)

type MonitorResultsAPIController struct {
	service IMonitorResultsService
}

func NewMonitorResultsAPIController(service IMonitorResultsService) MonitorResultsAPIController {
	return MonitorResultsAPIController{
		service: service,
	}
}

func (c *MonitorResultsAPIController) GetLatestMonitorResultByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorID, err := uuid.Parse(r.PathValue("monitorId"))
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID format"))
		return
	}

	result, svcErr := c.service.GetLatestMonitorResultByMonitorID(ctx, monitorID)
	if svcErr != nil {
		httpx.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, result)
}

func (c *MonitorResultsAPIController) GetMonitorResultsByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
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

	var from *time.Time
	if fromParam := r.URL.Query().Get("from"); fromParam != "" {
		parsed, err := time.Parse(time.RFC3339, fromParam)
		if err != nil {
			httpx.RespondError(
				ctx,
				w,
				http.StatusBadRequest,
				errors.New("invalid 'from' parameter format, expected RFC3339"),
			)
			return
		}
		from = &parsed
	}

	results, svcErr := c.service.GetMonitorResultsByMonitorID(ctx, monitorID, pagination, from)
	if svcErr != nil {
		httpx.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, results)
}
