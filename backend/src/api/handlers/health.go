package handlers

import (
	"github.com/m-milek/leszmonitor/api/api_util"
	"net/http"
	"time"
)

type healthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func GetHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := healthCheckResponse{
		Status:    "OK",
		Timestamp: time.Now(),
	}

	util.RespondJSON(w, http.StatusOK, response)
}
