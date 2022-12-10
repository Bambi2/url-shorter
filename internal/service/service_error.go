package service

type ServiceError struct {
	StatusCode int
	Msg        string
}

func (se *ServiceError) Error() string {
	return se.Msg
}
