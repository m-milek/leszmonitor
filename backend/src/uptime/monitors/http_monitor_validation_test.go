package monitors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateStatusCode(t *testing.T) {
	monitor := setupTestHttpMonitor()
	response := createMockResponse(200, "", nil)

	// Test matching status code
	monitorResponse := newHttpMonitorResponse()
	monitor.validateStatusCode(response, monitorResponse)
	assert.Equal(t, Success, monitorResponse.Base.Status)
	assert.Empty(t, monitorResponse.FailedAspects)

	// Test mismatched status code
	wrongStatus := 404
	monitor.ExpectedStatusCode = &wrongStatus
	monitorResponse = newHttpMonitorResponse()
	monitor.validateStatusCode(response, monitorResponse)
	assert.EqualValues(t, Failure, monitorResponse.Base.Status)
	assert.Contains(t, monitorResponse.FailedAspects, statusCodeAspect)
	assert.Contains(t, monitorResponse.Base.Failures[0], "Unexpected status code")
}

func TestValidateResponseTime(t *testing.T) {
	monitor := setupTestHttpMonitor()

	// Test response time within limit
	monitorResponse := newHttpMonitorResponse()
	elapsed := 500 * time.Millisecond
	monitor.validateResponseTime(elapsed, monitorResponse)
	assert.EqualValues(t, Success, monitorResponse.Base.Status)
	assert.Empty(t, monitorResponse.FailedAspects)

	// Test response time exceeding limit
	elapsed = 1500 * time.Millisecond
	monitorResponse = newHttpMonitorResponse()
	monitor.validateResponseTime(elapsed, monitorResponse)
	assert.EqualValues(t, Failure, monitorResponse.Base.Status)
	assert.Contains(t, monitorResponse.FailedAspects, responseTimeAspect)
	assert.Contains(t, monitorResponse.Base.Failures[0], "Response time exceeded")
}

func TestValidateResponseHeaders(t *testing.T) {
	monitor := setupTestHttpMonitor()

	// Test matching headers
	headers := map[string]string{"Content-Type": "application/json"}
	response := createMockResponse(200, "", headers)

	monitorResponse := newHttpMonitorResponse()
	monitor.validateResponseHeaders(response, monitorResponse)
	assert.EqualValues(t, Success, monitorResponse.Base.Status)
	assert.Empty(t, monitorResponse.FailedAspects)

	// Test mismatched headers
	headers = map[string]string{"Content-Type": "text/html"}
	response = createMockResponse(200, "", headers)

	monitorResponse = newHttpMonitorResponse()
	monitor.validateResponseHeaders(response, monitorResponse)
	assert.EqualValues(t, Failure, monitorResponse.Base.Status)
	assert.Contains(t, monitorResponse.FailedAspects, headersAspect)
	assert.Contains(t, monitorResponse.Base.Failures[0], "Header mismatch")
}

func TestValidateResponseBody(t *testing.T) {
	monitor := setupTestHttpMonitor()

	// Test matching body
	response := createMockResponse(200, "success message", nil)

	monitorResponse := newHttpMonitorResponse()
	monitor.validateResponseBody(response, monitorResponse)
	assert.EqualValues(t, Success, monitorResponse.Base.Status)
	assert.Empty(t, monitorResponse.FailedAspects)

	// Test non-matching body
	response = createMockResponse(200, "error message", nil)

	monitorResponse = newHttpMonitorResponse()
	monitor.validateResponseBody(response, monitorResponse)
	assert.EqualValues(t, Failure, monitorResponse.Base.Status)
	assert.Contains(t, monitorResponse.FailedAspects, bodyAspect)
	assert.Contains(t, monitorResponse.Base.Failures[0], "Response body does not match regex")
}
