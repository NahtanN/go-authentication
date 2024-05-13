package refresh_token

import (
	"context"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

type Request struct {
	Token string `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5..."`
}

type Handler struct {
	DB                     *pgxpool.Pool
	ValidateToken          func(token string) (*jwt.Token, bool)
	InvalidateTokensByUser func(userId uint32) error
	UpdateUserTokens       func(userId, parentTokenId uint32) (*auth_utils.Tokens, error)
}

func (handler *Handler) Exec(
	request *Request,
) (*auth_utils.Tokens, error) {
	// token, valid := middlewares.ValidateJWT(request.Token)
	token, valid := handler.ValidateToken(request.Token)

	if !valid || !token.Valid {
		return nil, &utils.CustomError{
			Message: "Refresh Token not valid.",
		}
	}

	rows, err := handler.DB.Query(
		context.Background(),
		"SELECT id, user_id, used FROM refresh_tokens WHERE token = $1",
		request.Token,
	)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to validate refresh token data.",
		}
	}
	defer rows.Close()

	var id, userId uint32
	var used bool

	for rows.Next() {
		err := rows.Scan(&id, &userId, &used)
		if err != nil {
			return nil, &utils.CustomError{
				Message: "Unable to parse refresh token data.",
			}
		}
	}

	if used || userId == 0 {
		_ = handler.InvalidateTokensByUser(userId)

		return nil, &utils.CustomError{
			Message: "Invalid Request",
		}
	}

	return handler.UpdateUserTokens(userId, id)
}
