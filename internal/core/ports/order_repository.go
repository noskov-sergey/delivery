package ports

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrOrderTrackerValueIsRequired = errors.New("tracker is required")
	ErrOrderNotFound               = errors.New("order not found")
)

type OrderRepository interface {
	Add(ctx context.Context, aggregate *order.Order) error
	Update(ctx context.Context, aggregate *order.Order) error
	Get(ctx context.Context, ID uuid.UUID) (*order.Order, error)
	GetRandomCreatedStatus(ctx context.Context) (*order.Order, error)
	GetAllAssignedStatus(ctx context.Context) ([]*order.Order, error)
}
