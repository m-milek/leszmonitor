package monitors

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpMonitor(t *testing.T) {
	tests := []struct {
		name                string
		httpMethod          string
		url                 string
		expectedStatusCodes []int
		expectError         bool
		errorMessage        string
	}{
		{
			name:                "Valid HTTP Monitor",
			httpMethod:          "GET",
			url:                 "https://example.com",
			expectedStatusCodes: []int{200},
			expectError:         false,
		},
		{
			name:                "Empty URL",
			httpMethod:          "GET",
			url:                 "",
			expectedStatusCodes: []int{200},
			expectError:         true,
			errorMessage:        "URL cannot be empty",
		},
		{
			name:                "Empty HTTP Method",
			httpMethod:          "",
			url:                 "https://example.com",
			expectedStatusCodes: []int{200},
			expectError:         true,
			errorMessage:        "HTTP method cannot be empty",
		},
		{
			name:                "Invalid Status Code",
			httpMethod:          "GET",
			url:                 "https://example.com",
			expectedStatusCodes: []int{600},
			expectError:         true,
			errorMessage:        "expected status codes must be between 100 and 599",
		},
		{
			name:                "Invalid URL Scheme",
			httpMethod:          "GET",
			url:                 "ftp://example.com",
			expectedStatusCodes: []int{200},
			expectError:         true,
			errorMessage:        "URL scheme must be either http or https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor, err := NewHttpMonitor(
				tt.httpMethod,
				tt.url,
				map[string]string{},
				"",
				tt.expectedStatusCodes,
				"",
				map[string]string{},
				1000,
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, monitor)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, monitor)
				assert.Equal(t, tt.httpMethod, monitor.HttpMethod)
				assert.Equal(t, tt.url, monitor.Url)
				assert.Equal(t, tt.expectedStatusCodes, monitor.ExpectedStatusCodes)
			}
		})
	}
}

func TestHttpMonitorValidate(t *testing.T) {
	monitor := setupTestHttpMonitor()
	err := monitor.validate()
	assert.NoError(t, err)

	// Test invalid URL
	monitor.Url = "invalid-url"
	err = monitor.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid URL format")

	// Reset and test invalid response time
	monitor = setupTestHttpMonitor()
	negativeTime := -100
	monitor.ExpectedResponseTime = &negativeTime
	err = monitor.validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected response time cannot be negative")
}

func TestCreateRequest(t *testing.T) {
	monitor := setupTestHttpMonitor()

	// Test basic request creation
	req, err := monitor.createRequest()
	assert.NoError(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://example.com", req.URL.String())
	assert.Equal(t, "application/json", req.Header.Get("Accept"))

	// Test with body
	monitor.Body = "test body"
	monitor.HttpMethod = "POST"
	req, err = monitor.createRequest()
	assert.NoError(t, err)
	assert.Equal(t, "POST", req.Method)

	// Read the body to verify
	bodyBytes, err := io.ReadAll(req.Body)
	assert.NoError(t, err)
	assert.Equal(t, "test body", string(bodyBytes))
}
