package results

import (
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
)

type IMonitorResult interface {
	GetID() uuid.UUID
	GetMonitorID() uuid.UUID
	GetStatus() kind.MonitorStatus
	GetIsManuallyTriggered() bool
	GetDurationMs() int64
	GetDetails() IMonitorResultDetails
	GetCreatedAt() time.Time
	AddFailure(fail string)
	SetDuration(duration int64)
	SetDetails(details IMonitorResultDetails)
	GetErrorDetails() ErrorDetails
}

type ErrorDetails struct {
	ErrorMessage string   `json:"errorMessage,omitempty"`
	Errors       []string `json:"errors,omitempty"`
	Failures     []string `json:"failures,omitempty"`
}

type baseMonitorResult struct {
	ID                  uuid.UUID          `json:"id"                     db:"id"`
	MonitorID           uuid.UUID          `json:"monitorId"              db:"monitor_id"`
	Status              kind.MonitorStatus `json:"status" db:"status"`
	IsManuallyTriggered bool               `json:"isManuallyTriggered"    db:"is_manually_triggered"`
	DurationMs          int64              `json:"durationMs"             db:"duration_ms"`
	ErrorDetailsJSON    []byte             `json:"-"                      db:"error_details"`
	ErrorDetails        *ErrorDetails      `json:"errorDetails,omitempty" db:"-"`
	CreatedAt           time.Time          `json:"createdAt"              db:"created_at"`
}

type MonitorResult struct {
	baseMonitorResult

	MonitorType string                `json:"monitorType" db:"kind"`
	DetailsJSON []byte                `json:"-"           db:"details"`
	Details     IMonitorResultDetails `json:"details"     db:"-"`
}

func NewMonitorResult(
	monitorID uuid.UUID,
	monitorType kind.ProbeType,
	status kind.MonitorStatus,
	isManuallyTriggered bool,
	durationMs int64,
	errorMessage string,
	details IMonitorResultDetails,
) MonitorResult {
	res := MonitorResult{
		baseMonitorResult: baseMonitorResult{
			ID:                  uuid.New(),
			MonitorID:           monitorID,
			Status:              status,
			IsManuallyTriggered: isManuallyTriggered,
			DurationMs:          durationMs,
			CreatedAt:           time.Now().UTC(),
		},
		MonitorType: string(monitorType),
		Details:     details,
	}
	if errorMessage != "" {
		res.ErrorDetails = &ErrorDetails{ErrorMessage: errorMessage}
	}
	return res
}

func (m *MonitorResult) GetID() uuid.UUID {
	return m.ID
}

func (m *MonitorResult) GetMonitorID() uuid.UUID {
	return m.MonitorID
}

func (m *MonitorResult) GetStatus() kind.MonitorStatus {
	return m.Status
}

func (m *MonitorResult) GetIsManuallyTriggered() bool {
	return m.IsManuallyTriggered
}

func (m *MonitorResult) GetDurationMs() int64 {
	return m.DurationMs
}

func (m *MonitorResult) GetErrorDetails() ErrorDetails {
	if m.ErrorDetails == nil {
		return ErrorDetails{}
	}
	return *m.ErrorDetails
}

func (m *MonitorResult) GetDetails() IMonitorResultDetails {
	return m.Details
}

func (m *MonitorResult) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *MonitorResult) AddFailure(fail string) {
	if m.ErrorDetails == nil {
		m.ErrorDetails = &ErrorDetails{}
	}
	m.ErrorDetails.Failures = append(m.ErrorDetails.Failures, fail)
	m.Status = kind.MonitorStatusDown
}

func (m *MonitorResult) SetDuration(duration int64) {
	m.DurationMs = duration
}

func (m *MonitorResult) SetDetails(details IMonitorResultDetails) {
	m.Details = details
}
