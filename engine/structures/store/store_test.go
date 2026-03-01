package store

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestStoreAllocateGet(t *testing.T) {
	store := NewStore[string](2)

	id, err := store.Allocate("alpha")
	if err != nil {
		t.Fatalf("Allocate returned error: %v", err)
	}

	got, ok := store.Get(id)
	if !ok {
		t.Fatalf("Get returned ok=false for valid id")
	}
	if got != "alpha" {
		t.Fatalf("Get returned %q, want %q", got, "alpha")
	}
}

func TestStoreAllocateFull(t *testing.T) {
	store := NewStore[int](1)

	if _, err := store.Allocate(1); err != nil {
		t.Fatalf("first Allocate returned error: %v", err)
	}

	_, err := store.Allocate(2)
	if err == nil {
		t.Fatalf("second Allocate returned nil error, want StoreFull")
	}

	var storeFull StoreFull
	if !errors.As(err, &storeFull) {
		t.Fatalf("second Allocate error = %T (%v), want StoreFull", err, err)
	}
}

func TestStoreDeallocateInvalidatesID(t *testing.T) {
	store := NewStore[int](1)

	id, err := store.Allocate(42)
	if err != nil {
		t.Fatalf("Allocate returned error: %v", err)
	}

	store.Deallocate(id)

	_, ok := store.Get(id)
	if ok {
		t.Fatalf("Get returned ok=true for deallocated id")
	}
}

func TestStoreReallocateBumpsGeneration(t *testing.T) {
	store := NewStore[string](1)

	id1, err := store.Allocate("first")
	if err != nil {
		t.Fatalf("first Allocate returned error: %v", err)
	}

	store.Deallocate(id1)

	id2, err := store.Allocate("second")
	if err != nil {
		t.Fatalf("second Allocate returned error: %v", err)
	}

	if id2.index != id1.index {
		t.Fatalf("reallocated index = %d, want %d", id2.index, id1.index)
	}
	if id2.generation <= id1.generation {
		t.Fatalf("generation did not increase: old=%d new=%d", id1.generation, id2.generation)
	}

	if _, ok := store.Get(id1); ok {
		t.Fatalf("old StoreID should be stale after reallocation")
	}

	got, ok := store.Get(id2)
	if !ok || got != "second" {
		t.Fatalf("Get with new StoreID = (%q, %v), want (%q, true)", got, ok, "second")
	}
}

func TestStoreDeallocateStaleIDDoesNotRemoveNewEntry(t *testing.T) {
	store := NewStore[int](1)

	id1, err := store.Allocate(7)
	if err != nil {
		t.Fatalf("first Allocate returned error: %v", err)
	}
	store.Deallocate(id1)

	id2, err := store.Allocate(9)
	if err != nil {
		t.Fatalf("second Allocate returned error: %v", err)
	}

	store.Deallocate(id1)

	got, ok := store.Get(id2)
	if !ok || got != 9 {
		t.Fatalf("stale deallocate affected new entry: Get(id2) = (%d, %v), want (9, true)", got, ok)
	}
}

func TestStoreGetOutOfBoundsID(t *testing.T) {
	store := NewStore[int](1)

	_, ok := store.Get(StoreID{index: 5, generation: 1})
	if ok {
		t.Fatalf("Get returned ok=true for out-of-bounds id")
	}
}

func TestStoreConcurrentStress(t *testing.T) {
	const (
		workers    = 32
		iterations = 2000
	)

	store := NewStore[int](workers)

	errCh := make(chan error, workers)
	var wg sync.WaitGroup

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)

		go func(worker int) {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				expected := worker*iterations + i

				id, err := store.Allocate(expected)
				if err != nil {
					var storeFull StoreFull
					if errors.As(err, &storeFull) {
						continue
					}

					errCh <- fmt.Errorf("worker %d allocate failed at iteration %d: %w", worker, i, err)
					return
				}

				got, ok := store.Get(id)
				if !ok {
					errCh <- fmt.Errorf("worker %d get miss at iteration %d", worker, i)
					store.Deallocate(id)
					return
				}

				if got != expected {
					errCh <- fmt.Errorf("worker %d get mismatch at iteration %d: got %d want %d", worker, i, got, expected)
					store.Deallocate(id)
					return
				}

				store.Deallocate(id)

				if _, ok := store.Get(id); ok {
					errCh <- fmt.Errorf("worker %d stale id remained valid at iteration %d", worker, i)
					return
				}
			}
		}(worker)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatal(err)
	}
}

func BenchmarkStoreAllocateDeallocate(b *testing.B) {
	store := NewStore[int](1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, err := store.Allocate(i)
		if err != nil {
			b.Fatalf("Allocate returned error: %v", err)
		}
		store.Deallocate(id)
	}
}

func BenchmarkStoreGetHit(b *testing.B) {
	store := NewStore[int](1)
	id, err := store.Allocate(123)
	if err != nil {
		b.Fatalf("Allocate returned error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, ok := store.Get(id)
		if !ok {
			b.Fatalf("Get returned ok=false for valid id")
		}
	}
}

func BenchmarkStoreGetMissStaleID(b *testing.B) {
	store := NewStore[int](1)
	id, err := store.Allocate(1)
	if err != nil {
		b.Fatalf("Allocate returned error: %v", err)
	}
	store.Deallocate(id)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.Get(id)
	}
}

func BenchmarkStoreParallelGetHit(b *testing.B) {
	store := NewStore[int](1)
	id, err := store.Allocate(1)
	if err != nil {
		b.Fatalf("Allocate returned error: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = store.Get(id)
		}
	})
}
