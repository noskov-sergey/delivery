package commands

import (
	"context"
	"delivery/internal/core/ports"
	"errors"
)

type MoveCouriersCommandHandler interface {
	Handle(context.Context, MoveCouriersCommand) error
}

var _ MoveCouriersCommandHandler = (*moveCouriersCommandHandler)(nil)

type moveCouriersCommandHandler struct {
	factory ports.UnitOfWork
}

func NewMoveCouriersCommandHandler(factory ports.UnitOfWork) (MoveCouriersCommandHandler, error) {
	if factory == nil {
		panic("factory is required")
	}

	return &moveCouriersCommandHandler{
		factory: factory,
	}, nil
}

func (h *moveCouriersCommandHandler) Handle(ctx context.Context, command MoveCouriersCommand) error {
	if !command.IsValid() {
		return errors.New("move couriers handler: invalid command")
	}

	orders, err := h.factory.OrderRepository().GetAllAssignedStatus(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		courier, err := h.factory.CourierRepository().Get(ctx, *order.CourierID())
		if err != nil {
			return err
		}

		err = courier.Move(order.Location())
		if err != nil {
			return err
		}

		if courier.Location().Equal(order.Location()) {
			err = courier.CompleteOrder(order)
			if err != nil {
				return err
			}
		}

		err = h.factory.OrderRepository().Update(ctx, order)
		if err != nil {
			return err
		}

		err = h.factory.CourierRepository().Update(ctx, courier)
		if err != nil {
			return err
		}
	}

	return h.factory.Commit(ctx)
}
