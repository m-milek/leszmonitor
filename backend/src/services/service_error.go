package services

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
