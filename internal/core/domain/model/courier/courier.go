package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
)

var (
	ErrCourierNameIsEmpty            = errors.New("courier name is empty")
	ErrCourierSpeedIsInvalid         = errors.New("courier speed is invalid")
	ErrCourierLocationIsEmpty        = errors.New("courier location is empty")
	ErrCourierCantGetVolume          = errors.New("courier can not get volume")
	ErrCourierCantTakeOrder          = errors.New("courier can not take order")
	ErrCourierOrderNotExists         = errors.New("courier order does not exist")
	ErrCourierTargetLocationNotValid = errors.New("courier target location is not valid")
)

type Courier struct {
	id           uuid.UUID
	name         string
	speed        int
	location     kernel.Location
	storagePlace []*StoragePlace
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
		storagePlace: []*StoragePlace{
			NewStoragePlaceStandard(),
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

	c.storagePlace = append(c.storagePlace, storage)

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

func (c *Courier) TakeOrder(order *order.Order) error {
	var take int8
	for _, storage := range c.storagePlace {
		can, err := storage.CanStore(order.Volume())
		if err != nil {
			continue
		}
		if !can {
			continue
		}

		err = storage.Store(order.ID(), order.Volume())
		if err != nil {
			continue
		}

		err = order.Assign(c.ID())
		if err != nil {
			_ = storage.Clear(order.ID())
			continue
		}

		take++
		break
	}

	if take == 0 {
		return ErrCourierCantTakeOrder
	}

	return nil
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	storage, err := c.findStoragePlaceByOrderID(order.ID())
	if err != nil {
		return ErrCourierOrderNotExists
	}

	err = storage.Clear(order.ID())
	if err != nil {
		return err
	}

	err = order.Complete()
	if err != nil {
		return err
	}

	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if !target.IsValid() {
		return 0, ErrCourierTargetLocationNotValid
	}

	dist, err := c.Location().CalculateDistance(target)
	if err != nil {
		return 0, err
	}

	return float64(dist / c.Speed()), nil
}

func (c *Courier) Move(target kernel.Location) error {
	if !target.IsValid() {
		return ErrCourierTargetLocationNotValid
	}

	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := kernel.NewLocation(uint8(newX), uint8(newY))
	if err != nil {
		return err
	}
	c.location = newLocation

	return nil
}

func (c *Courier) findStoragePlaceByOrderID(orderID uuid.UUID) (*StoragePlace, error) {
	for _, storage := range c.storagePlace {
		if *storage.OrderID() == orderID {
			return storage, nil
		}
	}

	return nil, ErrCourierOrderNotExists
}
