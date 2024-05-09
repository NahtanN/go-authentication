package current_user

import (
	"context"
	"net/http"

	"github.com/nahtann/go-lab/internal/context_values"
	"github.com/nahtann/go-lab/internal/interfaces"
	"github.com/nahtann/go-lab/internal/storage/database/models"
	"github.com/nahtann/go-lab/internal/utils"
)

type Request struct {
	ID uint32 `validate:"required"`
}

func RequestParser(request *Request, req *http.Request) error {
	id := req.Context().Value(context_values.UserIdKey)

	if id == nil {
		return &utils.CustomError{
			Message: "Unable to fetch current user data.",
		}
	}

	request.ID = uint32(id.(float64))

	return nil
}

type Handler struct {
	DB interfaces.Pgx
}

func (handler *Handler) Exec(request *Request) (*models.UserModel, error) {
	user := models.UserModel{
		Id: request.ID,
	}

	rows, err := handler.DB.Query(
		context.Background(),
		"SELECT username, email, created_at FROM users WHERE id = $1",
		request.ID,
	)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to retrieve current user data.",
		}
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, &utils.CustomError{
				Message: "Unable to parse current user data.",
			}
		}
	}

	return &user, nil
}
