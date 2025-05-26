package monitors

import (
	"fmt"
)

type BaseMonitor struct {
	Name        string      `json:"name" bson:"name"`
	Description string      `json:"description" bson:"description"`
	Interval    int         `json:"interval" bson:"interval"` // in seconds
	Timeout     int         `json:"timeout" bson:"timeout"`   // in seconds
	OwnerId     string      `json:"owner_id" bson:"owner_id"`
	Type        MonitorType `json:"type" bson:"type"`
}

// BaseMonitorResponse is the base of any monitor response.
// A monitor run can either succeed, fail, or error.
// 1. Success: The monitor ran successfully and returned the expected results.
// 2. Failure: The monitor ran but did not return the expected results (e.g., HTTP status code, response time).
// 3. Error: The monitor encountered an error while running (e.g., network issues, timeout).
type BaseMonitorResponse struct {
	Status    MonitorResponseStatus `json:"status" bson:"status"`
	Duration  int64                 `json:"duration" bson:"duration"`                     // in milliseconds
	Timestamp int64                 `json:"timestamp" bson:"timestamp"`                   // Unix timestamp in seconds
	Errors    []string              `json:"errors,omitempty" bson:"errors,omitempty"`     // List of errors encountered during the monitor run
	Failures  []string              `json:"failures,omitempty" bson:"failures,omitempty"` // List of specific failures, if any
}

func (b *BaseMonitorResponse) GetStatus() MonitorResponseStatus {
	return b.Status
}
func (b *BaseMonitorResponse) GetDuration() int64 {
	return b.Duration
}
func (b *BaseMonitorResponse) GetTimestamp() int64 {
	return b.Timestamp
}
func (b *BaseMonitorResponse) GetErrors() []string {
	return b.Errors
}
func (b *BaseMonitorResponse) GetFailures() []string {
	return b.Failures
}

type IMonitorResponse interface {
	GetStatus() MonitorResponseStatus
	GetDuration() int64
	GetTimestamp() int64
	GetErrors() []string
	GetFailures() []string
}

type IMonitor interface {
	Run(client HttpClient) (IMonitorResponse, error)
	GetName() string
	GetDescription() string
	GetInterval() int
	GetTimeout() int
}

type MonitorType string

const (
	Http MonitorType = "http"
)

type MonitorResponseStatus string

const (
	Success MonitorResponseStatus = "success"
	Failure                       = "failure"
	Error                         = "error"
)

func (m *BaseMonitor) validate() error {
	if m.Name == "" {
		return fmt.Errorf("monitor name cannot be empty")
	}
	if m.Interval <= 0 {
		return fmt.Errorf("monitor interval must be greater than zero")
	}
	if m.Timeout <= 0 {
		return fmt.Errorf("monitor timeout must be greater than zero")
	}
	if m.Type == "" {
		return fmt.Errorf("monitor type cannot be empty")
	}
	return nil
}

func (b *BaseMonitorResponse) addErrorMsg(err string) {
	if b.Errors == nil {
		b.Errors = make([]string, 0)
	}
	b.Errors = append(b.Errors, err)
	b.Status = Error
}

func (b *BaseMonitorResponse) addFailureMsg(err string) {
	if b.Failures == nil {
		b.Failures = make([]string, 0)
	}
	b.Failures = append(b.Failures, err)

	// If the status is not already set to Error, set it to Failure. Error always takes precedence.
	if b.Status != Error {
		b.Status = Failure
	}
}
