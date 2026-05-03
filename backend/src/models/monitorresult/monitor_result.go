package monitorresult

import (
	"github.com/google/uuid"
	consts "github.com/m-milek/leszmonitor/models/consts"
)

type IMonitorResult interface {
	GetMonitorID() uuid.UUID
	GetIsSuccess() bool
	GetIsManuallyTriggered() bool
	GetDurationMs() int64
	GetErrorMessage() string
	GetDetails() IMonitorResultDetails
	GetCreatedAt() string
	AddError(err string)
	AddFailure(fail string)
	SetDuration(duration int64)
	SetDetails(details IMonitorResultDetails)
	IsError() bool
	GetErrors() []string
}

type errorDetails struct {
	ErrorMessage string   `json:"errorMessage,omitempty"`
	Errors       []string `json:"errors,omitempty"`
	Failures     []string `json:"failures,omitempty"`
}

type baseMonitorResult struct {
	MonitorID           uuid.UUID    `json:"monitorId"           db:"monitor_id"`
	IsSuccess           bool         `json:"isSuccess"           db:"is_success"`
	IsManuallyTriggered bool         `json:"isManuallyTriggered" db:"is_manually_triggered"`
	DurationMs          int64        `json:"durationMs"          db:"duration_ms"`
	ErrorMessage        string       `json:"errorMessage,omitempty" db:"error_message"`
	ErrorDetails        errorDetails `json:"errorDetails,omitempty" db:"-"`
	CreatedAt           string       `json:"createdAt"           db:"created_at"`
}

type MonitorResult struct {
	baseMonitorResult
	MonitorType string                `json:"monitorType" db:"kind"`
	DetailsJSON []byte                `json:"-"           db:"details"`
	Details     IMonitorResultDetails `json:"details"     db:"-"`
}

func NewMonitorResult(monitorID uuid.UUID, monitorType consts.MonitorConfigType, isSuccess bool, isManuallyTriggered bool, durationMs int64, errorMessage string, details IMonitorResultDetails, createdAt string) *MonitorResult {
	return &MonitorResult{
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

func (m *MonitorResult) GetMonitorID() uuid.UUID {
	return m.MonitorID
}

func (m *MonitorResult) GetIsSuccess() bool {
	return m.IsSuccess
}

func (m *MonitorResult) GetIsManuallyTriggered() bool {
	return m.IsManuallyTriggered
}

func (m *MonitorResult) GetDurationMs() int64 {
	return m.DurationMs
}

func (m *MonitorResult) GetErrorMessage() string {
	return m.ErrorMessage
}

func (m *MonitorResult) GetDetails() IMonitorResultDetails {
	return m.Details
}

func (m *MonitorResult) GetCreatedAt() string {
	return m.CreatedAt
}

func (m *MonitorResult) AddError(err string) {
	m.ErrorDetails.Errors = append(m.ErrorDetails.Errors, err)
	m.IsSuccess = false
}

func (m *MonitorResult) AddFailure(fail string) {
	m.ErrorDetails.Failures = append(m.ErrorDetails.Failures, fail)
	m.IsSuccess = false
}

func (m *MonitorResult) SetDuration(duration int64) {
	m.DurationMs = duration
}

func (m *MonitorResult) SetDetails(details IMonitorResultDetails) {
	m.Details = details
}

func (m *MonitorResult) IsError() bool {
	return len(m.ErrorDetails.Errors) > 0 || m.ErrorMessage != ""
}

func (m *MonitorResult) GetErrors() []string {
	return m.ErrorDetails.Errors
}
