package interfaces

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Pgx interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
