package app

import (
	"net/http"
	"time"

	"github.com/m-milek/leszmonitor/platform/httpx"
)

type healthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// @Summary	Health check
// @Tags		health
// @Produce	json
// @Success	200	{object}	healthCheckResponse
// @Failure	401	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/health [get]
func GetHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := healthCheckResponse{
		Status:    "OK",
		Timestamp: time.Now(),
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, response)
}
