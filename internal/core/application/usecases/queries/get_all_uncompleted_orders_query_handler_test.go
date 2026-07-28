package queries

import (
	"context"
	"database/sql"
	postgre "delivery/internal/adapters/out/postgres"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"testing"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func Test_GetUncompletedOrdersShouldWork(t *testing.T) {
	assert := assert.New(t)

	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)
	assert.NotNil(t, db)

	// Создаем UnitOfWork
	uow, err := postgre.NewUnitOfWork(db)
	assert.NoError(err)
	assert.NotNil(uow)

	// Создаем курьера
	name := "test1speed"
	speed := 1
	loc, err := kernel.NewLocation(9, 9)
	assert.NoError(err)

	courier1speed, err := courier.NewCourier(name, speed, loc)
	assert.NoError(err)

	// Создаем второго курьера
	name4 := "test4speed"
	speed4 := 4
	loc4, err := kernel.NewLocation(1, 1)
	assert.NoError(err)

	courier4speed, err := courier.NewCourier(name4, speed4, loc4)
	assert.NoError(err)

	//Создаем заказ №1
	order1Id := uuid.New()
	loc19, err := kernel.NewLocation(1, 9)
	order1, err := order.NewOrder(order1Id, loc19, 7)
	assert.NoError(err)

	//Создаем заказ №2
	order2Id := uuid.New()
	loc22, err := kernel.NewLocation(2, 2)
	order2, err := order.NewOrder(order2Id, loc22, 7)
	assert.NoError(err)

	// Присваиваем заказы
	order1.Assign(courier1speed.ID())
	order2.Assign(courier4speed.ID())

	err = order2.Complete()
	assert.NoError(err)

	//Вызываем repo Add
	err = uow.OrderRepository().Add(ctx, order1)
	assert.NoError(err)
	err = uow.OrderRepository().Add(ctx, order2)
	assert.NoError(err)

	//Создаем handler
	handler, err := NewGetAllUncompletedOrderQueryHandler(uow)
	assert.NoError(err)

	query, err := NewGetAllUncompletedOrdersQuery()
	assert.NoError(err)
	orders, err := handler.Handle(ctx, query)
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(1, len(orders.Orders))
	assert.Equal(order1.ID(), orders.Orders[0].ID)
}

func setupTest(t *testing.T) (context.Context, *sql.DB, error) {
	ctx := context.Background()
	dbName := "delivery_test"
	dbUser := "delivery_test"
	dbPassword := "delivery_test"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	assert.NoError(t, err)

	psqlInfo, err := postgresContainer.ConnectionString(context.Background(), "sslmode=disable")
	assert.NoError(t, err, "error creating connection string")

	println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	assert.NoError(t, err, "error connecting to postgres")

	err = goose.Up(db, "../../../../../migrations")
	assert.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	return ctx, db, nil
}
