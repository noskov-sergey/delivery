package queries

import (
	"context"
	"delivery/internal/core/ports"
	"errors"
	"fmt"
)

type GetAllQueriersQueryHandler interface {
	Handle(context.Context, GetAllQueriersQuery) (*GetAllCouriersResponse, error)
}

var _ GetAllQueriersQueryHandler = (*getAllQueriersQueryHandler)(nil)

type getAllQueriersQueryHandler struct {
	factory ports.UnitOfWork
}

func NewGetAllQueriersQueryHandler(factory ports.UnitOfWork) (GetAllQueriersQueryHandler, error) {
	if factory == nil {
		panic("factory is required")
	}

	return &getAllQueriersQueryHandler{
		factory: factory,
	}, nil
}

func (h *getAllQueriersQueryHandler) Handle(ctx context.Context, query GetAllQueriersQuery) (*GetAllCouriersResponse, error) {
	if !query.IsValid() {
		return nil, errors.New("get all queriers: invalid command")
	}

	db := h.factory.Db()

	q := `SELECT id, name_, location_x, location_y FROM couriers`
	result, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("failed to get courier")
	}

	response := GetAllCouriersResponse{}
	for result.Next() {
		dto := CourierDTO{}

		if err := result.Scan(&dto.ID, &dto.Name, &dto.LocationX, &dto.LocationY); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		response.Couriers = append(response.Couriers, DtoToDomain(dto))
	}

	return &response, nil
}
