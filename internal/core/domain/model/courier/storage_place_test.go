package courier

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_StoragePlaceReturnWhenParamsAreCorrectOnCreated(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		expected  StoragePlace
		expectErr error
	}{
		"box": {
			name:        "box",
			totalVolume: 10,

			expected: StoragePlace{
				name:        "box",
				totalVolume: 10,
			},
			expectErr: nil,
		},
		"baggage": {
			name:        "baggage",
			totalVolume: 100,

			expected: StoragePlace{
				name:        "baggage",
				totalVolume: 100,
			},
			expectErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewStoragePlace(test.name, test.totalVolume)

			assert.Equal(t, got.Name(), test.expected.Name())
			assert.Equal(t, got.TotalVolume(), test.expected.TotalVolume())
			assert.NoError(t, err)
		})
	}
}

func Test_StoragePlaceReturnWhenParamsAreInCorrectOnCreated(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		expectErr error
	}{
		"box error": {
			name:        "box",
			totalVolume: -1,

			expectErr: ErrStorageVolumeEmpty,
		},
		"baggage name error": {
			name:        "",
			totalVolume: 5,

			expectErr: ErrStorageNameEmpty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewStoragePlace(test.name, test.totalVolume)

			assert.ErrorIs(t, err, test.expectErr)
		})
	}
}

func Test_StoragePlaceEqual(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		expected bool
	}{
		"equal": {
			name:        "box",
			totalVolume: 10,

			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage, err := NewStoragePlace(test.name, test.totalVolume)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, storage.Equal(*storage))
		})
	}
}

func Test_StoragePlaceNotEqual(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		expected bool
	}{
		"not equal": {
			name:        "box",
			totalVolume: 10,

			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage, err := NewStoragePlace(test.name, test.totalVolume)
			assert.NoError(t, err)

			other, err := NewStoragePlace(test.name, test.totalVolume)

			assert.NoError(t, err)
			assert.Equal(t, test.expected, storage.Equal(*other))
		})
	}
}

func Test_StorageCanStoreWithCorrectParams(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		storeVolume int

		expected  bool
		expectErr error
	}{
		"good case": {
			name:        "box",
			totalVolume: 10,

			storeVolume: 5,

			expected:  true,
			expectErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage, err := NewStoragePlace(test.name, test.totalVolume)
			assert.NoError(t, err)

			c, err := storage.CanStore(5)
			assert.Equal(t, test.expectErr, err)
			assert.Equal(t, test.expected, c)
		})
	}
}

func Test_StorageCanStoreWithInCorrectParams(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		storeVolume int
		occupied    bool

		expected  bool
		expectErr error
	}{
		"volume is bigger case": {
			name:        "box",
			totalVolume: 10,

			storeVolume: 20,
			occupied:    false,

			expected:  false,
			expectErr: ErrStorageIsSmaller,
		},
		"occupied": {
			name:        "box",
			totalVolume: 10,

			storeVolume: 5,
			occupied:    true,

			expected:  false,
			expectErr: ErrStorageIsNotEmpty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage, err := NewStoragePlace(test.name, test.totalVolume)
			assert.NoError(t, err)

			if test.occupied {
				err = storage.Store(uuid.New(), 1)
				assert.NoError(t, err)
			}

			c, err := storage.CanStore(test.storeVolume)
			assert.ErrorIs(t, err, test.expectErr)
			assert.Equal(t, test.expected, c)
		})
	}
}

func Test_StorageStoreWithCorrectParams(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		storeVolume int
		storeUUID   uuid.UUID

		expectErr error
	}{
		"good case": {
			name:        "box",
			totalVolume: 10,

			storeVolume: 5,
			storeUUID:   uuid.New(),

			expectErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage, err := NewStoragePlace(test.name, test.totalVolume)
			assert.NoError(t, err)

			err = storage.Store(test.storeUUID, test.storeVolume)

			assert.Equal(t, test.expectErr, err)
			assert.Equal(t, storage.orderID, &test.storeUUID)
		})
	}
}

func Test_StorageStoreWithInCorrectParams(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		storeVolume int
		occupied    bool

		expected  bool
		expectErr error
	}{
		"volume is bigger case": {
			name:        "box",
			totalVolume: 10,

			storeVolume: 20,
			occupied:    false,

			expected:  false,
			expectErr: ErrStorageIsSmaller,
		},
		"occupied": {
			name:        "box",
			totalVolume: 10,

			storeVolume: 5,
			occupied:    true,

			expected:  false,
			expectErr: ErrStorageIsNotEmpty,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage, err := NewStoragePlace(test.name, test.totalVolume)
			assert.NoError(t, err)

			if test.occupied {
				err = storage.Store(uuid.New(), 1)
				assert.NoError(t, err)
			}

			c, err := storage.CanStore(test.storeVolume)
			assert.ErrorIs(t, err, test.expectErr)
			assert.Equal(t, test.expected, c)
		})
	}
}

func Test_StorageClearWithCorrectParams(t *testing.T) {
	tests := map[string]struct {
		name        string
		totalVolume int

		storeUUID uuid.UUID

		expectErr error
	}{
		"good case": {
			name:        "box",
			totalVolume: 10,

			storeUUID: uuid.New(),

			expectErr: nil,
		},
		"bad case": {
			name:        "box",
			totalVolume: 10,

			storeUUID: uuid.New(),

			expectErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			storage, err := NewStoragePlace(test.name, test.totalVolume)
			assert.NoError(t, err)

			err = storage.Store(test.storeUUID, 1)
			assert.NoError(t, err)

			if test.expectErr != nil {
				test.storeUUID = uuid.Nil
			}

			err = storage.Clear(test.storeUUID)

			assert.Equal(t, test.expectErr, err)
		})
	}
}
