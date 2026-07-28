package commands

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type AddStoragePlaceCommand struct {
	courierID   uuid.UUID
	name        string
	totalVolume int

	valid bool
}

func NewAddStoragePlaceCommand(courierID uuid.UUID, name string, totalVolume int) (AddStoragePlaceCommand, error) {
	if courierID == uuid.Nil {
		return AddStoragePlaceCommand{}, errors.New("courierID cannot be nil")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return AddStoragePlaceCommand{}, errors.New("name cannot be empty")
	}

	if totalVolume <= 0 {
		return AddStoragePlaceCommand{}, errors.New("total volume must be greater than zero")
	}

	return AddStoragePlaceCommand{
		courierID:   courierID,
		name:        name,
		totalVolume: totalVolume,
		valid:       true,
	}, nil
}

func (a *AddStoragePlaceCommand) IsValid() bool        { return a.valid }
func (a *AddStoragePlaceCommand) CourierID() uuid.UUID { return a.courierID }
func (a *AddStoragePlaceCommand) Name() string         { return a.name }
func (a *AddStoragePlaceCommand) TotalVolume() int     { return a.totalVolume }
