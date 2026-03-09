package transform

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/yohamta/donburi"
)

// --- Helper: Build a chain hierarchy ---
func newChainHierarchy(world donburi.World, size int) []donburi.Entity {
	entities := make([]donburi.Entity, size)

	for i := 0; i < size; i++ {
		entities[i] = world.Create(TransformArchetype[:]...)
	}

	for i := 0; i < size; i++ {
		entry := world.Entry(entities[i])
		h := TransformHierarchyComponent.Get(entry)

		if i == 0 {
			h.Parent = donburi.Null
		} else {
			h.Parent = entities[i-1]
		}

		if i < size-1 {
			h.FirstChild = entities[i+1]
		} else {
			h.FirstChild = donburi.Null
		}

		h.NextSibling = donburi.Null
		h.PrevSibling = donburi.Null
	}

	return entities
}

// --- Helper: Build a star hierarchy ---
func newStarHierarchy(leafCount int) (donburi.World, donburi.Entity, []donburi.Entity) {
	world := donburi.NewWorld()
	root := world.Create(TransformArchetype[:]...)
	leaves := make([]donburi.Entity, leafCount)

	for i := 0; i < leafCount; i++ {
		leaves[i] = world.Create(TransformArchetype[:]...)
	}

	rootH := TransformHierarchyComponent.Get(world.Entry(root))
	if leafCount > 0 {
		rootH.FirstChild = leaves[0]
	}

	for i := 0; i < leafCount; i++ {
		entry := world.Entry(leaves[i])
		h := TransformHierarchyComponent.Get(entry)

		h.Parent = root
		h.FirstChild = donburi.Null

		if i > 0 {
			h.PrevSibling = leaves[i-1]
		} else {
			h.PrevSibling = donburi.Null
		}

		if i < leafCount-1 {
			h.NextSibling = leaves[i+1]
		} else {
			h.NextSibling = donburi.Null
		}
	}

	return world, root, leaves
}

// --- Helper: Build independent roots ---
func newIndependentRoots(count int) (donburi.World, []donburi.Entity) {
	world := donburi.NewWorld()
	roots := make([]donburi.Entity, count)

	for i := 0; i < count; i++ {
		roots[i] = world.Create(TransformArchetype[:]...)
		entry := world.Entry(roots[i])
		h := TransformHierarchyComponent.Get(entry)

		h.Parent = donburi.Null
		h.FirstChild = donburi.Null
		h.NextSibling = donburi.Null
		h.PrevSibling = donburi.Null
	}

	return world, roots
}

// --- Helper: Cache entries to avoid world.Entry inside loops ---
func cacheEntries(world donburi.World, entities []donburi.Entity) []*donburi.Entry {
	entries := make([]*donburi.Entry, len(entities))
	for i, e := range entities {
		entries[i] = world.Entry(e)
	}
	return entries
}

// --- Helper: Reset state for a single dirty root ---
func resetStateSingleDirty(entries []*donburi.Entry, dirtyIndex int) {
	for i, entry := range entries {
		state := TransformStateComponent.Get(entry)
		if i == dirtyIndex {
			state.DirtyFlag = DirtyLocal
		} else {
			state.DirtyFlag = 0
		}
		state.Queued = false
	}
}

// --- Helper: Reset state for a dirty slice ---
func resetStateForSlice(entries []*donburi.Entry, dirty []donburi.Entity) {
	// Clear all
	for _, entry := range entries {
		state := TransformStateComponent.Get(entry)
		state.DirtyFlag = 0
		state.Queued = false
	}
	// Mark dirty ones
	for _, e := range dirty {
		entry := entries[e] // requires mapping if needed
		state := TransformStateComponent.Get(entry)
		state.DirtyFlag = DirtyLocal
	}
}

// --- Helper: Precompute random dirty patterns ---
func precomputeDirtyPatterns(entities []donburi.Entity, pct float64, seed int64, variants int) [][]donburi.Entity {
	rng := rand.New(rand.NewSource(seed))
	patterns := make([][]donburi.Entity, variants)

	for i := 0; i < variants; i++ {
		count := int(float64(len(entities)) * pct)
		dirty := make([]donburi.Entity, 0, count+1)

		for _, e := range entities {
			if rng.Float64() < pct {
				dirty = append(dirty, e)
			}
		}
		patterns[i] = dirty
	}

	return patterns
}

// --- Helper: Precompute shuffled dirty orders ---
func precomputeShuffledOrders(base []donburi.Entity, seed int64, variants int) [][]donburi.Entity {
	rng := rand.New(rand.NewSource(seed))
	orders := make([][]donburi.Entity, variants)

	for i := 0; i < variants; i++ {
		order := append([]donburi.Entity(nil), base...)
		rng.Shuffle(len(order), func(a, b int) {
			order[a], order[b] = order[b], order[a]
		})
		orders[i] = order
	}

	return orders
}

func BenchmarkResolveHierarchySingleDirtyRootDeepChain(b *testing.B) {
	depths := []int{100, 1000, 2500, 5000}
	world := donburi.NewWorld()

	for _, depth := range depths {
		chain := newChainHierarchy(world, depth)
		entries := cacheEntries(world, chain)
		states := make([]*TransformState, len(chain))
		dirty := []donburi.Entity{chain[0]}

		restoreFunc := func() {
			for i, state := range states {
				if i == 0 {
					state.DirtyFlag = DirtyLocal
				} else {
					state.DirtyFlag = 0
				}
				state.Queued = false
			}
		}

		for i, e := range chain {
			entry := world.Entry(e)
			entries[i] = entry
			states[i] = TransformStateComponent.Get(entry)
		}

		b.Run("depth_"+strconv.Itoa(depth), func(b *testing.B) {
			restoreFunc()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ResolveHierarchy(world, dirty)
				restoreFunc()
			}
		})
	}
}
