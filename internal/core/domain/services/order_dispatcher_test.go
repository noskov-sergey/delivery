package services

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_DispatchCanTakeRealWinnerWithSpeed(t *testing.T) {
	valera := makeGoodUserWithParams("Валера", 1, 1, 1)
	jora := makeGoodUserWithParams("Жора", 2, 1, 1)

	uid := uuid.New()
	o := makeOrderWithParams(uid, 5, 5, 7)

	svc := NewOrderDispatchService()

	// Должен победить Жора, т.к они в одном месте, но Жора быстрее
	c, err := svc.Dispatch(&o, []*courier.Courier{&valera, &jora})

	assert.NoError(t, err)
	assert.Equal(t, jora.ID(), c.ID())
	assert.Equal(t, jora.ID(), *o.CourierID())
	assert.Equal(t, order.StatusAssigned, o.Status())
}

func Test_DispatchCanTakeRealWinnerWithVolume(t *testing.T) {
	valera := makeGoodUserWithParams("Валера", 1, 1, 1)
	jora := makeGoodUserWithParams("Жора", 2, 1, 1)

	uidPreload := uuid.New()
	oPreload := makeOrderWithParams(uidPreload, 5, 5, 7)

	uidGood := uuid.New()
	o := makeOrderWithParams(uidGood, 5, 5, 7)

	// Загружаем сумку Жоре, должен победить Валера хотя он медленнее
	_ = jora.TakeOrder(&oPreload)

	svc := NewOrderDispatchService()

	c, err := svc.Dispatch(&o, []*courier.Courier{&valera, &jora})

	assert.NoError(t, err)
	assert.Equal(t, valera.ID(), c.ID())
	assert.Equal(t, valera.ID(), *o.CourierID())
	assert.Equal(t, order.StatusAssigned, o.Status())
}

func Test_DispatchCantTakeRealWinner(t *testing.T) {
	valera := makeGoodUserWithParams("Валера", 1, 1, 1)

	uidPreload := uuid.New()
	oPreload := makeOrderWithParams(uidPreload, 5, 5, 7)

	uidGood := uuid.New()
	o := makeOrderWithParams(uidGood, 5, 5, 7)

	// сумка занята
	_ = valera.TakeOrder(&oPreload)

	svc := NewOrderDispatchService()

	c, err := svc.Dispatch(&o, []*courier.Courier{&valera})

	assert.ErrorIs(t, err, ErrCantFoundCourier)
	assert.Nil(t, c)
}

func Test_DispatchBadOrder(t *testing.T) {
	valera := makeGoodUserWithParams("Валера", 1, 1, 1)

	uidPreload := uuid.New()
	oPreload := makeOrderWithParams(uidPreload, 5, 5, 7)

	// сумка занята
	_ = valera.TakeOrder(&oPreload)

	svc := NewOrderDispatchService()

	c, err := svc.Dispatch(nil, []*courier.Courier{&valera})

	assert.ErrorIs(t, err, ErrOrderIsNil)
	assert.Nil(t, c)
}

func makeGoodUserWithParams(name string, speed int, x, y int) courier.Courier {
	location, _ := kernel.NewLocation(uint8(x), uint8(y))
	c, _ := courier.NewCourier(name, speed, location)

	return *c
}

func makeOrderWithParams(u uuid.UUID, x, y int, vol int) order.Order {
	location, _ := kernel.NewLocation(uint8(x), uint8(y))
	o, _ := order.NewOrder(u, location, vol)

	return *o
}
