package utils

import (
	"net/http"
)

type DefaultResponse struct {
	Message string `json:"message"`
}

func HttpServerError(w http.ResponseWriter) error {
	message := DefaultResponse{
		Message: "Server Error",
	}

	return WriteJSON(w, http.StatusInternalServerError, message)
}

func HttpServerInvalidRequest(w http.ResponseWriter) error {
	message := DefaultResponse{
		Message: "Invalid Request",
	}

	return WriteJSON(w, http.StatusInternalServerError, message)
}
