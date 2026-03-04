package transform

import (
	"math"
	"testing"

	"github.com/adm87/stellar/ecs/hierarchy"
	"github.com/yohamta/donburi"
)

var benchmarkTransformX, benchmarkTransformY float64

func TestUpdate_RootTransformSetsWorldFromLocal(t *testing.T) {
	world := donburi.NewWorld()

	e := world.Create(TransformComponent, ComputedTransformComponent)
	entry := world.Entry(e)

	SetPosition(world, entry, 7, -2)

	Update(world)

	if IsDirty(world, entry) {
		t.Fatal("expected root transform to be clean after update")
	}

	worldMatrix := GetWorldTransform(world, entry)
	x, y := worldMatrix.Apply(0, 0)

	if math.Abs(x-7) > 1e-9 || math.Abs(y+2) > 1e-9 {
		t.Fatalf("root world origin = (%v, %v), want (7, -2)", x, y)
	}
}

func TestUpdate_HierarchyUpdatesChildWhenParentChanges(t *testing.T) {
	world := donburi.NewWorld()

	root := world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent)
	child := world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent)

	rootEntry := world.Entry(root)
	childEntry := world.Entry(child)

	if err := hierarchy.SetParent(world, childEntry, rootEntry); err != nil {
		t.Fatalf("set child parent failed: %v", err)
	}

	SetPosition(world, rootEntry, 2, 3)
	SetPosition(world, childEntry, 4, 5)

	Update(world)

	firstWorld := GetWorldTransform(world, childEntry)
	x, y := firstWorld.Apply(0, 0)
	if math.Abs(x-6) > 1e-9 || math.Abs(y-8) > 1e-9 {
		t.Fatalf("first child world origin = (%v, %v), want (6, 8)", x, y)
	}

	SetPosition(world, rootEntry, 10, 20)
	Update(world)

	secondWorld := GetWorldTransform(world, childEntry)
	x, y = secondWorld.Apply(0, 0)
	if math.Abs(x-14) > 1e-9 || math.Abs(y-25) > 1e-9 {
		t.Fatalf("second child world origin = (%v, %v), want (14, 25)", x, y)
	}
}

func TestUpdate_TraversesThroughNonTransformHierarchyNodes(t *testing.T) {
	world := donburi.NewWorld()

	root := world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent)
	group := world.Create(hierarchy.HierarchyComponent)
	leaf := world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent)

	rootEntry := world.Entry(root)
	groupEntry := world.Entry(group)
	leafEntry := world.Entry(leaf)

	if err := hierarchy.SetParent(world, groupEntry, rootEntry); err != nil {
		t.Fatalf("set group parent failed: %v", err)
	}
	if err := hierarchy.SetParent(world, leafEntry, groupEntry); err != nil {
		t.Fatalf("set leaf parent failed: %v", err)
	}

	SetPosition(world, rootEntry, 10, 5)
	SetPosition(world, leafEntry, 3, 4)

	Update(world)

	worldMatrix := GetWorldTransform(world, leafEntry)
	x, y := worldMatrix.Apply(0, 0)

	if math.Abs(x-13) > 1e-9 || math.Abs(y-9) > 1e-9 {
		t.Fatalf("leaf world origin = (%v, %v), want (13, 9)", x, y)
	}
}

func TestQueries(t *testing.T) {
	world := donburi.NewWorld()

	// 1. Solo: Should match TransformQuery
	world.Create(TransformComponent, ComputedTransformComponent)

	// 2. Hierarchical: Should match TransformHierarchyQuery
	world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent)

	// 3. Logic Folder: Should match neither (missing Transform components)
	world.Create(hierarchy.HierarchyComponent)

	// 4. Incomplete: Should match neither (missing Computed component)
	world.Create(TransformComponent)

	// Verify TransformQuery
	soloCount := 0
	TransformQuery.Each(world, func(e *donburi.Entry) {
		soloCount++
	})
	if soloCount != 1 {
		t.Errorf("TransformQuery: expected 1 match, got %d", soloCount)
	}

	// Verify TransformHierarchyQuery
	hierarchyCount := 0
	TransformHierarchyQuery.Each(world, func(e *donburi.Entry) {
		hierarchyCount++
	})
	if hierarchyCount != 1 {
		t.Errorf("TransformHierarchyQuery: expected 1 match, got %d", hierarchyCount)
	}
}

func TestQuery_Structure(t *testing.T) {
	world := donburi.NewWorld()

	// Parent (Hierarchical)
	parent := world.Entry(world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent))

	// Folder (Hierarchy only - No Transform)
	folder := world.Entry(world.Create(hierarchy.HierarchyComponent))
	hierarchy.SetParent(world, folder, parent)

	// Child (Hierarchical)
	child := world.Entry(world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent))
	hierarchy.SetParent(world, child, folder)

	// TransformHierarchyQuery should find 2 entities (parent and child)
	// It should NOT find the folder.
	count := 0
	TransformHierarchyQuery.Each(world, func(e *donburi.Entry) {
		count++
	})

	if count != 2 {
		t.Errorf("expected 2 hierarchical transforms, got %d (likely missed child or included folder)", count)
	}
}

func BenchmarkUpdate_RootEntity(b *testing.B) {
	world := donburi.NewWorld()

	e := world.Create(TransformComponent, ComputedTransformComponent)
	entry := world.Entry(e)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SetPosition(world, entry, float64(i), float64(-i))
		Update(world)

		worldMatrix := GetWorldTransform(world, entry)
		benchmarkTransformX, benchmarkTransformY = worldMatrix.Apply(0, 0)
	}
}

func BenchmarkUpdate_TwoLevelHierarchy(b *testing.B) {
	world := donburi.NewWorld()

	root := world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent)
	child := world.Create(TransformComponent, ComputedTransformComponent, hierarchy.HierarchyComponent)

	rootEntry := world.Entry(root)
	childEntry := world.Entry(child)

	if err := hierarchy.SetParent(world, childEntry, rootEntry); err != nil {
		b.Fatalf("set child parent failed: %v", err)
	}

	SetPosition(world, childEntry, 10, 10)
	Update(world)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SetPosition(world, rootEntry, float64(i), float64(i))
		Update(world)

		worldMatrix := GetWorldTransform(world, childEntry)
		benchmarkTransformX, benchmarkTransformY = worldMatrix.Apply(0, 0)
	}
}
