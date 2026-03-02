package assets

import (
	"sync"

	"github.com/adm87/stellar/errs"
)

// AssetID is a unique identifier for an asset in the AssetStore,
// consisting of an index and a generation count.
type AssetID struct {
	index      uint32
	generation uint32
}

func (id AssetID) Idx() uint32 {
	return id.index
}

func (id AssetID) Gen() uint32 {
	return id.generation
}

// AssetSlot represents a slot in the AssetStore, containing the asset data,
// its generation count, and whether it is active.
type AssetSlot[T any] struct {
	data       T
	generation uint32
	active     bool
}

// AssetStore is a generic store for managing assets of type T, using an
// index-based system with generation counts to ensure safe access and removal.
type AssetStore[T any] struct {
	entries []AssetSlot[T]
	free    []uint32
	cap     uint32
	mu      sync.RWMutex
}

// NewAssetStore creates a new AssetStore with the specified capacity.
// The capacity must be greater than 0, panics otherwise.
//
// AssetSlot indices support a max generation count of 2^32 - 1.
// Add returns an error if a slot would overflow this limit.
func NewAssetStore[T any](capacity uint32) *AssetStore[T] {
	if capacity == 0 {
		panic("asset store capacity must be greater than 0")
	}

	entries := make([]AssetSlot[T], capacity)
	free := make([]uint32, capacity)

	for i := range capacity {
		free[i] = capacity - 1 - i
	}

	return &AssetStore[T]{
		entries: entries,
		free:    free,
		cap:     capacity,
	}
}

// Add adds a new asset to the store and returns its AssetID.
// If the store is at maximum capacity, returns an error.
func (s *AssetStore[T]) Add(data T) (AssetID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.free) == 0 {
		var zero AssetID
		return zero, errs.MaxCapacity{
			Message:  "asset store is at maximum capacity",
			Capacity: int(s.cap),
		}
	}

	index := s.free[len(s.free)-1]
	s.free = s.free[:len(s.free)-1]

	gen := s.entries[index].generation

	if gen == ^uint32(0) {
		var zero AssetID
		return zero, errs.MaxCapacity{
			Message:  "asset slot generation count overflow",
			Capacity: int(s.cap),
		}
	}

	nextGen := gen + 1

	s.entries[index] = AssetSlot[T]{
		data:       data,
		generation: nextGen,
		active:     true,
	}

	return AssetID{
		index:      index,
		generation: nextGen,
	}, nil
}

// Remove removes the asset with the given AssetID from the store.
// If the AssetID is invalid or already removed, returns an error.
func (s *AssetStore[T]) Remove(id AssetID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id.index >= uint32(len(s.entries)) {
		return errs.IndexOutOfBounds{
			Message: "asset ID index is out of bounds",
		}
	}

	slot := &s.entries[id.index]

	if !slot.active || slot.generation != id.generation {
		return errs.InvalidOperation{
			Message: "asset ID is invalid or already removed",
		}
	}

	var zero T
	slot.data = zero

	slot.active = false
	s.free = append(s.free, id.index)

	return nil
}

// Get retrieves the asset associated with the given AssetID.
// If the AssetID is invalid or removed, returns an error.
func (s *AssetStore[T]) Get(id AssetID) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if id.index >= uint32(len(s.entries)) {
		var zero T
		return zero, errs.IndexOutOfBounds{
			Message: "asset ID index is out of bounds",
		}
	}

	slot := s.entries[id.index]

	if !slot.active || slot.generation != id.generation {
		var zero T
		return zero, errs.InvalidOperation{
			Message: "asset ID is invalid or removed",
		}
	}

	return slot.data, nil
}
