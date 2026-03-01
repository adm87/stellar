package store

import "sync"

type StoreFull struct{}

func (e StoreFull) Error() string {
	return "store is full"
}

// StoreID is a unique identifier for a resource in the store.
type StoreID struct {
	index      uint32
	generation uint32
}

type storeSlot[T any] struct {
	data       T
	generation uint32
	active     bool
}

// Store is a generic resource store that manages resources of type T.
type Store[T any] struct {
	slots []storeSlot[T]
	free  []uint32
	mu    sync.RWMutex
}

// NewStore creates a new Store with the specified capacity.
func NewStore[T any](capacity uint32) *Store[T] {
	slots := make([]storeSlot[T], capacity)
	free := make([]uint32, capacity)

	for i := range capacity {
		free[i] = i
	}

	return &Store[T]{
		slots: slots,
		free:  free,
	}
}

// Allocate adds a new resource to the store and returns its StoreID. If the store is full, it returns an error.
func (c *Store[T]) Allocate(data T) (StoreID, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.free) == 0 {
		return StoreID{}, StoreFull{}
	}

	index := c.free[len(c.free)-1]
	c.free = c.free[:len(c.free)-1]

	c.slots[index] = storeSlot[T]{
		data:       data,
		generation: c.slots[index].generation + 1,
		active:     true,
	}

	return StoreID{
		index:      index,
		generation: c.slots[index].generation,
	}, nil
}

// Get retrieves a resource from the store by its StoreID. It returns the resource and a boolean indicating whether the resource was found and is valid.
func (c *Store[T]) Get(id StoreID) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if id.index >= uint32(len(c.slots)) {
		var zero T
		return zero, false
	}

	slot := c.slots[id.index]

	if !slot.active || slot.generation != id.generation {
		var zero T
		return zero, false
	}

	return slot.data, true
}

// Deallocate removes a resource from the store by its StoreID. It marks the slot as inactive and adds it back to the free list.
//
// Note: Deallocate does not release the resource from memory; it simply marks the slot as available for reuse.
// The caller is responsible for ensuring that any necessary cleanup of the resource is performed before calling Deallocate.
func (c *Store[T]) Deallocate(id StoreID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if id.index >= uint32(len(c.slots)) {
		return
	}

	slot := &c.slots[id.index]

	if !slot.active || slot.generation != id.generation {
		return
	}

	var zero T
	slot.data = zero

	slot.active = false
	c.free = append(c.free, id.index)
}
