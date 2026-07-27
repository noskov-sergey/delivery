package commands

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type CreateOrderCommand struct {
	orderID uuid.UUID
	street  string
	volume  int

	valid bool
}

func NewCreateOrderCommand(orderID uuid.UUID, street string, volume int) (CreateOrderCommand, error) {
	if orderID == uuid.Nil {
		return CreateOrderCommand{}, errors.New("orderID cannot be nil")
	}

	street = strings.TrimSpace(street)
	if street == "" {
		return CreateOrderCommand{}, errors.New("street cannot be empty")
	}

	if volume <= 0 {
		return CreateOrderCommand{}, errors.New("volume must be greater than zero")
	}

	return CreateOrderCommand{
		orderID: orderID,
		street:  street,
		volume:  volume,
		valid:   true,
	}, nil
}

func (c *CreateOrderCommand) IsValid() bool      { return c.valid }
func (c *CreateOrderCommand) OrderID() uuid.UUID { return c.orderID }
func (c *CreateOrderCommand) Street() string     { return c.street }
func (c *CreateOrderCommand) Volume() int        { return c.volume }
