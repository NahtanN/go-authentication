package database

import (
	"github.com/nahtann/go-authentication/internal/storage/database/models"
)

type UserRepository interface {
	Create(user *models.UserModel) error
	UserExistsByColumn(
		column, value string,
	) (bool, error)
	FindFirst(user models.UserModel) IQueryBuilder
}

type RefreshTokenRepository interface {
	Create(token models.RefreshTokenModel) error
	FindFirst() IQueryBuilder
	Update(token models.RefreshTokenModel) error
}
