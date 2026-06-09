package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"errors"
)

var ErrOrderIsNil = errors.New("order is nil")
var ErrCantFoundCourier = errors.New("can't find courier")

type OrderDispatchService interface {
	Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

var _ OrderDispatchService = &orderDispatchService{}

type orderDispatchService struct{}

func NewOrderDispatchService() OrderDispatchService {
	return &orderDispatchService{}
}

func (o *orderDispatchService) Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if order == nil {
		return nil, ErrOrderIsNil
	}

	var minimum = 99.00
	var winner *courier.Courier
	for _, courier := range couriers {
		could, err := courier.CanTakeOrder(*order)
		if err != nil {
			continue
		}

		if !could {
			continue
		}

		count, err := courier.CalculateTimeToLocation(order.Location())
		if err != nil {
			continue
		}

		if count < minimum {
			minimum = count
			winner = courier
		}
	}

	if winner == nil {
		return nil, ErrCantFoundCourier
	}

	err := winner.TakeOrder(order)
	if err != nil {
		return nil, err
	}

	return winner, nil
}
