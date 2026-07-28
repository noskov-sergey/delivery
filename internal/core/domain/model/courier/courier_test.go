package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CourierReturnWhenParamsAreCorrectOnCreated(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	courier, err := NewCourier(name, speed, location)

	assert.NoError(t, err)
	assert.Equal(t, name, courier.Name())
	assert.Equal(t, speed, courier.Speed())
	assert.Equal(t, location, courier.Location())
	assert.NotNil(t, courier.storagePlace[0])
}

func Test_CourierReturnWhenParamsAreInCorrectOnCreated(t *testing.T) {
	goodLocation, _ := kernel.NewLocation(1, 1)
	wrongLocation := kernel.Location{}

	tests := map[string]struct {
		name     string
		speed    int
		location kernel.Location

		expectErr error
	}{
		"empty name": {
			name:     "",
			speed:    3,
			location: goodLocation,

			expectErr: ErrCourierNameIsEmpty,
		},
		"empty speed": {
			name:     "Макс",
			speed:    0,
			location: goodLocation,

			expectErr: ErrCourierSpeedIsInvalid,
		},
		"empty location": {
			name:     "Макс",
			speed:    3,
			location: wrongLocation,

			expectErr: ErrCourierLocationIsEmpty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			couirer, err := NewCourier(test.name, test.speed, test.location)

			assert.Nil(t, couirer)
			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}

func Test_CourierAddStoragePlaceWhenParamsAreCorrect(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	courier, err := NewCourier(name, speed, location)

	assert.NoError(t, err)
	assert.Equal(t, name, courier.Name())
	assert.Equal(t, speed, courier.Speed())
	assert.Equal(t, location, courier.Location())
	assert.NotNil(t, courier.storagePlace[0])
	assert.Equal(t, 1, len(courier.storagePlace))

	err = courier.AddStoragePlace("котомка", 3)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(courier.storagePlace))
}

func Test_CourierAddStoragePlaceParamsAreInCorrect(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	valera, _ := NewCourier(name, speed, location)

	tests := map[string]struct {
		courier Courier
		name    string
		volume  int

		expectErr error
	}{
		"empty name": {
			courier: *valera,
			name:    "",
			volume:  3,

			expectErr: ErrStorageNameEmpty,
		},
		"empty volume": {
			courier: *valera,
			name:    "котомка",
			volume:  0,

			expectErr: ErrStorageVolumeEmpty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.courier.AddStoragePlace(test.name, test.volume)

			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}

func Test_CourierCanTakeOrderWhenParamsAreCorrect(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	goodOrder, _ := order.NewOrder(uuid.New(), location, 9)

	courier, err := NewCourier(name, speed, location)

	assert.NoError(t, err)
	assert.Equal(t, name, courier.Name())
	assert.Equal(t, speed, courier.Speed())
	assert.Equal(t, location, courier.Location())
	assert.NotNil(t, courier.storagePlace[0])
	assert.Equal(t, 1, len(courier.storagePlace))

	c, err := courier.CanTakeOrder(*goodOrder)
	assert.NoError(t, err)
	assert.True(t, c)
}

func Test_CourierCanTakeOrderParamsAreInCorrect(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	valera, _ := NewCourier(name, speed, location)

	// dima already has order
	dima, _ := NewCourier(name, speed, location)
	_ = dima.AddStoragePlace(name, 5)

	badOrder, _ := order.NewOrder(uuid.New(), location, 11)

	tests := map[string]struct {
		courier Courier
		order   order.Order

		expectErr error
	}{
		"too big order": {
			courier: *valera,
			order:   *badOrder,

			expectErr: ErrCourierCantGetVolume,
		},
		"already has order": {
			courier: *dima,
			order:   *badOrder,

			expectErr: ErrCourierCantGetVolume,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c, err := test.courier.CanTakeOrder(test.order)

			assert.ErrorIs(t, err, test.expectErr)
			assert.False(t, c)
		})
	}
}

func Test_CourierTakeOrderWhenParamsAreCorrect(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	goodOrder, _ := order.NewOrder(uuid.New(), location, 9)

	courier, err := NewCourier(name, speed, location)

	assert.NoError(t, err)
	assert.Equal(t, name, courier.Name())
	assert.Equal(t, speed, courier.Speed())
	assert.Equal(t, location, courier.Location())
	assert.NotNil(t, courier.storagePlace[0])
	assert.Equal(t, 1, len(courier.storagePlace))

	err = courier.TakeOrder(goodOrder)
	assert.NoError(t, err)

}

func Test_CourierTakeOrderParamsAreInCorrect(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	valera, _ := NewCourier(name, speed, location)

	// dima already has order
	dima, _ := NewCourier(name, speed, location)

	goodOrder, _ := order.NewOrder(uuid.New(), location, 4)
	badOrder, _ := order.NewOrder(uuid.New(), location, 11)

	err := dima.TakeOrder(goodOrder)
	assert.NoError(t, err)

	tests := map[string]struct {
		courier Courier
		order   order.Order

		expectErr error
	}{
		"too big order": {
			courier: *valera,
			order:   *badOrder,

			expectErr: ErrCourierCantTakeOrder,
		},
		"already has order": {
			courier: *dima,
			order:   *goodOrder,

			expectErr: ErrCourierCantTakeOrder,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.courier.TakeOrder(&test.order)

			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}

func Test_CourierCompleteOrderWhenParamsAreCorrect(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	goodOrder, _ := order.NewOrder(uuid.New(), location, 9)

	courier, err := NewCourier(name, speed, location)

	assert.NoError(t, err)
	assert.Equal(t, name, courier.Name())
	assert.Equal(t, speed, courier.Speed())
	assert.Equal(t, location, courier.Location())
	assert.NotNil(t, courier.storagePlace[0])
	assert.Equal(t, 1, len(courier.storagePlace))

	err = courier.TakeOrder(goodOrder)
	assert.NoError(t, err)

	err = goodOrder.Assign(courier.ID())
	assert.NoError(t, err)

	err = courier.CompleteOrder(goodOrder)
	assert.NoError(t, err)

	assert.Equal(t, order.StatusCompleted, goodOrder.Status())
}

func Test_CourierCompleteOrderParamsAreInCorrect(t *testing.T) {
	name := "Валера"
	speed := 3
	location, _ := kernel.NewLocation(2, 5)

	valera, _ := NewCourier(name, speed, location)

	goodOrder, _ := order.NewOrder(uuid.New(), location, 4)
	wrongOrder, _ := order.NewOrder(uuid.New(), location, 4)

	_ = valera.TakeOrder(goodOrder)

	tests := map[string]struct {
		courier Courier
		order   order.Order

		expectErr error
	}{
		"too big order": {
			courier: *valera,
			order:   *wrongOrder,

			expectErr: ErrCourierOrderNotExists,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.courier.CompleteOrder(&test.order)

			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}

func Test_CourierCalculateTimeToLocationWhenParamsAreCorrect(t *testing.T) {
	name := "Валера"
	speed := 2
	location, _ := kernel.NewLocation(1, 1)
	targetLocation, _ := kernel.NewLocation(5, 5)

	courier, err := NewCourier(name, speed, location)

	count, err := courier.CalculateTimeToLocation(targetLocation)
	assert.NoError(t, err)
	assert.Equal(t, float64(4), count)
}

func Test_CourierCalculateTimeToLocationWhenParamsAreInCorrect(t *testing.T) {
	name := "Валера"
	speed := 2
	location, _ := kernel.NewLocation(1, 1)
	targetLocation := kernel.Location{}

	courier, err := NewCourier(name, speed, location)

	count, err := courier.CalculateTimeToLocation(targetLocation)

	assert.ErrorIs(t, err, ErrCourierTargetLocationNotValid)
	assert.Equal(t, float64(0), count)
}

func Test_CourierMoveWhenParamsAreCorrect(t *testing.T) {
	name := "Валера"
	speed := 2
	location, _ := kernel.NewLocation(1, 1)
	targetLocation, _ := kernel.NewLocation(5, 5)

	courier, err := NewCourier(name, speed, location)

	err = courier.Move(targetLocation)
	assert.NoError(t, err)
	assert.Equal(t, 3, courier.location.X())
	assert.Equal(t, 1, courier.location.Y())

	err = courier.Move(targetLocation)
	assert.NoError(t, err)
	assert.Equal(t, 5, courier.location.X())
	assert.Equal(t, 1, courier.location.Y())

	err = courier.Move(targetLocation)
	assert.NoError(t, err)
	assert.Equal(t, 5, courier.location.X())
	assert.Equal(t, 3, courier.location.Y())
}
