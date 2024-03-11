package utils

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
