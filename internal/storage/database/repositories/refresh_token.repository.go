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
	DB    *pgxpool.Pool
	table string
	alias string
}

func NewRefreshTokenRepository(database *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		DB:    database,
		table: "refresh_tokens",
		alias: "rft",
	}
}

func (r *RefreshTokenRepository) FindFirst() database.IQueryBuilder {
	queryBuilder := database.QueryBuilder{
		DB:    r.DB,
		Model: models.RefreshTokenModel{},
		Table: r.table,
		Alias: r.alias,
	}

	queryBuilder.Query = fmt.Sprintf("SELECT * FROM %s %s LIMIT 1", r.table, r.alias)

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

func (r *RefreshTokenRepository) Update(token models.RefreshTokenModel) error {
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

	// var rowId string
	// sequence := 1
	valueSequence := []string{}
	for index, field := range queryData.SearchFields {
		if field == "id" {
			// rowId = queryData.SearchArgs[index].(string)
			continue
		}

		formattedSequence := fmt.Sprintf("%s = $%d", field, index+1)
		valueSequence = append(valueSequence, formattedSequence)
		// sequence += 1
	}

	query := fmt.Sprintf(
		"UPDATE refresh_tokens SET %s WHERE id = $1",
		strings.Join(valueSequence, ", "),
	)

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
