package services

import (
	"fmt"
	"net/http"
)

// ServiceError represents an error that occurred during a service operation.
// It includes an HTTP status code and the underlying error.
type ServiceError struct {
	Code int
	Err  error
}

// Error implements the error interface for ServiceError.
func (e ServiceError) Error() string {
	return e.Err.Error()
}

// NewNotFoundError creates a 404 Not Found ServiceError.
// It uses fmt.Errorf under the hood, so you can use %w to wrap existing errors.
func NewNotFoundError(format string, args ...any) *ServiceError {
	return &ServiceError{
		Code: http.StatusNotFound,
		Err:  fmt.Errorf(format, args...),
	}
}

// NewInternalError creates a 500 Internal Server Error ServiceError.
// It uses fmt.Errorf under the hood, so you can use %w to wrap existing errors.
func NewInternalError(format string, args ...any) *ServiceError {
	return &ServiceError{
		Code: http.StatusInternalServerError,
		Err:  fmt.Errorf(format, args...),
	}
}

// NewBadRequestError creates a 400 Bad Request ServiceError.
// It uses fmt.Errorf under the hood, so you can use %w to wrap existing errors.
func NewBadRequestError(format string, args ...any) *ServiceError {
	return &ServiceError{
		Code: http.StatusBadRequest,
		Err:  fmt.Errorf(format, args...),
	}
}

// NewForbiddenError creates a 403 Forbidden ServiceError.
// It uses fmt.Errorf under the hood, so you can use %w to wrap existing errors.
func NewForbiddenError(format string, args ...any) *ServiceError {
	return &ServiceError{
		Code: http.StatusForbidden,
		Err:  fmt.Errorf(format, args...),
	}
}

// NewUnauthorizedError creates a 401 Unauthorized ServiceError.
// It uses fmt.Errorf under the hood, so you can use %w to wrap existing errors.
func NewUnauthorizedError(format string, args ...any) *ServiceError {
	return &ServiceError{
		Code: http.StatusUnauthorized,
		Err:  fmt.Errorf(format, args...),
	}
}

// NewConflictError creates a 409 Conflict ServiceError.
// It uses fmt.Errorf under the hood, so you can use %w to wrap existing errors.
func NewConflictError(format string, args ...any) *ServiceError {
	return &ServiceError{
		Code: http.StatusConflict,
		Err:  fmt.Errorf(format, args...),
	}
}
