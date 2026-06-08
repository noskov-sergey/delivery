package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrCourierNameIsEmpty     = errors.New("courier name is empty")
	ErrCourierSpeedIsInvalid  = errors.New("courier speed is invalid")
	ErrCourierLocationIsEmpty = errors.New("courier location is empty")
	ErrCourierCantGetVolume   = errors.New("courier can not get volume")
)

type Courier struct {
	id           uuid.UUID
	name         string
	speed        int
	location     kernel.Location
	storagePlace []StoragePlace
}

func NewCourier(name string, speed int, loc kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, ErrCourierNameIsEmpty
	}

	if speed <= 0 {
		return nil, ErrCourierSpeedIsInvalid
	}

	if loc.IsEmpty() {
		return nil, ErrCourierLocationIsEmpty
	}

	return &Courier{
		id:       uuid.New(),
		name:     name,
		speed:    speed,
		location: loc,
		storagePlace: []StoragePlace{
			*NewStoragePlaceStandard(),
		},
	}, nil
}

func (c *Courier) ID() uuid.UUID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

// Logic Methods

func (c *Courier) AddStoragePlace(name string, volume int) error {
	storage, err := NewStoragePlace(name, volume)
	if err != nil {
		return fmt.Errorf("courier add storage place: %w", err)
	}

	c.storagePlace = append(c.storagePlace, *storage)

	return nil
}

func (c *Courier) CanTakeOrder(order order.Order) (bool, error) {
	var canTakeOrder int8
	for _, storage := range c.storagePlace {
		can, err := storage.CanStore(order.Volume())
		if err != nil {
			continue
		}
		if can {
			canTakeOrder++
		}
	}

	if canTakeOrder == 0 {
		return false, ErrCourierCantGetVolume
	}

	return true, nil
}

func (c *Courier) TakeOrder(order order.Order) error {
	return nil
}

func (c *Courier) CompleteOrder(order order.Order) error {
	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	return 0, nil
}

func (c *Courier) Move(target kernel.Location) error {
	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	return nil, nil
}
