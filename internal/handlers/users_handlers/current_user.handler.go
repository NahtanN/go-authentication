package users_handlers

import (
	"net/http"

	"github.com/nahtann/go-authentication/internal/context_values"
	"github.com/nahtann/go-authentication/internal/storage/database"
	"github.com/nahtann/go-authentication/internal/storage/database/models"
	"github.com/nahtann/go-authentication/internal/utils"
)

type CurrentUserHttpHandler struct {
	UserRepository database.UserRepository
}

func NewCurrentUserHttpHandler(userRepository database.UserRepository) *CurrentUserHttpHandler {
	return &CurrentUserHttpHandler{
		UserRepository: userRepository,
	}
}

func (handler *CurrentUserHttpHandler) Server(w http.ResponseWriter, r *http.Request) error {
	userId := r.Context().Value(context_values.UserIdKey)

	if userId == nil {
		message := utils.DefaultResponse{
			Message: "Unable to fetch current user data.",
		}

		return utils.WriteJSON(w, http.StatusBadRequest, message)
	}

	user, err := CurrentUser(handler.UserRepository, userId.(string))
	if err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	return utils.WriteJSON(w, http.StatusOK, user)
}

func CurrentUser(userRepository database.UserRepository, id string) (*models.UserModel, error) {
	user := models.UserModel{
		Id: id,
	}

	rows, err := userRepository.FindFirst(user).
		Select("username", "email", "created_at").
		Exec()
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
