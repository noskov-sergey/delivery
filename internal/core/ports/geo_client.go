package ports

import (
	"context"
	"delivery/internal/core/domain/kernel"
)

type GeoClient interface {
	GetLocation(context.Context, string) (kernel.Location, error)
	Close() error
}
