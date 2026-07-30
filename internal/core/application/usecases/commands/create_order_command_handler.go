package commands

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"errors"
	"fmt"
)

type CreateOrderCommandHandler interface {
	Handle(context.Context, CreateOrderCommand) error
}

var _ CreateOrderCommandHandler = (*createOrderCommandHandler)(nil)

type createOrderCommandHandler struct {
	factory   ports.UnitOfWork
	geoClient ports.GeoClient
}

func NewCreateOrderCommandHandler(factory ports.UnitOfWork, geoClient ports.GeoClient) (CreateOrderCommandHandler, error) {
	if factory == nil {
		panic("factory is required")
	}
	if geoClient == nil {
		panic("geoClient is required")
	}

	return &createOrderCommandHandler{
		factory:   factory,
		geoClient: geoClient,
	}, nil
}

func (h *createOrderCommandHandler) Handle(ctx context.Context, command CreateOrderCommand) error {
	if !command.IsValid() {
		return errors.New("create order handler: invalid command")
	}

	location, err := h.geoClient.GetLocation(ctx, command.Street())
	if err != nil {
		return fmt.Errorf("get location: %w", err)
	}

	order, err := order.NewOrder(command.OrderID(), location, command.Volume())
	if err != nil {
		return fmt.Errorf("new order: %w", err)
	}

	err = h.factory.OrderRepository().Add(ctx, order)
	if err != nil {
		return fmt.Errorf("add: %w", err)
	}

	return nil
}
