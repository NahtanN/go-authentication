package auth_handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nahtann/go-authentication/internal/storage/database"
	"github.com/nahtann/go-authentication/internal/storage/database/models"
	"github.com/nahtann/go-authentication/internal/utils"
)

type SignUpHttpHandler struct {
	UserRepository database.UserRepository
}

type SignupRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewSignUpHttpHandler(userRepository database.UserRepository) *SignUpHttpHandler {
	return &SignUpHttpHandler{
		UserRepository: userRepository,
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

	err = SingUp(handler.UserRepository, request)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	message := utils.DefaultResponse{
		Message: "Sign up successfully",
	}

	return utils.WriteJSON(w, http.StatusCreated, message)
}

func SingUp(userRepository database.UserRepository, request *SignupRequest) error {
	errorMessages := []string{}

	usernameExists, err := userRepository.UserExistsByColumn("username", request.Username)
	if err != nil {
		return err
	}

	if usernameExists {
		errorMessages = append(errorMessages, "Username already in use.")
	}

	emailExists, err := userRepository.UserExistsByColumn("email", request.Email)
	if err != nil {
		return err
	}

	if emailExists {
		errorMessages = append(errorMessages, "E-mail already in use.")
	}

	if usernameExists || emailExists {
		return &utils.CustomError{
			Message: strings.Join(errorMessages, " "),
		}
	}

	err = userRepository.Create(&models.UserModel{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return err
	}

	return nil
}
