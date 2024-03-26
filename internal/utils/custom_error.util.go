package utils

type CustomError struct {
	Message string `json:"message"`
}

func (e *CustomError) Error() string {
	return e.Message
}
