package order

import (
	"delivery/internal/core/domain/kernel"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_OrderReturnWhenParamsAreCorrectOnCreated(t *testing.T) {
	location, _ := kernel.NewLocation(2, 5)
	id := uuid.New()
	volume := 5

	order, err := NewOrder(id, location, volume)

	assert.NoError(t, err)
	assert.Equal(t, id, order.ID())
	assert.Equal(t, location, order.Location())
	assert.Equal(t, volume, order.Volume())
	assert.Equal(t, StatusCreated, order.Status())
	assert.Nil(t, order.CourierID())
}

func Test_OrderReturnWhenParamsAreInCorrectOnCreated(t *testing.T) {
	goodLocation, _ := kernel.NewLocation(1, 1)
	wrongLocation := kernel.Location{}

	tests := map[string]struct {
		id       uuid.UUID
		location kernel.Location
		volume   int

		expectErr error
	}{
		"empty uuid": {
			id:       uuid.Nil,
			location: goodLocation,
			volume:   5,

			expectErr: ErrOrderIdIsEmpty,
		},
		"empty location": {
			id:       uuid.New(),
			location: wrongLocation,
			volume:   5,

			expectErr: ErrOrderLocationIsEmpty,
		},
		"empty volume": {
			id:       uuid.New(),
			location: goodLocation,
			volume:   0,

			expectErr: ErrOrderVolumeIsEmpty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			order, err := NewOrder(test.id, test.location, test.volume)

			assert.Nil(t, order)
			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}

func Test_OrderAssignWhenParamsAreCorrect(t *testing.T) {
	location, _ := kernel.NewLocation(2, 5)
	id := uuid.New()
	courierId := uuid.New()
	volume := 5

	order, err := NewOrder(id, location, volume)

	assert.NoError(t, err)
	assert.Nil(t, order.CourierID())

	err = order.Assign(courierId)
	assert.NoError(t, err)
	assert.Equal(t, courierId, *order.CourierID())
}

func Test_OrderAssignWhenParamsAreInCorrect(t *testing.T) {
	goodLocation, _ := kernel.NewLocation(1, 1)

	tests := map[string]struct {
		id        uuid.UUID
		location  kernel.Location
		volume    int
		courierID uuid.UUID
		status    Status

		expectErr error
	}{
		"status completed": {
			id:        uuid.New(),
			location:  goodLocation,
			volume:    5,
			courierID: uuid.New(),
			status:    StatusCompleted,

			expectErr: ErrOrderIsCompleted,
		},
		"status assigned": {
			id:        uuid.New(),
			location:  goodLocation,
			volume:    5,
			courierID: uuid.New(),
			status:    StatusAssigned,

			expectErr: ErrOrderIsAssign,
		},
		"empty courier id": {
			id:        uuid.New(),
			location:  goodLocation,
			volume:    5,
			courierID: uuid.Nil,
			status:    StatusCreated,

			expectErr: ErrOrderCourierIdIsEmpty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			order, err := NewOrder(test.id, test.location, test.volume)
			assert.NoError(t, err)
			assert.NotNil(t, order)

			order.status = test.status

			err = order.Assign(test.courierID)
			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}

func Test_OrderCompleteWhenParamsAreCorrect(t *testing.T) {
	location, _ := kernel.NewLocation(2, 5)
	id := uuid.New()
	volume := 5

	order, err := NewOrder(id, location, volume)

	assert.NoError(t, err)
	assert.Nil(t, order.CourierID())

	order.status = StatusAssigned

	err = order.Complete()
	assert.NoError(t, err)
}

func Test_OrderCompleteWhenParamsAreInCorrect(t *testing.T) {
	goodLocation, _ := kernel.NewLocation(1, 1)

	tests := map[string]struct {
		id       uuid.UUID
		location kernel.Location
		volume   int
		status   Status

		expectErr error
	}{
		"status completed": {
			id:       uuid.New(),
			location: goodLocation,
			volume:   5,
			status:   StatusCompleted,

			expectErr: ErrOrderIsCompleted,
		},
		"status created": {
			id:       uuid.New(),
			location: goodLocation,
			volume:   5,
			status:   StatusCreated,

			expectErr: ErrOrderIsNotAssign,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			order, err := NewOrder(test.id, test.location, test.volume)
			assert.NoError(t, err)
			assert.NotNil(t, order)

			order.status = test.status

			err = order.Complete()
			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}
