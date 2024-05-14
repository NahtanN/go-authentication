package refresh_token

import (
	"context"

	"github.com/nahtann/go-lab/internal/interfaces"
	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

type UpdateHandler struct {
	DB           interfaces.Pgx
	CreateTokens func(id uint32) (*auth_utils.Tokens, error)
}

func (handler *UpdateHandler) UserTokens(
	userId, parentTokenId uint32,
) (*auth_utils.Tokens, error) {
	tokens, err := handler.CreateTokens(userId)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to create tokens.",
		}
	}

	tx, err := handler.DB.Begin(context.Background())
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to start transaction.",
		}
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO refresh_tokens (parent_token_id, token, user_id, expires_at) VALUES ($1, $2, $3, $4)",
		parentTokenId,
		tokens.RefreshToken,
		userId,
		tokens.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to save tokens.",
		}
	}

	_, err = tx.Exec(
		context.Background(),
		"UPDATE refresh_tokens SET used = true WHERE id = $1",
		parentTokenId,
	)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to update refresh token.",
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to commit transaction.",
		}
	}

	return tokens, nil
}
