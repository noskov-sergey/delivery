package order_repo

import (
	"delivery/internal/core/domain/model/order"

	"github.com/google/uuid"
)

type OrderDTO struct {
	ID        uuid.UUID    `db:"id"`
	CourierID *uuid.UUID   `db:"courier_id"`
	LocationX int          `db:"location_x"`
	LocationY int          `db:"location_y"`
	Volume    int          `db:"volume"`
	Status    order.Status `db:"status"`
}

func (OrderDTO) TableName() string {
	return "orders"
}
