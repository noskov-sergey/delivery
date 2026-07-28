package commands

import (
	"context"
	"database/sql"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"testing"

	postgre "delivery/internal/adapters/out/postgres"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func Test_MoveCouriersShouldWork(t *testing.T) {
	assert := assert.New(t)

	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)
	assert.NotNil(t, db)

	// Создаем UnitOfWork
	uow, err := postgre.NewUnitOfWork(db)
	assert.NoError(err)
	assert.NotNil(uow)

	// Создаем новый заказ
	orderCreated, err := order.NewOrder(uuid.New(), kernel.RandomLocation(), 7)
	assert.NoError(err)

	err = uow.OrderRepository().Add(ctx, orderCreated)
	assert.NoError(err)

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

	courier1speed.TakeOrder(order1)
	courier4speed.TakeOrder(order2)

	// Создаем проверочные координаты
	loc1Should, err := kernel.NewLocation(8, 9)
	assert.NoError(err)

	loc2Should, err := kernel.NewLocation(2, 2)
	assert.NoError(err)

	//Вызываем repo Add
	err = uow.OrderRepository().Add(ctx, order1)
	assert.NoError(err)
	err = uow.OrderRepository().Add(ctx, order2)
	assert.NoError(err)

	err = uow.CourierRepository().Add(ctx, courier1speed)
	assert.NoError(err)
	err = uow.CourierRepository().Add(ctx, courier4speed)
	assert.NoError(err)

	//Создаем handler
	handler, err := NewMoveCouriersCommandHandler(uow)
	assert.NoError(err)

	cmd, err := NewMoveCouriersCommand()
	assert.NoError(err)
	handler.Handle(ctx, cmd)

	//Вызываем repo Get
	couirier1result, err := uow.CourierRepository().Get(ctx, courier1speed.ID())
	assert.NoError(err)

	couirier4result, err := uow.CourierRepository().Get(ctx, courier4speed.ID())
	assert.NoError(err)

	order2res, err := uow.OrderRepository().Get(ctx, order2.ID())
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(courier1speed.ID(), couirier1result.ID())
	assert.Equal(loc1Should, couirier1result.Location())

	assert.Equal(loc2Should, couirier4result.Location())

	assert.Equal(order.StatusCompleted, order2res.Status())
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
