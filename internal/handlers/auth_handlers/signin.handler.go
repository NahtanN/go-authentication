package auth_handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nahtann/go-lab/internal/interfaces"
	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

type SignInHandler struct {
	DB              interfaces.Pgx
	VerifyPassword  func(password, hashedPassword string) (bool, error)
	CreateJwtTokens func(id uint32) (*auth_utils.Tokens, error)
}

type SingInHandlerInterface interface {
	Exec(request *SigninRequest) (*auth_utils.Tokens, error)
}

type SignInHttpHandler struct {
	SignInHandler
}

type SigninRequest struct {
	User     string `json:"user"     validate:"required" example:"nahtann@outlook.com"`
	Password string `json:"password" validate:"required" example:"#Asdf123"`
}

// @Description	Authenticate user and returns access and refresh tokens.
// @Tags			auth
// @Accept			json
// @Param			request	body	SigninRequest	true	"Request Body"
// @Produce		json
// @Success		201	{object}	Tokens
// @Failure		400	{object}	utils.CustomError	"Message: 'User or password invalid.'"
// @router			/auth/sign-in   [post]
func NewSignInHttpHandler(
	db interfaces.Pgx,
) *SignInHttpHandler {
	signInHandler := SignInHandler{
		DB:              db,
		VerifyPassword:  utils.VerifyPassword,
		CreateJwtTokens: auth_utils.CreateJwtTokens,
	}

	return &SignInHttpHandler{
		SignInHandler: signInHandler,
	}
}

func (h *SignInHttpHandler) Serve(w http.ResponseWriter, r *http.Request) error {
	request := new(SigninRequest)

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return utils.HttpServerInvalidRequest(w)
	}

	errorMessages := utils.Validate(request)
	if errorMessages != "" {
		message := utils.DefaultResponse{
			Message: errorMessages,
		}

		return utils.WriteJSON(w, http.StatusBadRequest, message)
	}

	tokens, err := h.Exec(request)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, tokens)
}

func (handler *SignInHandler) Exec(
	request *SigninRequest,
) (*auth_utils.Tokens, error) {
	rows, err := handler.DB.Query(
		context.Background(),
		"SELECT id, password FROM users WHERE email LIKE $1 OR username LIKE $1",
		request.User,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id uint32
	var password string

	for rows.Next() {
		err := rows.Scan(&id, &password)
		if err != nil {
			return nil, &utils.CustomError{
				Message: "Unable to parse data.",
			}
		}
	}

	defaultError := utils.CustomError{
		Message: "User or password invalid.",
	}

	if id == 0 {
		return nil, &defaultError
	}

	match, err := handler.VerifyPassword(request.Password, password)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to validate user.",
		}
	}
	if !match {
		return nil, &defaultError
	}

	tokens, err := handler.CreateJwtTokens(id)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to generate access token.",
		}
	}

	// Insert refresh token into database
	_, err = handler.DB.Exec(
		context.Background(),
		"INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)",
		tokens.RefreshToken,
		id,
		tokens.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
