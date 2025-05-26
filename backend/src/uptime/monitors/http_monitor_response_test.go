package monitors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHttpMonitorResponse(t *testing.T) {
	response := newHttpMonitorResponse()

	assert.Equal(t, Success, response.Base.Status)
	assert.Empty(t, response.Base.Errors)
	assert.Empty(t, response.Base.Failures)
	assert.Empty(t, response.FailedAspects)
	assert.Nil(t, response.RawHttpResponse)
}

func TestHttpMonitorResponseSetStatus(t *testing.T) {
	response := newHttpMonitorResponse()

	// Test setting status to Failure
	response.setStatus(Failure)
	assert.EqualValues(t, Failure, response.Base.Status)

	// Test setting status to Error
	response.setStatus(Error)
	assert.EqualValues(t, Error, response.Base.Status)
}

func TestHttpMonitorResponseAddErrorMsg(t *testing.T) {
	response := newHttpMonitorResponse()

	// Add an error message
	response.addErrorMsg("Test error")
	assert.Contains(t, response.Base.Errors, "Test error")
	assert.Len(t, response.Base.Errors, 1)

	// Add another error message
	response.addErrorMsg("Another error")
	assert.Contains(t, response.Base.Errors, "Another error")
	assert.Len(t, response.Base.Errors, 2)
}

func TestHttpMonitorResponseAddFailureMsg(t *testing.T) {
	response := newHttpMonitorResponse()

	// Add a failure message
	response.addFailureMsg("Test failure")
	assert.Contains(t, response.Base.Failures, "Test failure")
	assert.Len(t, response.Base.Failures, 1)

	// Add another failure message
	response.addFailureMsg("Another failure")
	assert.Contains(t, response.Base.Failures, "Another failure")
	assert.Len(t, response.Base.Failures, 2)
}

func TestHttpMonitorResponseAddFailedAspect(t *testing.T) {
	response := newHttpMonitorResponse()

	// Add one aspect
	response.addFailedAspect(StatusCodeAspect)
	assert.EqualValues(t, Failure, response.Base.Status)
	assert.Contains(t, response.FailedAspects, StatusCodeAspect)
	assert.Len(t, response.FailedAspects, 1)

	// Add another aspect
	response.addFailedAspect(BodyAspect)
	assert.EqualValues(t, Failure, response.Base.Status)
	assert.Contains(t, response.FailedAspects, StatusCodeAspect)
	assert.Contains(t, response.FailedAspects, BodyAspect)
	assert.Len(t, response.FailedAspects, 2)

	// Set status to Error and add another aspect
	response.setStatus(Error)
	response.addFailedAspect(HeadersAspect)
	assert.EqualValues(t, Error, response.Base.Status) // Status should remain Error
	assert.Len(t, response.FailedAspects, 3)
}
