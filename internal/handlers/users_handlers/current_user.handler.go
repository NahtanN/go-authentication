package users_handlers

import (
	"fmt"
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

	CurrentUser(handler.UserRepository, userId.(string))

	return nil
}

func CurrentUser(userRepository database.UserRepository, id string) {
	user := models.UserModel{}

	rows, err := userRepository.FindFirst(&models.UserModel{
		Id: id,
	}).Select(models.UserModel{}, "id", "username", "email", "created_at").Exec()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			fmt.Println(err)
			/*return "", &utils.CustomError{*/
			/*Message: "Unable to parse data.",*/
			/*}*/
		}
	}

	fmt.Println(user)
}
