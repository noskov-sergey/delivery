package queries

import (
	postgre "delivery/internal/adapters/out/postgres"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetAllCouriersShouldWork(t *testing.T) {
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

	//Вызываем repo Add
	err = uow.CourierRepository().Add(ctx, courier1speed)
	assert.NoError(err)
	err = uow.CourierRepository().Add(ctx, courier4speed)
	assert.NoError(err)

	//Создаем handler
	handler, err := NewGetAllQueriersQueryHandler(uow)
	assert.NoError(err)

	query, err := NewGetAllQueriersQuery()
	assert.NoError(err)
	res, err := handler.Handle(ctx, query)
	assert.NoError(err)

	// Проверяем эквивалентность
	assert.Equal(2, len(res.Couriers))
}
