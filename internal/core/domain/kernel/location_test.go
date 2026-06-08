package kernel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LocationReturnWhenParamsAreCorrectOnCreated(t *testing.T) {
	tests := map[string]struct {
		x uint8
		y uint8

		expected Location
	}{
		"1,1": {
			x: 1,
			y: 1,

			expected: Location{
				x: 1,
				y: 1,

				valid: true,
			},
		},
		"6,7": {
			x: 6,
			y: 7,

			expected: Location{
				x: 6,
				y: 7,

				valid: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewLocation(test.x, test.y)

			assert.NoError(t, err)
			assert.Equal(t, got, test.expected)
			assert.Equal(t, got.IsValid(), test.expected.IsValid())
		})
	}
}

func Test_LocationReturnErrorWhenParamsAreInCorrectOnCreated(t *testing.T) {
	tests := map[string]struct {
		x uint8
		y uint8

		expected error
	}{
		"0,1": {
			x: 0,
			y: 1,

			expected: ErrXInvalidValue,
		},
		"1,0": {
			x: 1,
			y: 0,

			expected: ErrYInvalidValue,
		},
		"15,1": {
			x: 15,
			y: 1,

			expected: ErrXInvalidValue,
		},
		"1,15": {
			x: 1,
			y: 15,

			expected: ErrYInvalidValue,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewLocation(test.x, test.y)

			if err.Error() != test.expected.Error() {
				t.Errorf("expected: %v, got: %v", test.expected, err)
			}
		})
	}
}

func Test_LocationCalculateDistance(t *testing.T) {
	tests := map[string]struct {
		first  Location
		second Location

		expected int
	}{
		"1,1 and 2,2": {
			first: Location{
				x:     1,
				y:     1,
				valid: true,
			},
			second: Location{
				x:     2,
				y:     2,
				valid: true,
			},
			expected: 2,
		},
		"2,6 and 4,9": {
			first: Location{
				x:     2,
				y:     6,
				valid: true,
			},
			second: Location{
				x:     4,
				y:     9,
				valid: true,
			},
			expected: 5,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			i, err := test.first.CalculateDistance(test.second)

			assert.NoError(t, err)

			if i != test.expected {
				t.Errorf("expected: %v, got: %v", test.expected, i)
			}
		})
	}
}

func Test_LocationCalculateDistanceWithError(t *testing.T) {
	tests := map[string]struct {
		first  Location
		second Location

		expected error
	}{
		"1,1 and 2,2": {
			first: Location{
				x:     1,
				y:     1,
				valid: true,
			},
			second: Location{
				x:     2,
				y:     2,
				valid: false,
			},
			expected: ErrInvalidValue,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := test.first.CalculateDistance(test.second)

			if err.Error() != test.expected.Error() {
				t.Errorf("expected: %v, got: %v", test.expected, err)
			}
		})
	}
}
