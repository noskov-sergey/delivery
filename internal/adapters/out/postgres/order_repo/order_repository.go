package order_repo

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Repository struct {
	tracker Tracker
}

func NewRepository(tracker Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, ports.ErrOrderTrackerValueIsRequired
	}

	return &Repository{
		tracker: tracker,
	}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *order.Order) error {
	r.tracker.Track(aggregate)

	dto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	query := `INSERT INTO orders (id, courier_id, location_x, location_y, volume, status) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := tx.ExecContext(ctx, query, dto.ID, dto.CourierID, dto.LocationX, dto.LocationY, dto.Volume, dto.Status)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *order.Order) error {
	r.tracker.Track(aggregate)

	dto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	query := `UPDATE orders SET courier_id = $1, location_x = $2, location_y = $3, volume = $4, status = $5 WHERE id = $6`
	_, err := tx.ExecContext(ctx, query, dto.CourierID, dto.LocationX, dto.LocationY, dto.Volume, dto.Status, dto.ID)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*order.Order, error) {
	dto := OrderDTO{}

	db := r.tracker.Db()

	query := `SELECT id, courier_id, location_x, location_y, volume, status FROM orders WHERE id = $1`
	result := db.QueryRowContext(ctx, query, ID)
	if result == nil {
		return nil, errors.New("failed to get order")
	}

	if err := result.Scan(&dto.ID, &dto.CourierID, &dto.LocationX, &dto.LocationY, &dto.Volume, &dto.Status); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return DtoToDomain(dto), nil
}

func (r *Repository) GetRandomCreatedStatus(ctx context.Context) (*order.Order, error) {
	dto := OrderDTO{}

	db := r.tracker.Db()

	query := `SELECT id, courier_id, location_x, location_y, volume, status FROM orders WHERE status = $1 ORDER BY random() LIMIT 1`
	result := db.QueryRowContext(ctx, query, order.StatusCreated)
	if result == nil {
		return nil, errors.New("failed to get random order")
	}

	if err := result.Scan(&dto.ID, &dto.CourierID, &dto.LocationX, &dto.LocationY, &dto.Volume, &dto.Status); err != nil {
		return nil, fmt.Errorf("failed to get random order: %w", err)
	}

	return DtoToDomain(dto), nil
}

func (r *Repository) GetAllAssignedStatus(ctx context.Context) ([]*order.Order, error) {
	db := r.tracker.Db()

	query := `SELECT id, courier_id, location_x, location_y, volume, status FROM orders WHERE status = $1`
	result, err := db.QueryContext(ctx, query, order.StatusAssigned)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer result.Close()

	aggregates := make([]*order.Order, 0)
	for result.Next() {
		dto := OrderDTO{}
		if err := result.Scan(&dto.ID, &dto.CourierID, &dto.LocationX, &dto.LocationY, &dto.Volume, &dto.Status); err != nil {
			return nil, fmt.Errorf("failed to get order: %w", err)
		}

		aggregates = append(aggregates, DtoToDomain(dto))
	}

	return aggregates, nil
}
