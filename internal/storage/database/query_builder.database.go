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

type QueryBuilderWhereMethod interface {
	Format(queryBuilder *QueryBuilder) string
}

type IQueryBuilder interface {
	Where(methods ...QueryBuilderWhereMethod) IQueryBuilder
	Select(v ...string) IQueryBuilder
	Exec() (pgx.Rows, error)
}

type QueryBuilder struct {
	DB     *pgxpool.Pool
	Model  interface{}
	Query  string
	Args   []any
	Errors []string
	Table  string
	Alias  string
}

func (qb *QueryBuilder) Where(methods ...QueryBuilderWhereMethod) IQueryBuilder {
	searchArgs := []string{}
	for _, method := range methods {
		search := method.Format(qb)

		searchArgs = append(searchArgs, search)
	}

	queryParts := strings.Split(qb.Query, qb.Alias)

	newQuery := fmt.Sprintf(
		"%s%s WHERE %s%s",
		queryParts[0],
		qb.Alias,
		strings.Join(searchArgs, " AND "),
		queryParts[1],
	)

	qb.Query = newQuery

	return qb
}

func (qb *QueryBuilder) Select(columns ...string) IQueryBuilder {
	if len(columns) <= 0 {
		return qb
	}

	modelType := reflect.TypeOf(qb.Model)

	if modelType.Kind() != reflect.Struct {
		qb.Errors = append(qb.Errors, "Select method with invalid model interface.")

		return qb
	}

	fields := []string{}
	for _, column := range columns {
		valid, databaseColumn := utils.ModelHasColumn(qb.Model, column)

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

	fmt.Println(qb.Query, qb.Args)

	rows, err := qb.DB.Query(context.Background(), qb.Query, qb.Args...)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to execute query.",
		}
	}

	return rows, nil
}
