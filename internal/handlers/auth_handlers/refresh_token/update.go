package refresh_token

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

type UpdateHandler struct {
	DB           *pgxpool.Pool
	CreateTokens func(id uint32) (*auth_utils.Tokens, error)
}

func (handler *UpdateHandler) UserTokens(
	userId, parentTokenId uint32,
) (*auth_utils.Tokens, error) {
	tokens, err := handler.CreateTokens(userId)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to generate access token.",
		}
	}

	_, err = handler.DB.Exec(
		context.Background(),
		"INSERT INTO refresh_tokens (parent_token_id, token, user_id, expires_at) VALUES ($1, $2, $3, $4)",
		parentTokenId,
		tokens.RefreshToken,
		userId,
		tokens.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	_, err = handler.DB.Exec(
		context.Background(),
		"UPDATE refresh_tokens SET used = true WHERE id = $1",
		parentTokenId,
	)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
