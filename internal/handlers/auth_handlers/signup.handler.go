package auth_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-authentication/internal/utils"
)

type UserModel struct {
	Id       string `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

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
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	message := utils.DefaultResponse{
		Message: "Sign up successfully",
	}

	return utils.WriteJSON(w, http.StatusCreated, message)
}

func SingUp(database *pgxpool.Pool, request *SignupRequest) error {
	errorMessages := []string{}

	usernameExists, err := UserExistsByColumn(database, "username", request.Username)
	if err != nil {
		return err
	}

	if usernameExists {
		errorMessages = append(errorMessages, "Username already in use.")
	}

	emailExists, err := UserExistsByColumn(database, "email", request.Email)
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

	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)"
	_, err = database.Exec(
		context.Background(),
		query,
		request.Username, request.Email, request.Password,
	)
	if err != nil {
		return &utils.CustomError{
			Message: "Unable to create user.",
		}
	}

	return nil
}

func UserExistsByColumn(database *pgxpool.Pool, column, value string) (bool, error) {
	valid := utils.ModelHasColumn(UserModel{}, column)

	if !valid {
		return false, &utils.CustomError{
			Message: "Model invalid.",
		}
	}

	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(%s) LIKE LOWER($1))",
		column,
	)

	var exists bool

	err := database.QueryRow(context.Background(), query, value).Scan(&exists)
	if err != nil {
		return false, &utils.CustomError{
			Message: "Unable to validate user.",
		}
	}

	return exists, nil
}
