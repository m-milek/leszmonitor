package monitors_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/stretchr/testify/assert"
)

func TestReceivePushHandler_InvalidInput(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name      string
		monitorID string
		query     string
		body      string
	}{
		{name: "Invalid monitor ID", monitorID: "not-a-uuid"},
		{name: "Invalid status", monitorID: validID, query: "?status=maybe"},
		{name: "Non-numeric latency", monitorID: validID, query: "?latency=fast"},
		{name: "Negative latency", monitorID: validID, query: "?latency=-1"},
		{name: "Body too large", monitorID: validID, body: strings.Repeat("a", 5001)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := monitors.NewMonitorAPIController(nil)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/push/"+tt.monitorID+tt.query, strings.NewReader(tt.body))
			req.SetPathValue("monitorId", tt.monitorID)
			rec := httptest.NewRecorder()

			controller.ReceivePushHandler(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}
