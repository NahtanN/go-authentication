package database

import (
	"github.com/nahtann/go-authentication/internal/storage/database/models"
)

/*type QueryBuilder interface {*/
/*Select(v ...string)*/
/*}*/

type UserRepository interface {
	Create(user *models.UserModel) error
	UserExistsByColumn(
		column, value string,
	) (bool, error)
	FindFirst(user *models.UserModel) IQueryBuilder
}
