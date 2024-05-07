package sign_in

import (
	"context"
	"fmt"

	"github.com/nahtann/go-lab/internal/interfaces"
	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

type Handler struct {
	DB              interfaces.Pgx
	VerifyPassword  func(password, hashedPassword string) (bool, error)
	CreateJwtTokens func(id uint32) (*auth_utils.Tokens, error)
}

type Request struct {
	User     string `json:"user"     validate:"required" example:"nahtann@outlook.com"`
	Password string `json:"password" validate:"required" example:"#Asdf123"`
}

func (handler *Handler) Exec(
	request *Request,
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
			fmt.Println(err)
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
			Message: "Unable to create access token.",
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
		return nil, &utils.CustomError{
			Message: "Internal error.",
		}
	}

	return tokens, nil
}
