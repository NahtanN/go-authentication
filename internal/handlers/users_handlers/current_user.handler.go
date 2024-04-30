package users_handlers

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-authentication/internal/context_values"
	"github.com/nahtann/go-authentication/internal/storage/database/models"
	"github.com/nahtann/go-authentication/internal/utils"
)

type CurrentUserHttpHandler struct {
	DB *pgxpool.Pool
}

func NewCurrentUserHttpHandler(db *pgxpool.Pool) *CurrentUserHttpHandler {
	return &CurrentUserHttpHandler{
		DB: db,
	}
}

func (handler *CurrentUserHttpHandler) Serve(w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value(context_values.UserIdKey)

	if userId == nil {
		message := utils.DefaultResponse{
			Message: "Unable to fetch current user data.",
		}

		return utils.WriteJSON(w, http.StatusBadRequest, message)
	}

	user, err := CurrentUser(handler.DB, uint32(userId.(float64)))
	if err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	return utils.WriteJSON(w, http.StatusOK, user)
}

func CurrentUser(db *pgxpool.Pool, id uint32) (*models.UserModel, error) {
	user := models.UserModel{
		Id: id,
	}

	rows, err := db.Query(
		context.Background(),
		"SELECT username, email, created_at FROM users WHERE id = $1",
		id,
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
