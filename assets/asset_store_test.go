package assets

import (
	"errors"
	"testing"

	"github.com/adm87/stellar/errs"
)

func TestAssetStoreAddGetRemoveAndReuse(t *testing.T) {
	store := NewAssetStore[int](2)

	id1, err := store.Add(42)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	value, err := store.Get(id1)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if value != 42 {
		t.Fatalf("Get returned %d, want 42", value)
	}

	if err := store.Remove(id1); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if _, err := store.Get(id1); err == nil {
		t.Fatal("expected Get on removed ID to fail")
	} else {
		var opErr errs.InvalidOperation
		if !errors.As(err, &opErr) {
			t.Fatalf("expected InvalidOperation, got %T", err)
		}
	}

	if err := store.Remove(id1); err == nil {
		t.Fatal("expected Remove on removed ID to fail")
	} else {
		var opErr errs.InvalidOperation
		if !errors.As(err, &opErr) {
			t.Fatalf("expected InvalidOperation, got %T", err)
		}
	}

	id2, err := store.Add(99)
	if err != nil {
		t.Fatalf("second Add failed: %v", err)
	}

	if id2.Idx() != id1.Idx() {
		t.Fatalf("expected slot reuse at index %d, got %d", id1.Idx(), id2.Idx())
	}
	if id2.Gen() <= id1.Gen() {
		t.Fatalf("expected generation to increase, old=%d new=%d", id1.Gen(), id2.Gen())
	}

	if _, err := store.Get(id1); err == nil {
		t.Fatal("expected stale ID to fail after slot reuse")
	}

	value, err = store.Get(id2)
	if err != nil {
		t.Fatalf("Get with new ID failed: %v", err)
	}
	if value != 99 {
		t.Fatalf("Get returned %d, want 99", value)
	}
}

func TestAssetStoreCapacityLimit(t *testing.T) {
	store := NewAssetStore[int](1)

	if _, err := store.Add(1); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}

	if _, err := store.Add(2); err == nil {
		t.Fatal("expected Add to fail at capacity")
	} else {
		var capErr errs.MaxCapacity
		if !errors.As(err, &capErr) {
			t.Fatalf("expected MaxCapacity, got %T", err)
		}
	}
}

func TestAssetStoreBoundsChecks(t *testing.T) {
	store := NewAssetStore[int](1)
	badID := AssetID{index: 100, generation: 1}

	if _, err := store.Get(badID); err == nil {
		t.Fatal("expected out-of-bounds Get to fail")
	} else {
		var boundsErr errs.IndexOutOfBounds
		if !errors.As(err, &boundsErr) {
			t.Fatalf("expected IndexOutOfBounds, got %T", err)
		}
	}

	if err := store.Remove(badID); err == nil {
		t.Fatal("expected out-of-bounds Remove to fail")
	} else {
		var boundsErr errs.IndexOutOfBounds
		if !errors.As(err, &boundsErr) {
			t.Fatalf("expected IndexOutOfBounds, got %T", err)
		}
	}
}

func TestAssetStoreGenerationOverflowGuard(t *testing.T) {
	store := NewAssetStore[int](1)
	store.entries[0].generation = ^uint32(0)

	if _, err := store.Add(7); err == nil {
		t.Fatal("expected Add to fail on generation overflow")
	} else {
		var capErr errs.MaxCapacity
		if !errors.As(err, &capErr) {
			t.Fatalf("expected MaxCapacity, got %T", err)
		}
	}
}

func BenchmarkAssetStoreAddRemove(b *testing.B) {
	store := NewAssetStore[int](1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, err := store.Add(i)
		if err != nil {
			b.Fatalf("Add failed: %v", err)
		}
		if err := store.Remove(id); err != nil {
			b.Fatalf("Remove failed: %v", err)
		}
	}
}

func BenchmarkAssetStoreGet(b *testing.B) {
	store := NewAssetStore[int](1)
	id, err := store.Add(123)
	if err != nil {
		b.Fatalf("Add failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := store.Get(id); err != nil {
			b.Fatalf("Get failed: %v", err)
		}
	}
}

func BenchmarkAssetStoreGetParallel(b *testing.B) {
	store := NewAssetStore[int](1)
	id, err := store.Add(123)
	if err != nil {
		b.Fatalf("Add failed: %v", err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := store.Get(id); err != nil {
				b.Fatalf("Get failed: %v", err)
			}
		}
	})
}

func BenchmarkAssetStoreRemove(b *testing.B) {
	store := NewAssetStore[int](1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, err := store.Add(i)
		if err != nil {
			b.Fatalf("Add failed: %v", err)
		}
		if err := store.Remove(id); err != nil {
			b.Fatalf("Remove failed: %v", err)
		}
	}
}
