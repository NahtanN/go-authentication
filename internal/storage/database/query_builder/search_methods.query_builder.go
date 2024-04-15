package query_builder

import (
	"fmt"
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

func (e *EqualsMethod) Format(qb *QueryBuilder) string {
	valid := qb.ValidColumnData(e.Column, e.Value)

	if !valid {
		return ""
	}

	search := fmt.Sprintf("%s = $%d", e.Column, len(qb.Args)+1)

	qb.Args = append(qb.Args, e.Value)

	return search
}
