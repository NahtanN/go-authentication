package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	database "github.com/nahtann/go-authentication/internal/storage/database/common"
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

func (r *RefreshTokenRepository) Create(token models.RefreshTokenModel) error {
	queryData, err := database.SetQueryData(token)
	if err != nil {
		return &utils.CustomError{
			Message: "Unable to set query data.",
		}
	}

	if len(queryData.SearchFields) == 0 || len(queryData.SearchArgs) == 0 {
		return &utils.CustomError{
			Message: "Query data args not set.",
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
		return &utils.CustomError{
			Message: "Unable to save refresh token.",
		}
	}

	return nil
}
