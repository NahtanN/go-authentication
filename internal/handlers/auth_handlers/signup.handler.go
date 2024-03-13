package auth_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-authentication/internal/utils"
)

type SignUpHttpHandler struct {
	database *pgxpool.Pool
}

type SignupRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewSignUpHttpHandler(database *pgxpool.Pool) *SignUpHttpHandler {
	return &SignUpHttpHandler{
		database: database,
	}
}

func (handler *SignUpHttpHandler) Serve(w http.ResponseWriter, r *http.Request) error {
	request := new(SignupRequest)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		apiError := utils.DefaultResponse{
			Message: "Invalid request",
		}

		return utils.WriteJSON(w, http.StatusBadRequest, apiError)

	}

	errorMessages := utils.Validate(request)
	if errorMessages != "" {
		message := utils.DefaultResponse{
			Message: errorMessages,
		}

		return utils.WriteJSON(w, http.StatusBadRequest, message)
	}

	err = SingUp(handler.database, request)
	if err != nil {
		apiError := utils.DefaultResponse{
			Message: "Unable to sign in",
		}

		return utils.WriteJSON(w, http.StatusBadRequest, apiError)
	}

	message := utils.DefaultResponse{
		Message: "Sign up successfully",
	}

	return utils.WriteJSON(w, http.StatusCreated, message)
}

func SingUp(database *pgxpool.Pool, request *SignupRequest) error {
	return nil
}
