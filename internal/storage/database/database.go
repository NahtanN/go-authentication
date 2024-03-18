package database

import "github.com/nahtann/go-authentication/internal/storage/database/models"

type UserRepository interface {
	Create(user *models.UserModel) error
	UserExistsByColumn(
		column, value string,
	) (bool, error)
}
