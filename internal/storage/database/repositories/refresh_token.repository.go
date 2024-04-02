package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-authentication/internal/storage/database"
	database_common "github.com/nahtann/go-authentication/internal/storage/database/common"
	"github.com/nahtann/go-authentication/internal/storage/database/models"
	"github.com/nahtann/go-authentication/internal/utils"
)

type RefreshTokenRepository struct {
	DB *pgxpool.Pool
}

func NewRefreshTokenRepository(database *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		DB: database,
	}
}

func (r *RefreshTokenRepository) FindFirst(
	refreshToken models.RefreshTokenModel,
) database.IQueryBuilder {
	queryBuilder := database.QueryBuilder{
		DB:    r.DB,
		Model: refreshToken,
	}

	queryData, err := database_common.SetQueryData(refreshToken)
	if err != nil {
		queryBuilder.Errors = append(queryBuilder.Errors, "Unable to set query data.")
		return &queryBuilder
	}

	if len(queryData.SearchFields) == 0 || len(queryData.SearchArgs) == 0 {
		queryBuilder.Errors = append(queryBuilder.Errors, "Query data args not seted.")
		return &queryBuilder
	}

	clause := []string{}
	for i, v := range queryData.SearchFields {
		search := fmt.Sprintf("%s = $%d", v, i+1)

		clause = append(clause, search)
	}

	where := strings.Join(clause, " OR ")

	query := fmt.Sprintf("SELECT * FROM refresh_tokens WHERE %s", where)

	queryBuilder.Query = query
	queryBuilder.Args = queryData.SearchArgs

	return &queryBuilder
}

func (r *RefreshTokenRepository) Create(token models.RefreshTokenModel) error {
	queryData, err := database_common.SetQueryData(token)
	if err != nil {
		return &utils.CustomError{
			Message: "Unable to set query data.",
		}
	}

	if len(queryData.SearchFields) == 0 || len(queryData.SearchArgs) == 0 {
		return &utils.CustomError{
			Message: "Query data args not seted.",
		}
	}

	valueSequence := []string{}
	for index := range queryData.SearchArgs {
		formattedSequence := fmt.Sprintf("$%d", index+1)
		valueSequence = append(valueSequence, formattedSequence)
	}

	query := fmt.Sprintf(
		"INSERT INTO refresh_tokens (%s) VALUES (%s)",
		strings.Join(queryData.SearchFields, ", "),
		strings.Join(valueSequence, ", "),
	)

	fmt.Println(query, queryData.SearchArgs)

	_, err = r.DB.Exec(
		context.Background(),
		query,
		queryData.SearchArgs...,
	)
	if err != nil {
		fmt.Println(err)
		return &utils.CustomError{
			Message: "Unable to save refresh token.",
		}
	}

	return nil
}
