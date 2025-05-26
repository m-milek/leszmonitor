package monitors

type IMonitorResponse interface {
	GetStatus() MonitorResponseStatus
	GetDuration() int64
	GetTimestamp() int64
	GetErrors() []string
	GetFailures() []string
}

type MonitorResponseStatus string

const (
	Success MonitorResponseStatus = "success"
	Failure                       = "failure"
	Error                         = "error"
)

// baseMonitorResponse is the base of any monitor response.
// A monitor run can either succeed, fail, or error.
// 1. Success: The monitor ran successfully and returned the expected results.
// 2. Failure: The monitor ran but did not return the expected results (e.g., HTTP status code, response time).
// 3. Error: The monitor encountered an error while running (e.g., network issues, timeout).
type baseMonitorResponse struct {
	Status    MonitorResponseStatus `json:"status" bson:"status"`
	Duration  int64                 `json:"duration" bson:"duration"`                     // in milliseconds
	Timestamp int64                 `json:"timestamp" bson:"timestamp"`                   // Unix timestamp in seconds
	Errors    []string              `json:"errors,omitempty" bson:"errors,omitempty"`     // List of errors encountered during the monitor run
	Failures  []string              `json:"failures,omitempty" bson:"failures,omitempty"` // List of specific failures, if any
}

func (b *baseMonitorResponse) addErrorMsg(err string) {
	if b.Errors == nil {
		b.Errors = make([]string, 0)
	}
	b.Errors = append(b.Errors, err)
	b.Status = Error
}

func (b *baseMonitorResponse) addFailureMsg(err string) {
	if b.Failures == nil {
		b.Failures = make([]string, 0)
	}
	b.Failures = append(b.Failures, err)

	// If the status is not already set to Error, set it to Failure. Error always takes precedence.
	if b.Status != Error {
		b.Status = Failure
	}
}

func (b *baseMonitorResponse) GetStatus() MonitorResponseStatus {
	return b.Status
}

func (b *baseMonitorResponse) GetDuration() int64 {
	return b.Duration
}

func (b *baseMonitorResponse) GetTimestamp() int64 {
	return b.Timestamp
}

func (b *baseMonitorResponse) GetErrors() []string {
	return b.Errors
}

func (b *baseMonitorResponse) GetFailures() []string {
	return b.Failures
}
