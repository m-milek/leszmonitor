package services

type ServiceError struct {
	Code int
	Err  error
}

func (e ServiceError) Error() string {
	return e.Err.Error()
}
