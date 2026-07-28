package queries

import (
	"delivery/internal/core/domain/kernel"

	"github.com/google/uuid"
)

type GetAllUncompletedOrdersResponse struct {
	Orders []Order
}

type Order struct {
	ID       uuid.UUID
	Location kernel.Location
}

type OrderDTO struct {
	ID        uuid.UUID `db:"id"`
	LocationX int       `db:"location_x"`
	LocationY int       `db:"location_y"`
}

func OrderDtoToDomain(dto OrderDTO) Order {
	location, _ := kernel.NewLocation(uint8(dto.LocationX), uint8(dto.LocationY))
	return Order{ID: dto.ID, Location: location}
}
