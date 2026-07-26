package courier_repo

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Repository struct {
	tracker Tracker
}

func NewRepository(tracker Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, ports.ErrCourierTrackerValueIsRequired
	}

	return &Repository{
		tracker: tracker,
	}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	r.tracker.Track(aggregate)

	courierDto, storageDto := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	query := "INSERT INTO couriers (id, name_, speed, location_x, location_y, storage_places) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := tx.ExecContext(ctx, query, courierDto.ID, courierDto.Name, courierDto.Speed, courierDto.LocationX,
		courierDto.LocationY, courierDto.StoragePlaces)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if len(storageDto) != 0 {
		var valueStrings []string
		var valueArgs []interface{}

		phCount := 1
		for _, s := range storageDto {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
				phCount, phCount+1, phCount+2, phCount+3, phCount+4))
			valueArgs = append(valueArgs, s.ID, s.Name, s.TotalVolume, s.OrderID, s.CourierID)
			phCount += 5
		}

		query = fmt.Sprintf("INSERT INTO storage_places (id, name_, total_volume, order_id, courier_id) VALUES %s",
			strings.Join(valueStrings, ", "))

		_, err := tx.ExecContext(ctx, query, valueArgs...)
		if err != nil {
			return fmt.Errorf("exec storage: %w", err)
		}
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	dto := CourierDTO{}

	db := r.tracker.Db()

	query := "SELECT id, name_, speed, location_x, location_y, storage_places FROM couriers WHERE id = $1"
	result := db.QueryRowContext(ctx, query, ID)
	if result == nil {
		return nil, errors.New("failed to get courier")
	}

	if err := result.Scan(&dto.ID, &dto.Name, &dto.Speed, &dto.LocationX, &dto.LocationY, &dto.StoragePlaces); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	storagePlacesDto := make([]StoragePlaceDTO, 0, len(dto.StoragePlaces))

	if len(dto.StoragePlaces) != 0 {
		placeholders := make([]string, len(dto.StoragePlaces))
		for i, _ := range placeholders {
			placeholders[i] = "$" + strconv.Itoa(i+1)
		}

		query = fmt.Sprintf("SELECT id, name_, total_volume, order_id, courier_id  FROM storage_places WHERE id IN (%s)", strings.Join(placeholders, ","))

		args := make([]any, len(dto.StoragePlaces))
		for i, id := range dto.StoragePlaces {
			args[i] = id
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			return nil, fmt.Errorf("query storage: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var storagePlace StoragePlaceDTO
			if err := rows.Scan(&storagePlace.ID, &storagePlace.Name, &storagePlace.TotalVolume,
				&storagePlace.OrderID, &storagePlace.CourierID); err != nil {
				return nil, err
			}
			storagePlacesDto = append(storagePlacesDto, storagePlace)
		}
	}

	return DtoToDomain(dto, storagePlacesDto), nil
}

func (r *Repository) GetAllFree(ctx context.Context) ([]*courier.Courier, error) {
	db := r.tracker.Db()

	query := `SELECT id, name_, speed, location_x, location_y, storage_places FROM couriers WHERE id NOT IN (SELECT courier_id FROM storage_places sp WHERE sp.courier_id = couriers.id AND sp.order_id != '00000000-0000-0000-0000-000000000000')`
	result, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all couriers: %w", err)
	}
	defer result.Close()
	if result == nil {
		return nil, errors.New("failed to get courier")
	}

	couriers := []*courier.Courier{}

	for result.Next() {
		dto := CourierDTO{}
		if err := result.Scan(&dto.ID, &dto.Name, &dto.Speed, &dto.LocationX, &dto.LocationY, &dto.StoragePlaces); err != nil {
			return nil, fmt.Errorf("failed to get courier: %w", err)
		}

		storagePlacesDto := make([]StoragePlaceDTO, 0, len(dto.StoragePlaces))
		if len(dto.StoragePlaces) != 0 {
			placeholders := make([]string, len(dto.StoragePlaces))
			for i := range placeholders {
				placeholders[i] = "$" + strconv.Itoa(i+1)
			}

			query = fmt.Sprintf("SELECT id, name_, total_volume, order_id, courier_id FROM storage_places WHERE id IN (%s)", strings.Join(placeholders, ","))

			args := make([]any, len(dto.StoragePlaces))
			for i, id := range dto.StoragePlaces {
				args[i] = id
			}

			rows, err := db.Query(query, args...)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var storagePlace StoragePlaceDTO
				if err := rows.Scan(&storagePlace.ID, &storagePlace.Name, &storagePlace.TotalVolume, &storagePlace.OrderID, &storagePlace.CourierID); err != nil {
					return nil, err
				}
				storagePlacesDto = append(storagePlacesDto, storagePlace)
			}
		}

		couriers = append(couriers, DtoToDomain(dto, storagePlacesDto))
	}

	return couriers, nil
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	r.tracker.Track(aggregate)

	courier, storages := DomainToDTO(aggregate)

	isInTransaction := r.tracker.InTx()
	if !isInTransaction {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	query := "UPDATE couriers SET name_ = $1, speed = $2, location_x = $3, location_y = $4, storage_places = $5 WHERE id = $6"
	_, err := tx.ExecContext(ctx, query, courier.Name, courier.Speed, courier.LocationX, courier.LocationY, courier.StoragePlaces, aggregate.ID())
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	for _, storage := range storages {
		query = "UPDATE storage_places SET name_ = $1, total_volume = $2, order_id = $3, courier_id = $4 WHERE id = $5"
		_, err := tx.ExecContext(ctx, query, storage.Name, storage.TotalVolume, storage.OrderID, storage.CourierID, storage.ID)
		if err != nil {
			return fmt.Errorf("exec storage place: %w", err)
		}
	}

	if !isInTransaction {
		err := r.tracker.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
