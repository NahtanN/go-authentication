package query_builder

import (
	"fmt"

	"github.com/nahtann/go-authentication/internal/storage/database"
)

type EqualsMethod struct {
	Column string
	Value  any
}

func Equals(column string, value any) *EqualsMethod {
	return &EqualsMethod{
		Column: column,
		Value:  value,
	}
}

func (e *EqualsMethod) Format(queryBuilder *database.QueryBuilder) string {
	search := fmt.Sprintf("%s = $%d", e.Column, len(queryBuilder.Args)+1)

	queryBuilder.Args = append(queryBuilder.Args, e.Value)

	return search
}
