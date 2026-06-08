package courier

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrStorageNameEmpty       = fmt.Errorf("storage name is empty")
	ErrStorageVolumeEmpty     = fmt.Errorf("storage volume is empty")
	ErrStorageIsNotEmpty      = fmt.Errorf("storage is not empty")
	ErrStorageIsSmaller       = fmt.Errorf("storage is smaller then new order")
	ErrStorageCantStore       = fmt.Errorf("can't store")
	ErrStorageOrderIsNotEqual = fmt.Errorf("order is not equal")
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, ErrStorageNameEmpty
	}

	if totalVolume <= 0 {
		return nil, ErrStorageVolumeEmpty
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}, nil
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) Equal(other StoragePlace) bool {
	if s.id == other.id {
		return true
	}

	return false
}

func (s *StoragePlace) CanStore(volume int) (bool, error) {
	if s.isOccupied() {
		return false, ErrStorageIsNotEmpty
	}

	if s.totalVolume < volume {
		return false, ErrStorageIsSmaller
	}

	return true, nil
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	c, err := s.CanStore(volume)
	if err != nil {
		return fmt.Errorf("can store: %w", err)
	}

	if !c {
		return ErrStorageCantStore
	}

	s.orderID = &orderID

	return nil
}

func (s *StoragePlace) Clear(orderID uuid.UUID) error {
	if *s.orderID != orderID {
		return ErrStorageOrderIsNotEqual
	}
	s.orderID = nil

	return nil
}

func (s *StoragePlace) isOccupied() bool {
	if s.orderID == nil {
		return false
	}
	return true
}
