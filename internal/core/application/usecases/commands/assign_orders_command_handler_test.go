package commands

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/services"
	"testing"

	postgre "delivery/internal/adapters/out/postgres"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_AssignOrdersShouldWork(t *testing.T) {
	assert := assert.New(t)

	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(err)
	assert.NotNil(t, db)

	// Создаем UnitOfWork
	uow, err := postgre.NewUnitOfWork(db)
	assert.NoError(err)
	assert.NotNil(uow)

	// Создаем dispatcher
	dispatcher := services.NewOrderDispatchService()

	//Создаем handler
	handler, err := NewAssignOrdersCommandHandler(uow, dispatcher)
	assert.NoError(err)

	// Создаем курьера
	name := "test1"
	speed := 3
	loc, err := kernel.NewLocation(9, 9)
	assert.NoError(err)

	courier1, err := courier.NewCourier(name, speed, loc)
	assert.NoError(err)

	// Создаем второго курьера
	name2 := "test2"
	speed2 := 3
	loc2, err := kernel.NewLocation(1, 1)
	assert.NoError(err)

	courier2, err := courier.NewCourier(name2, speed2, loc2)
	assert.NoError(err)

	//Создаем заказ №1
	order1Id := uuid.New()
	loc19, err := kernel.NewLocation(7, 7)
	order1, err := order.NewOrder(order1Id, loc19, 3)
	assert.NoError(err)

	//Вызываем repo Add
	err = uow.OrderRepository().Add(ctx, order1)
	assert.NoError(err)

	err = uow.CourierRepository().Add(ctx, courier1)
	assert.NoError(err)
	err = uow.CourierRepository().Add(ctx, courier2)
	assert.NoError(err)

	cmd, err := NewAssignOrdersCommand()
	assert.NoError(err)

	handler.Handle(ctx, cmd)

	//Вызываем repo Get
	courier1result, err := uow.CourierRepository().Get(ctx, courier1.ID())
	assert.NoError(err)

	order1res, err := uow.OrderRepository().Get(ctx, order1.ID())
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(courier1.ID(), courier1result.ID())
	assert.Equal(order1.ID(), order1res.ID())
	assert.Equal(courier1.ID(), *order1res.CourierID())
	assert.Equal(order.StatusAssigned, order1res.Status())
}
