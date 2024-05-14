package sign_up

import (
	"context"
	"fmt"
	"strings"

	"github.com/nahtann/go-lab/internal/interfaces"
	"github.com/nahtann/go-lab/internal/utils"
)

type Handler struct {
	DB           interfaces.Pgx
	HashPassword func(password string) (string, error)
}

type Request struct {
	Username string `json:"username" validate:"required" example:"NahtanN"`
	Email    string `json:"email"    validate:"required" example:"nahtann@outlook.com"`
	Password string `json:"password" validate:"required" example:"#Asdf123"`
}

func (handler *Handler) Exec(request *Request) (*utils.DefaultResponse, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(%s) LIKE LOWER($1))"

	errorMessages := []string{}

	var usernameExists, emailExists bool

	queryByUsername := fmt.Sprintf(query, "username")

	err := handler.DB.QueryRow(context.Background(), queryByUsername, request.Username).
		Scan(&usernameExists)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to validate user username.",
		}
	}
	if usernameExists {
		errorMessages = append(errorMessages, "Username already in use.")
	}

	queryByEmail := fmt.Sprintf(query, "email")

	err = handler.DB.QueryRow(context.Background(), queryByEmail, request.Email).
		Scan(&emailExists)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to validate user email.",
		}
	}
	if emailExists {
		errorMessages = append(errorMessages, "E-mail already in use.")
	}

	if usernameExists || emailExists {
		return nil, &utils.CustomError{
			Message: strings.Join(errorMessages, " "),
		}
	}

	hashPassword, err := handler.HashPassword(request.Password)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to validate password.",
		}
	}

	_, err = handler.DB.Exec(
		context.Background(),
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		request.Username, request.Email, hashPassword,
	)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to create user.",
		}
	}

	return &utils.DefaultResponse{
		Message: "Sign up successfully",
	}, nil
}
