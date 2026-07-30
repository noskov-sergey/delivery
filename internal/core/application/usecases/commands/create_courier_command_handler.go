package commands

import (
	"context"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"errors"
	"fmt"
)

type CreateCourierCommandHandler interface {
	Handle(context.Context, CreateCourierCommand) error
}

var _ CreateCourierCommandHandler = (*createCourierCommandHandler)(nil)

type createCourierCommandHandler struct {
	factory ports.UnitOfWork
}

func NewCreateCourierCommandHandler(factory ports.UnitOfWork) (CreateCourierCommandHandler, error) {
	if factory == nil {
		panic("factory is required")
	}

	return &createCourierCommandHandler{
		factory: factory,
	}, nil
}

func (h *createCourierCommandHandler) Handle(ctx context.Context, command CreateCourierCommand) error {
	if !command.IsValid() {
		return errors.New("create courier handler: invalid command")
	}

	courier, err := courier.NewCourier(command.name, command.speed, kernel.RandomLocation())
	if err != nil {
		return fmt.Errorf("new courier: %w", err)
	}

	err = h.factory.CourierRepository().Add(ctx, courier)
	if err != nil {
		return fmt.Errorf("add: %w", err)
	}

	return nil
}
