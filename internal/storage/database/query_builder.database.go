package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-authentication/internal/utils"
)

type IQueryBuilder interface {
	Select(model interface{}, v ...string) IQueryBuilder
	Exec() (pgx.Row, error)
}

type QueryBuilder struct {
	DB     *pgxpool.Pool
	Query  string
	Args   []string
	Errors []string
}

func (qb *QueryBuilder) Select(model interface{}, columns ...string) IQueryBuilder {
	if len(columns) <= 0 {
		return qb
	}

	modelType := reflect.TypeOf(model)

	fields := []string{}
	for _, column := range columns {
		valid, databaseColumn := utils.ModelHasColumn(model, column)

		if valid {
			fields = append(fields, databaseColumn)
			continue
		}

		error := fmt.Sprintf("Column `%s` does not exists on model `%s`", column, modelType.Name())

		qb.Errors = append(qb.Errors, error)
	}

	selectArgs := strings.Join(fields, ", ")

	query := strings.Replace(qb.Query, "*", selectArgs, 1)

	qb.Query = query

	return qb
}

func (qb *QueryBuilder) Exec() (pgx.Row, error) {
	errors := strings.Join(qb.Errors, " ")

	if errors != "" {
		return nil, &utils.CustomError{
			Message: errors,
		}
	}

	return qb.DB.QueryRow(context.Background(), qb.Query, qb.Args), nil
}
