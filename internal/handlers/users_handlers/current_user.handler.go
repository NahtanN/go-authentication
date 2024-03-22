package users_handlers

import (
	"fmt"
	"net/http"
)

type CurrentUserHttpHandler struct{}

func NewCurrentUserHttpHandler() *CurrentUserHttpHandler {
	return &CurrentUserHttpHandler{}
}

func (handler *CurrentUserHttpHandler) Server(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Curret user route")
	return nil
}
