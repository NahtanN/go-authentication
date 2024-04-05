package query_builder

import (
	"fmt"
	"strings"

	"github.com/nahtann/go-authentication/internal/storage/database"
)

type LogicalOperator struct {
	Operator string
	Methods  []database.QueryBuilderWhereMethod
}

func Operator(op string, methods ...database.QueryBuilderWhereMethod) *LogicalOperator {
	return &LogicalOperator{
		Operator: op,
		Methods:  methods,
	}
}

func (lo *LogicalOperator) Format(queryBuilder *database.QueryBuilder) string {
	searchArgs := []string{}
	for _, method := range lo.Methods {
		search := method.Format(queryBuilder)

		searchArgs = append(searchArgs, search)
	}

	formattedOperator := fmt.Sprintf(" %s ", lo.Operator)
	joinedArgs := strings.Join(searchArgs, formattedOperator)

	resultString := fmt.Sprintf("(%s)", joinedArgs)

	return resultString
}
