package wrappers

import (
	"encoding/json"
	"net/http"

	"github.com/nahtann/go-lab/internal/utils"
)

type HandlerInterface[K interface{}, V interface{}] interface {
	Exec(*K) (*V, error)
}

type HttpWrapper[T interface{}, R interface{}] struct {
	Handler         HandlerInterface[T, R]
	ValidateRequest func(s any) string
}

func (wrapper *HttpWrapper[T, R]) Serve(w http.ResponseWriter, r *http.Request) error {
	test := new(T)

	err := json.NewDecoder(r.Body).Decode(&test)
	if err != nil {
		return utils.HttpServerInvalidRequest(w)
	}

	errorMessages := wrapper.ValidateRequest(test)
	if errorMessages != "" {
		message := utils.DefaultResponse{
			Message: errorMessages,
		}

		return utils.WriteJSON(w, http.StatusBadRequest, message)
	}

	response, err := wrapper.Handler.Exec(test)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, response)
}
