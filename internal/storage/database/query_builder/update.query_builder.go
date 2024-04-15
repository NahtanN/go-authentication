package query_builder

import (
	"fmt"
	"strings"
)

type IQueryBuilderUpdate interface {
	Set(...QueryBuilderMethods) *QueryBuilder
	Where(...QueryBuilderMethods) *QueryBuilder
	Exec()
}

func (qb *QueryBuilder) Update(model ...TableModel) *QueryBuilder {
	tableData := qb.hasTableData()

	if !tableData && len(model) == 0 {
		qb.Errors = append(qb.Errors, "Query builder table data not set correctly.")
		return qb
	}

	qb.setTableData(model[0])

	query := fmt.Sprintf("UPDATE %s SET", qb.Table)

	qb.Query = query

	return qb
}

func (qb *QueryBuilder) Set(methods ...QueryBuilderMethods) *QueryBuilder {
	args := qb.setMethodsArgs(methods)

	setArgs := strings.Join(args, ", ")
	query := fmt.Sprintf("%s %s", qb.Query, setArgs)

	qb.Query = query

	return qb
}

func (qb *QueryBuilder) Where(methods ...QueryBuilderMethods) *QueryBuilder {
	args := qb.setMethodsArgs(methods)

	setArgs := strings.Join(args, ", ")
	query := fmt.Sprintf("%s WHERE %s", qb.Query, setArgs)

	qb.Query = query

	return qb
}
