package monitors

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHttpMonitor_Validate(t *testing.T) {
	// Test valid monitor
	t.Run("Valid HTTP Monitor", func(t *testing.T) {
		baseMonitor := createTestBaseMonitor()
		config, _ := NewHttpConfig("GET", "https://example.com", nil, "", []int{200}, "", nil, 1000)
		monitor := HttpMonitor{
			BaseMonitor: baseMonitor,
			Config:      *config,
		}
		err := monitor.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid HTTP Monitor Config", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		err := monitor.validate()
		assert.NoError(t, err)
	})

	// Test empty URL
	t.Run("Empty URL", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		monitor.Url = ""
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "URL cannot be empty")
	})

	// Test empty HTTP method
	t.Run("Empty HTTP Method", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		monitor.HttpMethod = ""
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP method cannot be empty")
	})

	// Test invalid HTTP method
	t.Run("Invalid HTTP Method", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		monitor.HttpMethod = "INVALID"
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid HTTP method")
	})

	// Test empty expected status codes
	t.Run("Empty Expected Status Codes", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		emptyCodes := []int{}
		monitor.ExpectedStatusCodes = emptyCodes
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected status codes cannot be empty")
	})

	// Test invalid status code range (too low)
	t.Run("Status Code Too Low", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		invalidCodes := []int{50, 200}
		monitor.ExpectedStatusCodes = invalidCodes
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected status codes must be between 100 and 599")
	})

	// Test invalid status code range (too high)
	t.Run("Status Code Too High", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		invalidCodes := []int{200, 600}
		monitor.ExpectedStatusCodes = invalidCodes
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected status codes must be between 100 and 599")
	})

	// Test invalid URL format
	t.Run("Invalid URL Format", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		monitor.Url = "not-a-valid-url"
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid URL format")
	})

	// Test invalid URL scheme
	t.Run("Invalid URL Scheme", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		monitor.Url = "ftp://example.com"
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "URL scheme must be either http or https")
	})

	// Test negative expected response time
	t.Run("Negative Expected Response Time", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		negativeTime := -500
		monitor.ExpectedResponseTime = &negativeTime
		err := monitor.validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected response time cannot be negative")
	})

	// Test valid status codes
	t.Run("Valid Status Codes", func(t *testing.T) {
		monitor := setupTestHttpMonitorConfig()
		validCodes := []int{200, 201, 204}
		monitor.ExpectedStatusCodes = validCodes
		err := monitor.validate()
		assert.NoError(t, err)
	})

	// Test all valid HTTP methods
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}
	for _, method := range validMethods {
		t.Run("Valid HTTP Method: "+method, func(t *testing.T) {
			monitor := setupTestHttpMonitorConfig()
			monitor.HttpMethod = method
			err := monitor.validate()
			assert.NoError(t, err)
		})
	}
}

func TestValidateStatusCode(t *testing.T) {
	setupTest := func() (*HttpConfig, *http.Response, *HttpMonitorResponse) {
		monitor := setupTestHttpMonitorConfig()
		response := createMockResponse(200, "", nil)
		monitorResponse := NewHttpMonitorResponse()
		return monitor, response, monitorResponse
	}

	t.Run("Matching Status Code", func(t *testing.T) {
		monitor, response, monitorResponse := setupTest()

		monitor.checkStatusCode(response, monitorResponse)

		assert.Equal(t, Success, monitorResponse.Status)
		assert.Empty(t, monitorResponse.FailedAspects)
	})

	t.Run("Mismatched Status Code", func(t *testing.T) {
		monitor, response, monitorResponse := setupTest()
		monitor.ExpectedStatusCodes = []int{404}

		monitor.checkStatusCode(response, monitorResponse)

		assert.EqualValues(t, Failure, monitorResponse.Status)
		assert.Contains(t, monitorResponse.FailedAspects, statusCodeAspect)
		assert.Contains(t, monitorResponse.Failures[0], "Unexpected status code")
	})
}

func TestValidateResponseTime(t *testing.T) {
	setupTest := func() (*HttpConfig, *HttpMonitorResponse) {
		monitor := setupTestHttpMonitorConfig()
		monitorResponse := NewHttpMonitorResponse()
		return monitor, monitorResponse
	}

	t.Run("Response Time Within Limit", func(t *testing.T) {
		monitor, monitorResponse := setupTest()
		elapsed := 500 * time.Millisecond

		monitor.checkResponseTime(elapsed, monitorResponse)

		assert.EqualValues(t, Success, monitorResponse.Status)
		assert.Empty(t, monitorResponse.FailedAspects)
	})

	t.Run("Response Time Exceeding Limit", func(t *testing.T) {
		monitor, monitorResponse := setupTest()
		elapsed := 1500 * time.Millisecond

		monitor.checkResponseTime(elapsed, monitorResponse)

		assert.EqualValues(t, Failure, monitorResponse.Status)
		assert.Contains(t, monitorResponse.FailedAspects, responseTimeAspect)
		assert.Contains(t, monitorResponse.Failures[0], "Response time exceeded")
	})
}

func TestValidateResponseHeaders(t *testing.T) {
	setupTest := func(headers map[string]string) (*HttpConfig, *http.Response, *HttpMonitorResponse) {
		monitor := setupTestHttpMonitorConfig()
		response := createMockResponse(200, "", headers)
		monitorResponse := NewHttpMonitorResponse()
		return monitor, response, monitorResponse
	}

	t.Run("Matching Headers", func(t *testing.T) {
		headers := map[string]string{"Content-Type": "application/json"}
		monitor, response, monitorResponse := setupTest(headers)

		monitor.checkResponseHeaders(response, monitorResponse)

		assert.EqualValues(t, Success, monitorResponse.Status)
		assert.Empty(t, monitorResponse.FailedAspects)
	})

	t.Run("Mismatched Headers", func(t *testing.T) {
		headers := map[string]string{"Content-Type": "text/html"}
		monitor, response, monitorResponse := setupTest(headers)

		monitor.checkResponseHeaders(response, monitorResponse)

		assert.EqualValues(t, Failure, monitorResponse.Status)
		assert.Contains(t, monitorResponse.FailedAspects, headersAspect)
		assert.Contains(t, monitorResponse.Failures[0], "Header mismatch")
	})
}

func TestValidateResponseBody(t *testing.T) {
	setupTest := func(body string) (*HttpConfig, *http.Response, *HttpMonitorResponse) {
		monitor := setupTestHttpMonitorConfig()
		response := createMockResponse(200, body, nil)
		monitorResponse := NewHttpMonitorResponse()
		return monitor, response, monitorResponse
	}

	t.Run("Matching Body", func(t *testing.T) {
		monitor, response, monitorResponse := setupTest("success message")

		monitor.checkResponseBody(response, monitorResponse)

		assert.EqualValues(t, Success, monitorResponse.Status)
		assert.Empty(t, monitorResponse.FailedAspects)
	})

	t.Run("Non-Matching Body", func(t *testing.T) {
		monitor, response, monitorResponse := setupTest("error message")

		monitor.checkResponseBody(response, monitorResponse)

		assert.EqualValues(t, Failure, monitorResponse.Status)
		assert.Contains(t, monitorResponse.FailedAspects, bodyAspect)
		assert.Contains(t, monitorResponse.Failures[0], "Response body does not match regex")
	})
}
