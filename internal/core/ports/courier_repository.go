package ports

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCourierTrackerValueIsRequired = errors.New("tracker is required")
)

type CourierRepository interface {
	Add(ctx context.Context, aggregate *courier.Courier) error
	Update(ctx context.Context, aggregate *courier.Courier) error
	Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error)
	GetAllFree(ctx context.Context) ([]*courier.Courier, error)
}
