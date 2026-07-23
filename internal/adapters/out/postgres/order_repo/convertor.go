package order_repo

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
)

func DomainToDTO(aggregate *order.Order) OrderDTO {
	return OrderDTO{
		ID:        aggregate.ID(),
		CourierID: aggregate.CourierID(),
		LocationX: aggregate.Location().X(),
		LocationY: aggregate.Location().Y(),
		Volume:    aggregate.Volume(),
		Status:    aggregate.Status(),
	}
}

func DtoToDomain(dto OrderDTO) *order.Order {
	location, _ := kernel.NewLocation(uint8(dto.LocationX), uint8(dto.LocationY))
	return order.RestoreOrder(dto.ID, dto.CourierID, location, dto.Volume, dto.Status)
}
