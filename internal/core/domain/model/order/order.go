package order

import (
	"delivery/internal/core/domain/kernel"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrOrderIsNotAssign     = errors.New("order is not assign")
	ErrOrderIsCompleted     = errors.New("order is completed")
	ErrOrderIsAssign        = errors.New("order is already assign")
	ErrOrderIdIsEmpty       = errors.New("order id is empty")
	ErrOrderVolumeIsEmpty   = errors.New("order volume is empty")
	ErrOrderLocationIsEmpty = errors.New("order location is empty")
)

type Order struct {
	id        uuid.UUID
	courierID *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func NewOrder(id uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if id == uuid.Nil {
		return nil, ErrOrderIdIsEmpty
	}

	if volume <= 0 {
		return nil, ErrOrderVolumeIsEmpty
	}

	if !location.IsEmpty() {
		return nil, ErrOrderLocationIsEmpty
	}

	return &Order{
		id:       id,
		location: location,
		volume:   volume,
		status:   StatusCreated,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if o.courierID != nil {
		return ErrOrderIsAssign
	}

	switch o.status {
	case StatusCreated:
		o.courierID = &courierID
		return nil
	case StatusCompleted:
		return ErrOrderIsCompleted
	case StatusAssigned:
		return ErrOrderIsAssign
	}

	return nil
}

func (o *Order) Complete() error {
	switch o.status {
	case StatusCreated:
		return ErrOrderIsNotAssign
	case StatusCompleted:
		return ErrOrderIsCompleted
	case StatusAssigned:
		o.status = StatusCompleted
		return nil
	}
	return nil
}
