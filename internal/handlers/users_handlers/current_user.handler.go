package users_handlers

import (
	"net/http"
)

type CurrentUserHttpHandler struct{}

func NewCurrentUserHttpHandler() *CurrentUserHttpHandler {
	return &CurrentUserHttpHandler{}
}

func (handler *CurrentUserHttpHandler) Server(w http.ResponseWriter, r *http.Request) error {
	// userId := r.Context().Value(middlewares.UserIdKey)

	return nil
}
