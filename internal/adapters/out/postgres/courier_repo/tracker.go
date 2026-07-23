package courier_repo

import (
	"context"
	"database/sql"
	"delivery/internal/pkg/ddd"
)

type Tracker interface {
	Tx() *sql.Tx
	Db() *sql.DB
	InTx() bool
	Track(agg ddd.AggregateRoot)
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
}
