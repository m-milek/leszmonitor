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

// @Summary	Get the latest monitor result
// @Tags		monitors
// @Produce	json
// @Param		monitorId	path		string	true	"Monitor ID"
// @Success	200			{object}	MonitorResult
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId}/results/latest [get]
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

// @Summary	List monitor results
// @Tags		monitors
// @Produce	json
// @Param		monitorId	path		string	true	"Monitor ID"
// @Param		page		query		int		false	"Page number"								default(1)
// @Param		perPage		query		int		false	"Items per page"							default(20)
// @Param		from		query		string	false	"Only results after this time (RFC3339)"	format(date-time)
// @Success	200			{array}		MonitorResult
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId}/results [get]
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
