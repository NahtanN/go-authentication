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
	Exec() (pgx.Rows, error)
}

type QueryBuilder struct {
	DB     *pgxpool.Pool
	Query  string
	Args   []any
	Errors []string
}

func (qb *QueryBuilder) Select(model interface{}, columns ...string) IQueryBuilder {
	if len(columns) <= 0 {
		return qb
	}

	modelType := reflect.TypeOf(model)

	if modelType.Kind() != reflect.Struct {
		qb.Errors = append(qb.Errors, "Select method with invalid model interface.")

		return qb
	}

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

func (qb *QueryBuilder) Exec() (pgx.Rows, error) {
	errors := strings.Join(qb.Errors, " ")

	if errors != "" {
		return nil, &utils.CustomError{
			Message: errors,
		}
	}

	rows, err := qb.DB.Query(context.Background(), qb.Query, qb.Args...)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to execute query.",
		}
	}

	return rows, nil
}
