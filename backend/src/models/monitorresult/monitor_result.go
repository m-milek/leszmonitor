package monitorresult

import (
	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/monitors"
)

type IMonitorResult interface {
	GetMonitorID() uuid.UUID
	GetIsSuccess() string
	GetIsManuallyTriggered() bool
	GetDurationMs() int64
	GetErrorMessage() string
	GetDetails() IMonitorResultDetails
	GetCreatedAt() string
}

type errorDetails struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type baseMonitorResult struct {
	MonitorID           uuid.UUID    `json:"monitorId"`
	IsSuccess           string       `json:"isSuccess"`
	IsManuallyTriggered bool         `json:"isManuallyTriggered"`
	DurationMs          int64        `json:"durationMs"`
	ErrorMessage        string       `json:"errorMessage,omitempty"`
	ErrorDetails        errorDetails `json:"errorDetails,omitempty"`
	CreatedAt           string       `json:"createdAt"`
}

type monitorResult struct {
	baseMonitorResult
	MonitorType string                `json:"monitorType"`
	Details     IMonitorResultDetails `json:"details"`
}

func newMonitorResult(monitorID uuid.UUID, monitorType monitors.MonitorConfigType, isSuccess string, isManuallyTriggered bool, durationMs int64, errorMessage string, details IMonitorResultDetails, createdAt string) *monitorResult {
	return &monitorResult{
		baseMonitorResult: baseMonitorResult{
			MonitorID:           monitorID,
			IsSuccess:           isSuccess,
			IsManuallyTriggered: isManuallyTriggered,
			DurationMs:          durationMs,
			ErrorMessage:        errorMessage,
			CreatedAt:           createdAt,
		},
		MonitorType: string(monitorType),
		Details:     details,
	}
}

func (m *monitorResult) GetMonitorID() uuid.UUID {
	return m.MonitorID
}

func (m *monitorResult) GetIsSuccess() string {
	return m.IsSuccess
}

func (m *monitorResult) GetIsManuallyTriggered() bool {
	return m.IsManuallyTriggered
}

func (m *monitorResult) GetDurationMs() int64 {
	return m.DurationMs
}

func (m *monitorResult) GetErrorMessage() string {
	return m.ErrorMessage
}

func (m *monitorResult) GetDetails() IMonitorResultDetails {
	return m.Details
}

func (m *monitorResult) GetCreatedAt() string {
	return m.CreatedAt
}
