package monitors

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	_const "github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHttpProbeFromReader(t *testing.T) {
	jsonInput := `{
		"url": "http://example.com",
		"method": "GET",
		"expectedStatusCodes": [200]
	}`

	probe, err := ProbeFromJSON(jsonInput, _const.HTTPConfigType)

	require.NoError(t, err)
	assert.NotNil(t, probe)

	httpMonitor, ok := probe.(*HTTPProbe)
	assert.True(t, ok)
	assert.Equal(t, "http://example.com", httpMonitor.URL)
}

func TestHttpMonitorFromReaderInvalidJSON(t *testing.T) {
	jsonInput := `invalid json`

	monitor, err := ProbeFromJSON(jsonInput, _const.HTTPConfigType)

	require.Error(t, err)
	assert.Nil(t, monitor)
}

func TestHttpMonitorFromReaderMissingType(t *testing.T) {
	jsonInput := `{
		"url": "http://example.com"
	}`

	monitor, err := ProbeFromJSON(jsonInput, "")

	require.Error(t, err)
	assert.Nil(t, monitor)
}

func TestHttpConfigValidate(t *testing.T) {
	t.Run("Valid config", func(t *testing.T) {
		config := HTTPProbe{
			Method:              "GET",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
		}
		assert.NoError(t, config.Validate())
	})

	t.Run("Empty URL", func(t *testing.T) {
		config := HTTPProbe{
			Method:              "GET",
			URL:                 "",
			ExpectedStatusCodes: []int{200},
		}
		require.Error(t, config.Validate())
		assert.Contains(t, config.Validate().Error(), "URL cannot be empty")
	})

	t.Run("Invalid URL", func(t *testing.T) {
		config := HTTPProbe{
			Method:              "GET",
			URL:                 "invalid-url",
			ExpectedStatusCodes: []int{200},
		}
		require.Error(t, config.Validate())
		assert.Contains(t, config.Validate().Error(), "invalid URL")
	})

	t.Run("Empty Method", func(t *testing.T) {
		config := HTTPProbe{
			Method:              "",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
		}
		require.Error(t, config.Validate())
		assert.Contains(t, config.Validate().Error(), "method cannot be empty")
	})

	t.Run("Invalid Method", func(t *testing.T) {
		config := HTTPProbe{
			Method:              "INVALID",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
		}
		require.Error(t, config.Validate())
		assert.Contains(t, config.Validate().Error(), "invalid HTTP method")
	})

	t.Run("Invalid Body Regex", func(t *testing.T) {
		config := HTTPProbe{
			Method:              "GET",
			URL:                 "http://example.com",
			ExpectedStatusCodes: []int{200},
			ExpectedBodyRegex:   "[invalid-regex",
		}
		require.Error(t, config.Validate())
		assert.Contains(t, config.Validate().Error(), "invalid body regex")
	})
}

func TestHttpMonitorRunSuccess(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}

	probe := setupTestHTTPProbe()

	successResponse := createMockResponse(200, "success", map[string]string{
		"Content-Type": "application/json",
	})

	mockHTTPClient.On("Do", mock.Anything).Return(successResponse, nil).Once()
	httpClientOrMock = mockHTTPClient

	response := probe.Run(context.Background(), uuid.Nil)

	assert.True(t, response.GetIsSuccess())
	assert.Empty(t, response.GetErrorDetails().Errors)

	details, ok := response.GetDetails().(*monitorresult.HTTPResultDetails)
	assert.True(t, ok)
	assert.Equal(t, 200, details.StatusCode)
	assert.Empty(t, response.GetErrorDetails().Failures)

	mockHTTPClient.AssertExpectations(t)
}

func TestHttpMonitorRunFailure(t *testing.T) {
	mockClient := new(MockHTTPClient)

	probe := setupTestHTTPProbe()

	failedResponse := createMockResponse(404, "not found", map[string]string{
		"Content-Type": "text/plain",
	})

	mockClient.On("Do", mock.Anything).Return(failedResponse, nil).Once()
	httpClientOrMock = mockClient

	response := probe.Run(context.Background(), uuid.Nil)

	assert.False(t, response.GetIsSuccess())

	assert.Contains(t, response.GetErrorDetails().Failures[0], "Unexpected status code")

	mockClient.AssertExpectations(t)
}

func TestHttpMonitorRunError(t *testing.T) {
	mockClient := new(MockHTTPClient)

	probe := setupTestHTTPProbe()

	mockClient.On("Do", mock.Anything).Return(nil, errors.New("connection refused")).Once()
	httpClientOrMock = mockClient

	response := probe.Run(context.Background(), uuid.Nil)

	assert.False(t, response.GetIsSuccess())
	assert.NotEmpty(t, response.GetErrorDetails().Errors)
	assert.Contains(t, response.GetErrorDetails().Errors[0], "connection refused")

	mockClient.AssertExpectations(t)
}

func TestHttpMonitorRunMultipleFailures(t *testing.T) {
	mockClient := new(MockHTTPClient)

	probe := setupTestHTTPProbe()
	probe.ExpectedBodyRegex = "success"
	probe.ExpectedHeaders = map[string]string{"X-Test": "Value"}

	failedResponse := createMockResponse(404, "error", map[string]string{
		"Content-Type": "text/html",
	})

	mockClient.On("Do", mock.Anything).Return(failedResponse, nil).Once()
	httpClientOrMock = mockClient

	response := probe.Run(context.Background(), uuid.Nil)

	assert.False(t, response.GetIsSuccess())

	failures := response.GetErrorDetails().Failures
	assert.Len(t, failures, 3)

	hasStatusCode := false
	hasBody := false
	hasHeaders := false
	for _, f := range failures {
		if strings.Contains(f, "status code") {
			hasStatusCode = true
		}
		if strings.Contains(f, "body") {
			hasBody = true
		}
		if strings.Contains(f, "Header mismatch") {
			hasHeaders = true
		}
	}
	assert.True(t, hasStatusCode)
	assert.True(t, hasBody)
	assert.True(t, hasHeaders)

	mockClient.AssertExpectations(t)
}

// MockHTTPClient is a mock implementation of the httpClient interface.
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// createMockResponse creates a mock [http.Response] for testing.
func createMockResponse(statusCode int, body string, headers map[string]string) *http.Response {
	header := make(http.Header)
	for k, v := range headers {
		header.Add(k, v)
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     header,
	}
}

// setupTestHTTPProbe returns a default HTTPProbe for testing.
func setupTestHTTPProbe() *HTTPProbe {
	return &HTTPProbe{
		Method:              "GET",
		URL:                 "http://example.com",
		ExpectedStatusCodes: []int{200},
	}
}
