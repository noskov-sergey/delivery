package ports

import (
	"context"
	"database/sql"
	"delivery/internal/pkg/ddd"
)

type UnitOfWork interface {
	Tx() *sql.Tx
	Db() *sql.DB
	InTx() bool
	Track(aggregate ddd.AggregateRoot)

	Begin(ctx context.Context) error
	Commit(ctx context.Context) error

	CourierRepository() CourierRepository
	OrderRepository() OrderRepository
}
