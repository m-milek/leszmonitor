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

type IMonitor interface {
	Run() error
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

func (m *BaseMonitorResponse) addErrorMsg(err string) {
	if m.Errors == nil {
		m.Errors = make([]string, 0)
	}
	m.Errors = append(m.Errors, err)
	m.Status = Error
}

func (m *BaseMonitorResponse) addFailureMsg(err string) {
	if m.Failures == nil {
		m.Failures = make([]string, 0)
	}
	m.Failures = append(m.Failures, err)

	// If the status is not already set to Error, set it to Failure. Error always takes precedence.
	if m.Status != Error {
		m.Status = Failure
	}
}
