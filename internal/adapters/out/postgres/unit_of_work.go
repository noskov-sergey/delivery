package postgres

import (
	"context"
	"database/sql"
	courierrepo "delivery/internal/adapters/out/postgres/courier_repo"
	orderrepo "delivery/internal/adapters/out/postgres/order_repo"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"errors"
	"fmt"

	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	ErrDBValueIsRequired = errors.New("db value is required")
)

type UnitOfWork struct {
	tx                *sql.Tx
	db                *sql.DB
	trackedAggregates []ddd.AggregateRoot
	courierRepository ports.CourierRepository
	orderRepository   ports.OrderRepository
}

func NewUnitOfWork(db *sql.DB) (ports.UnitOfWork, error) {
	var (
		uow = &UnitOfWork{db: db}
		err error
	)

	if db == nil {
		return nil, ErrDBValueIsRequired
	}

	uow.courierRepository, err = courierrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}

	uow.orderRepository, err = orderrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}

	return uow, nil
}

func (u *UnitOfWork) Tx() *sql.Tx { return u.tx }

func (u *UnitOfWork) Db() *sql.DB { return u.db }

func (u *UnitOfWork) InTx() bool { return u.tx != nil }

func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *UnitOfWork) CourierRepository() ports.CourierRepository {
	return u.courierRepository
}

func (u *UnitOfWork) OrderRepository() ports.OrderRepository {
	return u.orderRepository
}

func (u *UnitOfWork) Begin(ctx context.Context) (err error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	u.tx = tx

	return nil
}

func (u *UnitOfWork) Commit(ctx context.Context) (err error) {
	if u.tx == nil {
		return fmt.Errorf("cannot commit without transaction")
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, u.tx.Rollback())
			if err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
				log.Error(err)
			}
		}

		u.clearTx()
	}()

	if err = u.tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *UnitOfWork) clearTx() { u.tx, u.trackedAggregates = nil, nil }
