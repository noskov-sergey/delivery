package courier_repo

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"

	"github.com/google/uuid"
)

func DomainToDTO(aggregate *courier.Courier) (CourierDTO, []StoragePlaceDTO) {
	storages := make([]StoragePlaceDTO, len(aggregate.StoragePlaces()))
	storagesIDs := make([]string, len(aggregate.StoragePlaces()))

	for i, storage := range aggregate.StoragePlaces() {
		orderID := uuid.Nil
		if storage.OrderID() != nil {
			orderID = *storage.OrderID()
		}

		storages[i] = StoragePlaceDTO{
			ID:          storage.ID(),
			Name:        storage.Name(),
			TotalVolume: storage.TotalVolume(),
			OrderID:     orderID,
			CourierID:   aggregate.ID(),
		}

		storagesIDs[i] = storage.ID().String()
	}

	return CourierDTO{
		ID:            aggregate.ID(),
		Name:          aggregate.Name(),
		Speed:         aggregate.Speed(),
		LocationX:     aggregate.Location().X(),
		LocationY:     aggregate.Location().Y(),
		StoragePlaces: storagesIDs,
	}, storages
}

func DtoToDomain(dto CourierDTO, storageDto []StoragePlaceDTO) *courier.Courier {
	location, _ := kernel.NewLocation(uint8(dto.LocationX), uint8(dto.LocationY))
	storagePlaces := make([]*courier.StoragePlace, len(storageDto))
	for i, storagePlace := range storageDto {
		storageNew := courier.RestoreStoragePlace(
			storagePlace.ID, storagePlace.Name, storagePlace.TotalVolume, &storagePlace.OrderID)

		storagePlaces[i] = storageNew
	}

	return courier.RestoreCourier(dto.ID, dto.Name, dto.Speed, location, storagePlaces)
}
