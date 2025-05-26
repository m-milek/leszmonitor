package monitors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var httpClient HttpClient = newHttpClient()

func TestHttpMonitorRunSuccess(t *testing.T) {
	// Save the original client and restore it after the test
	originalClient := httpClient
	defer func() { httpClient = originalClient }()

	// Create a mock client
	mockClient := new(MockHTTPClient)
	httpClient = mockClient

	monitor := setupTestHttpMonitor()

	// Test successful response
	successResponse := createMockResponse(200, "success", map[string]string{
		"Content-Type": "application/json",
	})

	mockClient.On("Do", mock.Anything).Return(successResponse, nil).Once()

	response, err := monitor.Run(mockClient)

	if response, ok := response.(*HttpMonitorResponse); ok {
		assert.NoError(t, err)
		assert.EqualValues(t, Success, response.Base.Status)
		assert.Empty(t, response.FailedAspects)
	} else {
		t.Fatalf("Expected HttpMonitorResponse, got %T", response)
	}

	mockClient.AssertExpectations(t)
}

func TestHttpMonitorRunFailure(t *testing.T) {
	// Save the original client and restore it after the test
	originalClient := httpClient
	defer func() { httpClient = originalClient }()

	// Create a mock client
	mockClient := new(MockHTTPClient)
	httpClient = mockClient

	monitor := setupTestHttpMonitor()

	// Test failed response - wrong status code
	failedResponse := createMockResponse(404, "success", map[string]string{
		"Content-Type": "application/json",
	})

	mockClient.On("Do", mock.Anything).Return(failedResponse, nil).Once()

	response, err := monitor.Run(mockClient)

	if response, ok := response.(*HttpMonitorResponse); ok {
		assert.NoError(t, err)
		assert.EqualValues(t, Failure, response.Base.Status)
		assert.Contains(t, response.FailedAspects, StatusCodeAspect)
		assert.Len(t, response.FailedAspects, 1)
	} else {
		t.Fatalf("Expected HttpMonitorResponse, got %T", response)
	}

	mockClient.AssertExpectations(t)
}

func TestHttpMonitorRunError(t *testing.T) {
	// Save the original client and restore it after the test
	originalClient := httpClient
	defer func() { httpClient = originalClient }()

	// Create a mock client
	mockClient := new(MockHTTPClient)
	httpClient = mockClient

	monitor := setupTestHttpMonitor()

	// Test error response - HTTP client error
	mockClient.On("Do", mock.Anything).Return(nil, errors.New("connection refused")).Once()

	response, err := monitor.Run(mockClient)

	if response, ok := response.(*HttpMonitorResponse); ok {
		assert.Error(t, err)
		assert.EqualValues(t, Error, response.Base.Status)
		assert.Contains(t, response.Base.Errors[0], "connection refused")
	} else {
		t.Fatalf("Expected HttpMonitorResponse, got %T", response)
	}

	mockClient.AssertExpectations(t)
}

func TestHttpMonitorRunMultipleFailures(t *testing.T) {
	// Save the original client and restore it after the test
	originalClient := httpClient
	defer func() { httpClient = originalClient }()

	// Create a mock client
	mockClient := new(MockHTTPClient)
	httpClient = mockClient

	monitor := setupTestHttpMonitor()

	// Test response with multiple failures
	failedResponse := createMockResponse(404, "error", map[string]string{
		"Content-Type": "text/html",
	})

	mockClient.On("Do", mock.Anything).Return(failedResponse, nil).Once()

	response, err := monitor.Run(mockClient)

	if response, ok := response.(*HttpMonitorResponse); ok {
		assert.NoError(t, err)
		assert.EqualValues(t, Failure, response.Base.Status)
		assert.Contains(t, response.FailedAspects, StatusCodeAspect)
		assert.Contains(t, response.FailedAspects, BodyAspect)
		assert.Contains(t, response.FailedAspects, HeadersAspect)
		assert.Len(t, response.FailedAspects, 3)
	} else {
		t.Fatalf("Expected HttpMonitorResponse, got %T", response)
	}

	mockClient.AssertExpectations(t)
}
