package commands

import (
	"context"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"errors"
)

type AssignOrdersCommandHandler interface {
	Handle(context.Context, AssignOrdersCommand) error
}

var _ AssignOrdersCommandHandler = (*assignOrdersCommandHandler)(nil)

type assignOrdersCommandHandler struct {
	factory    ports.UnitOfWork
	dispatcher services.OrderDispatchService
}

func NewAssignOrdersCommandHandler(factory ports.UnitOfWork, dispatcher services.OrderDispatchService) (AssignOrdersCommandHandler, error) {
	if factory == nil {
		panic("factory is required")
	}
	if dispatcher == nil {
		panic("dispatcher is required")
	}

	return &assignOrdersCommandHandler{
		factory:    factory,
		dispatcher: dispatcher,
	}, nil
}

func (h *assignOrdersCommandHandler) Handle(ctx context.Context, command AssignOrdersCommand) error {
	if !command.IsValid() {
		return errors.New("assign orders: invalid command")
	}

	aggregate, err := h.factory.OrderRepository().GetRandomCreatedStatus(ctx)
	if err != nil {
		return err
	}

	couriers, err := h.factory.CourierRepository().GetAllFree(ctx)
	if err != nil {
		return err
	}

	orderKeeper, err := h.dispatcher.Dispatch(aggregate, couriers)
	if err != nil {
		return err
	}

	err = h.factory.Begin(ctx)
	if err != nil {
		return err
	}

	err = h.factory.OrderRepository().Update(ctx, aggregate)
	if err != nil {
		return err
	}

	err = h.factory.CourierRepository().Update(ctx, orderKeeper)
	if err != nil {
		return err
	}

	return h.factory.Commit(ctx)
}
