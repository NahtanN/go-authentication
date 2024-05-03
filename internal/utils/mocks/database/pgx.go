package mock_database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type MockDB struct {
	QueryFunc func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

func (m *MockDB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, args...)
	}
	return nil, nil
}
