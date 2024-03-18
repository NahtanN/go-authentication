package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-authentication/internal/storage/database/models"
	"github.com/nahtann/go-authentication/internal/utils"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(database *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: database,
	}
}

func (r *UserRepository) Create(user *models.UserModel) error {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)"

	_, err := r.DB.Exec(
		context.Background(),
		query,
		user.Username, user.Email, user.Password,
	)
	if err != nil {
		return &utils.CustomError{
			Message: "Unable to create user.",
		}
	}

	return nil
}

func (r *UserRepository) UserExistsByColumn(
	column, value string,
) (bool, error) {
	valid := utils.ModelHasColumn(models.UserModel{}, column)

	if !valid {
		return false, &utils.CustomError{
			Message: "Model invalid.",
		}
	}

	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(%s) LIKE LOWER($1))",
		column,
	)

	var exists bool

	err := r.DB.QueryRow(context.Background(), query, value).Scan(&exists)
	if err != nil {
		return false, &utils.CustomError{
			Message: "Unable to validate user.",
		}
	}

	return exists, nil
}
