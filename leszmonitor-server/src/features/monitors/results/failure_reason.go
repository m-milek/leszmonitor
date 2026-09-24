package results

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type FailureReason string

const (
	FailureReasonHTTPRequestFailed          FailureReason = "HTTP_REQUEST_FAILED"
	FailureReasonHTTPResponseBodyReadFailed FailureReason = "HTTP_RESPONSE_BODY_READ_FAILED"
	FailureReasonHTTPStatusCodeMismatch     FailureReason = "HTTP_STATUS_CODE_MISMATCH"
	FailureReasonHTTPResponseBodyMismatch   FailureReason = "HTTP_RESPONSE_BODY_MISMATCH"
	FailureReasonHTTPResponseHeaderMismatch FailureReason = "HTTP_RESPONSE_HEADER_MISMATCH"
	FailureReasonHTTPResponseTimeExceeded   FailureReason = "HTTP_RESPONSE_TIME_EXCEEDED"

	FailureReasonDNSLookupFailed          FailureReason = "DNS_LOOKUP_FAILED"
	FailureReasonDNSInvalidSRVHostname    FailureReason = "DNS_INVALID_SRV_HOSTNAME"
	FailureReasonDNSExpectedRecordMissing FailureReason = "DNS_EXPECTED_RECORD_MISSING"

	FailureReasonTCPConnectionFailed FailureReason = "TCP_CONNECTION_FAILED"
)

type FailureCause string

const (
	FailureCauseTimeout           FailureCause = "TIMEOUT"
	FailureCauseConnectionRefused FailureCause = "CONNECTION_REFUSED"
	FailureCauseUnreachable       FailureCause = "UNREACHABLE"
	FailureCauseDNS               FailureCause = "DNS"
	FailureCauseTLS               FailureCause = "TLS"
	FailureCauseNotFound          FailureCause = "NOT_FOUND"
	FailureCauseTemporary         FailureCause = "TEMPORARY"
	FailureCauseOther             FailureCause = "OTHER"
)

type CauseFailureDetails struct {
	Cause FailureCause `json:"cause"`
}

type StatusCodeMismatchDetails struct {
	Got      int   `json:"got"`
	Expected []int `json:"expected"`
}

type ResponseTimeExceededDetails struct {
	GotMs         int64 `json:"gotMs"`
	ExpectedMaxMs int   `json:"expectedMaxMs"`
}

type HeaderMismatch struct {
	Name     string `json:"name"`
	Got      string `json:"got"`
	Expected string `json:"expected"`
}

type HeaderMismatchDetails struct {
	Headers []HeaderMismatch `json:"headers"`
}

type BodyMismatchDetails struct {
	Pattern string `json:"pattern"`
}

type MissingRecordsDetails struct {
	Missing []string `json:"missing"`
}

type Failure struct {
	Reason  FailureReason `json:"reason"`
	Details any           `json:"details,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type Failures []Failure

func (f Failures) Value() (driver.Value, error) {
	if len(f) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (f *Failures) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*f = nil
		return nil
	case string:
		return json.Unmarshal([]byte(v), f)
	case []byte:
		return json.Unmarshal(v, f)
	default:
		return fmt.Errorf("cannot scan %T into Failures", src)
	}
}
