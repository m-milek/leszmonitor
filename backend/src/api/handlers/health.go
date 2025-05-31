package handlers

import (
	"github.com/m-milek/leszmonitor/api/api_util"
	"net/http"
	"time"
)

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func GetHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Status:    "OK",
		Timestamp: time.Now(),
	}

	util.RespondJSON(w, http.StatusOK, response)
}
