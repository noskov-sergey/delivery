package kernel

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

const (
	_min = 1
	_max = 10
)

var (
	ErrXInvalidValue = errors.New("x must be between 0 and 10")
	ErrYInvalidValue = errors.New("y must be between 0 and 10")
	ErrInvalidValue  = errors.New("other location not valid")
)

type Location struct {
	x uint8
	y uint8

	valid bool
}

func NewLocation(x uint8, y uint8) (Location, error) {
	if x < _min || x > _max {
		return Location{}, ErrXInvalidValue
	}

	if y < _min || y > _max {
		return Location{}, ErrYInvalidValue
	}

	return Location{
		x: x,
		y: y,

		valid: true,
	}, nil
}

func RandomLocation() Location {
	return Location{
		x: uint8(rand.Intn(_max-_min) + _min),
		y: uint8(rand.Intn(_max-_min) + _min),

		valid: true,
	}
}

func (l Location) Equal(other Location) bool {
	return l.x == other.x && l.y == other.y
}

func (l Location) IsValid() bool {
	return l.valid
}

func (l Location) CalculateDistance(other Location) (int, error) {
	if !other.IsValid() {
		return 0, ErrInvalidValue
	}

	return int(math.Abs(float64(int8(l.x)-int8(other.x))) + math.Abs(float64(int8(l.y)-int8(other.y)))), nil
}

func (l Location) String() string {
	return fmt.Sprintf("{x:%d,y:%d}", l.x, l.y)
}
