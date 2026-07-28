package queries

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"errors"
	"fmt"
)

type GetAllUncompletedOrdersQueryHandler interface {
	Handle(context.Context, GetAllUncompletedOrdersQuery) (*GetAllUncompletedOrdersResponse, error)
}

var _ GetAllUncompletedOrdersQueryHandler = (*getAllUncompletedOrdersQueryHandler)(nil)

type getAllUncompletedOrdersQueryHandler struct {
	factory ports.UnitOfWork
}

func NewGetAllUncompletedOrderQueryHandler(factory ports.UnitOfWork) (GetAllUncompletedOrdersQueryHandler, error) {
	if factory == nil {
		panic("factory is required")
	}

	return &getAllUncompletedOrdersQueryHandler{
		factory: factory,
	}, nil
}

func (h *getAllUncompletedOrdersQueryHandler) Handle(ctx context.Context, query GetAllUncompletedOrdersQuery) (*GetAllUncompletedOrdersResponse, error) {
	if !query.IsValid() {
		return nil, errors.New("get all queriers: invalid command")
	}

	db := h.factory.Db()

	q := `SELECT id, location_x, location_y FROM orders WHERE status IN ($1, $2)`
	result, err := db.QueryContext(ctx, q, order.StatusAssigned, order.StatusCreated)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("failed to get orders")
	}

	response := GetAllUncompletedOrdersResponse{}
	for result.Next() {
		dto := OrderDTO{}

		if err := result.Scan(&dto.ID, &dto.LocationX, &dto.LocationY); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		response.Orders = append(response.Orders, OrderDtoToDomain(dto))
	}

	return &response, nil
}
