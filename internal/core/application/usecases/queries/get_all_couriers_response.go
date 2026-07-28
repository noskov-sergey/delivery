package queries

import (
	"delivery/internal/core/domain/kernel"

	"github.com/google/uuid"
)

type GetAllCouriersResponse struct {
	Couriers []Courier
}

type Courier struct {
	ID       uuid.UUID
	Name     string
	Location kernel.Location
}

type CourierDTO struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	LocationX int       `db:"location_x"`
	LocationY int       `db:"location_y"`
}

func DtoToDomain(dto CourierDTO) Courier {
	location, _ := kernel.NewLocation(uint8(dto.LocationX), uint8(dto.LocationY))
	return Courier{ID: dto.ID, Name: dto.Name, Location: location}
}
