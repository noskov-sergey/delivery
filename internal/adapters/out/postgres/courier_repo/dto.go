package courier_repo

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CourierDTO struct {
	ID            uuid.UUID      `db:"id"`
	Name          string         `db:"name"`
	Speed         int            `db:"speed"`
	LocationX     int            `db:"location_x"`
	LocationY     int            `db:"location_y"`
	StoragePlaces pq.StringArray `db:"storage_places"`
}

func (CourierDTO) TableName() string {
	return "couriers"
}

type LocationDTO struct {
	X, Y int
}

type StoragePlaceDTO struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	TotalVolume int       `db:"total_volume"`
	OrderID     uuid.UUID `db:"order_id"`
	CourierID   uuid.UUID `db:"courier_id"`
}

func (StoragePlaceDTO) TableName() string {
	return "storage_places"
}
