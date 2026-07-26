package postgres

import (
	"context"
	"database/sql"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"testing"

	//"delivery/internal/adapters/out/postgres/courier_repo"
	//"delivery/internal/adapters/out/postgres/order_repo"
	//"delivery/internal/core/domain/model/courier"
	//"delivery/internal/core/domain/model/"
	//"delivery/internal/core/domain/model/order"
	//"delivery/internal/pkg/testcnts"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func Test_CourierRepositoryShouldAddCourier(t *testing.T) {
	assert := assert.New(t)

	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)
	assert.NotNil(t, db)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(err)
	assert.NotNil(uow)

	// Создаем курьера
	name := "test"
	speed := 5
	loc := kernel.RandomLocation()
	courier, err := courier.NewCourier(name, speed, loc)
	assert.NoError(err)

	//Добавляем курьеру место
	err = courier.AddStoragePlace("коробка", 7)
	assert.NoError(err)

	//Создаем заказ
	orderId := uuid.New()
	boxOrder, err := order.NewOrder(orderId, loc, 7)
	assert.NoError(err)

	//Курьер берет заказ
	err = courier.TakeOrder(boxOrder)
	assert.NoError(err)

	//Проверяем, что заказ взят
	assert.Equal(len(courier.StoragePlaces()), 2)
	assert.Equal(courier.StoragePlaces()[1].Name(), "коробка")

	//Вызываем repo Add
	err = uow.CourierRepository().Add(ctx, courier)
	assert.NoError(err)

	//Вызываем repo Get
	got, err := uow.CourierRepository().Get(ctx, courier.ID())
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(courier.ID(), got.ID())
	assert.Equal(courier.Name(), got.Name())
	assert.Equal(courier.Speed(), got.Speed())
	assert.Equal(courier.Location().X(), got.Location().X())
	assert.Equal(courier.Location().Y(), got.Location().Y())
	assert.Equal(len(courier.StoragePlaces()), len(got.StoragePlaces()))
}

func Test_CourierRepositoryShouldGetAllFreeCourier(t *testing.T) {
	assert := assert.New(t)

	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)
	assert.NotNil(t, db)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(err)
	assert.NotNil(uow)

	// Создаем курьера
	name := "free"
	speed := 5
	loc := kernel.RandomLocation()
	courierFree, err := courier.NewCourier(name, speed, loc)
	assert.NoError(err)

	// Создаем занятого курьера
	nameBusy := "busy"
	courierBusy, err := courier.NewCourier(nameBusy, speed, loc)
	assert.NoError(err)

	//Создаем заказ
	orderId := uuid.New()
	boxOrder, err := order.NewOrder(orderId, loc, 10)
	assert.NoError(err)

	//Курьер берет заказ
	err = courierBusy.TakeOrder(boxOrder)
	assert.NoError(err)

	//Вызываем repo Add
	err = uow.CourierRepository().Add(ctx, courierBusy)
	assert.NoError(err)

	//Вызываем repo Add
	err = uow.CourierRepository().Add(ctx, courierFree)
	assert.NoError(err)

	//Вызываем repo Get
	got, err := uow.CourierRepository().GetAllFree(ctx)
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(1, len(got))
	assert.Equal(courierFree.ID(), got[0].ID())
	assert.Equal(courierFree.Name(), got[0].Name())
	assert.Equal(courierFree.Speed(), got[0].Speed())
	assert.Equal(courierFree.Location().X(), got[0].Location().X())
	assert.Equal(courierFree.Location().Y(), got[0].Location().Y())
	assert.Equal(len(courierFree.StoragePlaces()), len(got[0].StoragePlaces()))
}

func Test_OrderRepositoryShouldAddOrder(t *testing.T) {
	assert := assert.New(t)

	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)
	assert.NotNil(t, db)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(err)
	assert.NotNil(uow)

	// Создаем курьера
	name := "test"
	speed := 5
	loc := kernel.RandomLocation()
	courier, err := courier.NewCourier(name, speed, loc)
	assert.NoError(err)

	//Создаем заказ
	orderId := uuid.New()
	order, err := order.NewOrder(orderId, loc, 7)
	assert.NoError(err)

	order.Assign(courier.ID())

	//Вызываем repo Add
	err = uow.OrderRepository().Add(ctx, order)
	assert.NoError(err)

	//Вызываем repo Get
	got, err := uow.OrderRepository().Get(ctx, order.ID())
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(order.ID(), got.ID())
	assert.Equal(order.Location(), got.Location())
	assert.Equal(order.Status(), got.Status())
	assert.Equal(order.Volume(), got.Volume())
	assert.Equal(order.CourierID(), got.CourierID())
}

func Test_OrderRepositoryShouldGetAssigned(t *testing.T) {
	assert := assert.New(t)

	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)
	assert.NotNil(t, db)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(err)
	assert.NotNil(uow)

	// Создаем новый заказ
	orderCreated, err := order.NewOrder(uuid.New(), kernel.RandomLocation(), 7)
	assert.NoError(err)

	//Вызываем repo Add
	err = uow.OrderRepository().Add(ctx, orderCreated)
	assert.NoError(err)

	// Создаем курьера
	name := "test"
	speed := 5
	loc := kernel.RandomLocation()
	courier, err := courier.NewCourier(name, speed, loc)
	assert.NoError(err)

	//Создаем заказ
	orderId := uuid.New()
	order, err := order.NewOrder(orderId, loc, 7)
	assert.NoError(err)

	order.Assign(courier.ID())

	//Вызываем repo Add
	err = uow.OrderRepository().Add(ctx, order)
	assert.NoError(err)

	//Вызываем repo Get
	got, err := uow.OrderRepository().GetAllAssignedStatus(ctx)
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(1, len(got))
	assert.Equal(order.ID(), got[0].ID())
	assert.Equal(order.Location(), got[0].Location())
	assert.Equal(order.Status(), got[0].Status())
	assert.Equal(order.Volume(), got[0].Volume())
	assert.Equal(order.CourierID(), got[0].CourierID())
}

func Test_OrderRepositoryShouldGetRandom(t *testing.T) {
	assert := assert.New(t)

	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)
	assert.NotNil(t, db)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(err)
	assert.NotNil(uow)

	// Создаем новый заказ
	orderCreated, err := order.NewOrder(uuid.New(), kernel.RandomLocation(), 7)
	assert.NoError(err)

	err = uow.OrderRepository().Add(ctx, orderCreated)
	assert.NoError(err)

	// Создаем курьера
	name := "test"
	speed := 5
	loc := kernel.RandomLocation()
	courier, err := courier.NewCourier(name, speed, loc)
	assert.NoError(err)

	//Создаем заказ
	orderId := uuid.New()
	order, err := order.NewOrder(orderId, loc, 7)
	assert.NoError(err)

	order.Assign(courier.ID())

	//Вызываем repo Add
	err = uow.OrderRepository().Add(ctx, order)
	assert.NoError(err)

	//Вызываем repo Get
	got, err := uow.OrderRepository().GetRandomCreatedStatus(ctx)
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(orderCreated.ID(), got.ID())
	assert.Equal(orderCreated.Location(), got.Location())
	assert.Equal(orderCreated.Status(), got.Status())
	assert.Equal(orderCreated.Volume(), got.Volume())
	assert.Equal(orderCreated.CourierID(), got.CourierID())
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

	err = goose.Up(db, "../../../../migrations")
	assert.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	return ctx, db, nil
}
