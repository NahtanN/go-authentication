package sign_in

import (
	"encoding/json"
	"net/http"

	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

type HandlerInterface interface {
	Exec(*SigninRequest) (*auth_utils.Tokens, error)
}

type HttpWrapper struct {
	Handler         HandlerInterface
	ValidateRequest func(s any) string
}

// @Description	Authenticate user and returns access and refresh tokens.
// @Tags			auth
// @Accept			json
// @Param			request	body	SigninRequest	true	"Request Body"
// @Produce		json
// @Success		201	{object}	Tokens
// @Failure		400	{object}	utils.CustomError	"Message: 'User or password invalid.'"
// @router			/auth/sign-in   [post]
func (wrapper *HttpWrapper) Serve(w http.ResponseWriter, r *http.Request) error {
	request := new(SigninRequest)

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return utils.HttpServerInvalidRequest(w)
	}

	errorMessages := wrapper.ValidateRequest(request)
	if errorMessages != "" {
		message := utils.DefaultResponse{
			Message: errorMessages,
		}

		return utils.WriteJSON(w, http.StatusBadRequest, message)
	}

	tokens, err := wrapper.Handler.Exec(request)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, tokens)
}
