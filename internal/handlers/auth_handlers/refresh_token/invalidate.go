package refresh_token

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-lab/internal/utils"
)

type InvalidateHandler struct {
	DB *pgxpool.Pool
}

func (handler *InvalidateHandler) TokensByUser(
	userId uint32,
) error {
	_, err := handler.DB.Exec(
		context.Background(),
		"UPDATE refresh_tokens SET used = true WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return &utils.CustomError{
			Message: "Invalid Request",
		}
	}

	return nil
}
