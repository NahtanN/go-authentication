package repositories

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-authentication/internal/storage/database"
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

func (r *UserRepository) FindFirst(user models.UserModel) database.IQueryBuilder {
	queryBuilder := database.QueryBuilder{
		DB:    r.DB,
		Model: user,
	}

	userTypes := reflect.TypeOf(user)
	userValues := reflect.ValueOf(user)

	if userTypes.Kind() != reflect.Struct {
		queryBuilder.Errors = append(queryBuilder.Errors, "User model invalid.")

		return &queryBuilder
	}

	fieldCount := 1
	searchFields := []string{}
	for i := 0; i < userTypes.NumField(); i++ {
		modelField := userTypes.Field(i)

		field := modelField.Tag.Get("db")
		value := userValues.Field(i).Interface()

		if modelField.Type == reflect.TypeOf(time.Time{}) && value.(time.Time).IsZero() {
			continue
		}

		clause := fmt.Sprintf("%s = $%d", field, fieldCount)

		if field != "" && value != "" {
			searchFields = append(searchFields, clause)
			queryBuilder.Args = append(queryBuilder.Args, value)
			fieldCount++
		}
	}

	where := strings.Join(searchFields, " OR ")

	query := fmt.Sprintf("SELECT * FROM users WHERE %s", where)

	queryBuilder.Query = query

	return &queryBuilder
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
	valid, _ := utils.ModelHasColumn(models.UserModel{}, column)

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
