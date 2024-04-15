package query_builder

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-authentication/internal/utils"
)

type TableModel interface {
	Table() (string, string)
}

type QueryBuilder struct {
	DB        *pgxpool.Pool
	Model     interface{}
	ModelName string
	Query     string
	Args      []any
	Errors    []string
	Table     string
	Alias     string
}

type QueryBuilderMethods interface {
	Format(queryBuilder *QueryBuilder) string
}

func NewQueryBuilder(db *pgxpool.Pool, tableModel ...TableModel) *QueryBuilder {
	qb := QueryBuilder{
		DB: db,
	}

	if len(tableModel) > 0 {
		qb.setTableData(tableModel[0])
	}

	return &qb
}

func (qb *QueryBuilder) hasTableData() bool {
	if qb.Table == "" || qb.Alias == "" {
		return false
	}

	return true
}

func (qb *QueryBuilder) setMethodsArgs(methods []QueryBuilderMethods) []string {
	args := []string{}

	for _, fn := range methods {
		result := fn.Format(qb)

		args = append(args, result)
	}

	return args
}

func (qb *QueryBuilder) ValidColumnData(column string, data any) bool {
	valid, _ := utils.ModelHasColumn(qb.Model, column)

	if !valid {
		error := fmt.Sprintf("Column `%s` does not exists on model `%s`", column, qb.ModelName)
		qb.Errors = append(qb.Errors, error)
	}

	return valid
}

func (qb *QueryBuilder) setTableData(model TableModel) *QueryBuilder {
	modelType := reflect.TypeOf(qb.Model)

	if modelType.Kind() != reflect.Struct {
		qb.Errors = append(qb.Errors, "Invalid model interface.")

		return qb
	}

	table, alias := model.Table()

	qb.Table = table
	qb.Alias = alias

	qb.Model = model
	qb.ModelName = modelType.Name()

	return qb
}
