package controllers

import (
	"net/http"
	"time"

	util "github.com/m-milek/leszmonitor/api/api_util"
)

type healthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func GetHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := healthCheckResponse{
		Status:    "OK",
		Timestamp: time.Now(),
	}

	util.RespondJSON(ctx, w, http.StatusOK, response)
}
