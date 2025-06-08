package monitors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHttpMonitorRunSuccess(t *testing.T) {
	mockHttpClient := &MockHTTPClient{}

	monitor := setupTestHttpMonitorConfig()

	// Test successful response
	successResponse := createMockResponse(200, "success", map[string]string{
		"Content-Type": "application/json",
	})

	mockHttpClient.On("Do", mock.Anything).Return(successResponse, nil).Once()
	httpClientOrMock = mockHttpClient

	response := monitor.run()

	if response, ok := response.(*HttpMonitorResponse); ok {
		assert.EqualValues(t, Success, response.Base.Status)
		assert.Empty(t, response.FailedAspects)
	} else {
		t.Fatalf("Expected HttpMonitorResponse, got %T", response)
	}

	mockHttpClient.AssertExpectations(t)
}

func TestHttpMonitorRunFailure(t *testing.T) {
	mockClient := new(MockHTTPClient)

	monitor := setupTestHttpMonitorConfig()

	// Test failed response - wrong status code
	failedResponse := createMockResponse(404, "success", map[string]string{
		"Content-Type": "application/json",
	})

	mockClient.On("Do", mock.Anything).Return(failedResponse, nil).Once()
	httpClientOrMock = mockClient

	response := monitor.run()

	if response, ok := response.(*HttpMonitorResponse); ok {
		assert.EqualValues(t, Failure, response.Base.Status)
		assert.Contains(t, response.FailedAspects, statusCodeAspect)
		assert.Len(t, response.FailedAspects, 1)
	} else {
		t.Fatalf("Expected HttpMonitorResponse, got %T", response)
	}

	mockClient.AssertExpectations(t)
}

func TestHttpMonitorRunError(t *testing.T) {
	mockClient := new(MockHTTPClient)

	monitor := setupTestHttpMonitorConfig()

	// Test error response - HTTP client error
	mockClient.On("Do", mock.Anything).Return(nil, errors.New("connection refused")).Once()
	httpClientOrMock = mockClient

	response := monitor.run()

	if response, ok := response.(*HttpMonitorResponse); ok {
		assert.EqualValues(t, Error, response.Base.Status)
		assert.Contains(t, response.Base.Errors[0], "connection refused")
	} else {
		t.Fatalf("Expected HttpMonitorResponse, got %T", response)
	}

	mockClient.AssertExpectations(t)
}

func TestHttpMonitorRunMultipleFailures(t *testing.T) {
	mockClient := new(MockHTTPClient)

	monitor := setupTestHttpMonitorConfig()

	// Test response with multiple failures
	failedResponse := createMockResponse(404, "error", map[string]string{
		"Content-Type": "text/html",
	})

	mockClient.On("Do", mock.Anything).Return(failedResponse, nil).Once()
	httpClientOrMock = mockClient

	response := monitor.run()

	if response, ok := response.(*HttpMonitorResponse); ok {
		assert.EqualValues(t, Failure, response.Base.Status)
		assert.Contains(t, response.FailedAspects, statusCodeAspect)
		assert.Contains(t, response.FailedAspects, bodyAspect)
		assert.Contains(t, response.FailedAspects, headersAspect)
		assert.Len(t, response.FailedAspects, 3)
	} else {
		t.Fatalf("Expected HttpMonitorResponse, got %T", response)
	}

	mockClient.AssertExpectations(t)
}
