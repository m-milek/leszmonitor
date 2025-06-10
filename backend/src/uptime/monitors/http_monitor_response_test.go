package monitors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHttpMonitorResponse(t *testing.T) {
	response := NewHttpMonitorResponse()

	assert.Equal(t, Success, response.Status)
	assert.Empty(t, response.Errors)
	assert.Empty(t, response.Failures)
	assert.Empty(t, response.FailedAspects)
}

func TestHttpMonitorResponseSetStatus(t *testing.T) {
	response := NewHttpMonitorResponse()

	// Test setting status to Failure
	response.setStatus(Failure)
	assert.EqualValues(t, Failure, response.Status)

	// Test setting status to Error
	response.setStatus(Error)
	assert.EqualValues(t, Error, response.Status)
}

func TestHttpMonitorResponseAddErrorMsg(t *testing.T) {
	response := NewHttpMonitorResponse()

	// Add an error message
	response.addErrorMsg("Test error")
	assert.Contains(t, response.Errors, "Test error")
	assert.Len(t, response.Errors, 1)

	// Add another error message
	response.addErrorMsg("Another error")
	assert.Contains(t, response.Errors, "Another error")
	assert.Len(t, response.Errors, 2)
}

func TestHttpMonitorResponseAddFailureMsg(t *testing.T) {
	response := NewHttpMonitorResponse()

	// Add a failure message
	response.addFailureMsg("Test failure")
	assert.Contains(t, response.Failures, "Test failure")
	assert.Len(t, response.Failures, 1)

	// Add another failure message
	response.addFailureMsg("Another failure")
	assert.Contains(t, response.Failures, "Another failure")
	assert.Len(t, response.Failures, 2)
}

func TestHttpMonitorResponseAddFailedAspect(t *testing.T) {
	response := NewHttpMonitorResponse()

	// Add one aspect
	response.addFailedAspect(statusCodeAspect)
	assert.EqualValues(t, Failure, response.Status)
	assert.Contains(t, response.FailedAspects, statusCodeAspect)
	assert.Len(t, response.FailedAspects, 1)

	// Add another aspect
	response.addFailedAspect(bodyAspect)
	assert.EqualValues(t, Failure, response.Status)
	assert.Contains(t, response.FailedAspects, statusCodeAspect)
	assert.Contains(t, response.FailedAspects, bodyAspect)
	assert.Len(t, response.FailedAspects, 2)

	// Set status to Error and add another aspect
	response.setStatus(Error)
	response.addFailedAspect(headersAspect)
	assert.EqualValues(t, Error, response.Status) // Status should remain Error
	assert.Len(t, response.FailedAspects, 3)
}
