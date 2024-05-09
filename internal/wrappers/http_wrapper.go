package wrappers

import (
	"net/http"

	"github.com/nahtann/go-lab/internal/utils"
)

type RequestParser[T interface{}] func(*T, *http.Request) error

type HandlerInterface[K interface{}, V interface{}] interface {
	Exec(*K) (*V, error)
}

// HttpWrapper expects two type definitions.
// R (the first one) should be of type Handler Request
// E (the second one) should be of type handler Exec return value
type HttpWrapper[R interface{}, E interface{}] struct {
	Handler         HandlerInterface[R, E]
	RequestParsers  []RequestParser[R]
	ValidateRequest func(s any) string
}

func (wrapper *HttpWrapper[R, E]) Serve(w http.ResponseWriter, r *http.Request) error {
	request := new(R)

	for _, parser := range wrapper.RequestParsers {
		err := parser(request, r)
		if err != nil {
			return utils.HttpServerInvalidRequest(w)
		}
	}

	if wrapper.ValidateRequest != nil {
		errorMessages := wrapper.ValidateRequest(request)
		if errorMessages != "" {
			message := utils.DefaultResponse{
				Message: errorMessages,
			}

			return utils.WriteJSON(w, http.StatusBadRequest, message)
		}
	}

	response, err := wrapper.Handler.Exec(request)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, response)
}
