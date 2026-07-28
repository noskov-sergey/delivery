package commands

import (
	"context"
	"delivery/internal/core/ports"
	"errors"
)

type AddStoragePlaceCommandHandler interface {
	Handle(context.Context, AddStoragePlaceCommand) error
}

var _ AddStoragePlaceCommandHandler = (*addStoragePlaceCommandHandler)(nil)

type addStoragePlaceCommandHandler struct {
	factory ports.UnitOfWork
}

func NewAddStoragePlaceCommandHandler(factory ports.UnitOfWork) (AddStoragePlaceCommandHandler, error) {
	if factory == nil {
		panic("factory is required")
	}

	return &addStoragePlaceCommandHandler{
		factory: factory,
	}, nil
}

func (h *addStoragePlaceCommandHandler) Handle(ctx context.Context, command AddStoragePlaceCommand) error {
	if !command.IsValid() {
		return errors.New("create add storage place handler: invalid command")
	}

	courier, err := h.factory.CourierRepository().Get(ctx, command.CourierID())
	if err != nil {
		return err
	}

	err = courier.AddStoragePlace(command.name, command.totalVolume)
	if err != nil {
		return err
	}

	err = h.factory.CourierRepository().Update(ctx, courier)
	if err != nil {
		return err
	}

	return h.factory.Commit(ctx)
}
