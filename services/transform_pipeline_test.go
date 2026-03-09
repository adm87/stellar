package services_test

import (
	"math"
	"testing"

	"github.com/adm87/stellar/ecs/transform"
	"github.com/adm87/stellar/services"
	"github.com/adm87/stellar/timing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func TestTransformPipelineComputesWorldMatrixFromParentChain(t *testing.T) {
	world := donburi.NewWorld()
	svc := services.NewTransformService(world, timing.NewTime(60))

	parent := world.Create(transform.TransformArchetype[:]...)
	child := world.Create(transform.TransformArchetype[:]...)

	parentEntry := world.Entry(parent)
	childEntry := world.Entry(child)

	svc.SetPosition(parentEntry, 10, 20)
	svc.SetPosition(childEntry, 3, 4)
	if err := svc.SetParent(childEntry, parentEntry); err != nil {
		t.Fatalf("SetParent failed: %v", err)
	}

	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())

	assertMatrixOrigin(t, svc.GetWorldMatrix(parentEntry), 10, 20)
	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 13, 24)
}

func TestTransformPipelinePropagatesParentMovementToDescendants(t *testing.T) {
	world := donburi.NewWorld()
	svc := services.NewTransformService(world, timing.NewTime(60))

	parent := world.Create(transform.TransformArchetype[:]...)
	child := world.Create(transform.TransformArchetype[:]...)

	parentEntry := world.Entry(parent)
	childEntry := world.Entry(child)

	svc.SetPosition(parentEntry, 5, 7)
	svc.SetPosition(childEntry, 2, 3)
	if err := svc.SetParent(childEntry, parentEntry); err != nil {
		t.Fatalf("SetParent failed: %v", err)
	}

	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())

	// Move only the parent, then resolve with whatever dirty batch the service emits.
	svc.SetPosition(parentEntry, 20, 30)
	dirty := svc.ConsumeDirtyEntities()
	transform.ResolveHierarchy(world, dirty)

	assertMatrixOrigin(t, svc.GetWorldMatrix(parentEntry), 20, 30)
	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 22, 33)
}

func TestTransformPipelineReparentingUsesNewParentWorldTransform(t *testing.T) {
	world := donburi.NewWorld()
	svc := services.NewTransformService(world, timing.NewTime(60))

	parentA := world.Create(transform.TransformArchetype[:]...)
	parentB := world.Create(transform.TransformArchetype[:]...)
	child := world.Create(transform.TransformArchetype[:]...)

	parentAEntry := world.Entry(parentA)
	parentBEntry := world.Entry(parentB)
	childEntry := world.Entry(child)

	svc.SetPosition(parentAEntry, 10, 0)
	svc.SetPosition(parentBEntry, 100, 0)
	svc.SetPosition(childEntry, 1, 0)

	if err := svc.SetParent(childEntry, parentAEntry); err != nil {
		t.Fatalf("initial SetParent failed: %v", err)
	}
	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())
	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 11, 0)

	if err := svc.SetParent(childEntry, parentBEntry); err != nil {
		t.Fatalf("reparent SetParent failed: %v", err)
	}

	// Reverse order to assert pipeline behavior is stable for unordered batches.
	dirty := svc.ConsumeDirtyEntities()
	reverseEntities(dirty)
	transform.ResolveHierarchy(world, dirty)

	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 101, 0)
}

func TestTransformPipelineDeepChainRootMovePropagatesToAllDescendants(t *testing.T) {
	world := donburi.NewWorld()
	svc := services.NewTransformService(world, timing.NewTime(60))

	parent := world.Create(transform.TransformArchetype[:]...)
	child := world.Create(transform.TransformArchetype[:]...)
	grandchild := world.Create(transform.TransformArchetype[:]...)
	greatGrandchild := world.Create(transform.TransformArchetype[:]...)

	parentEntry := world.Entry(parent)
	childEntry := world.Entry(child)
	grandchildEntry := world.Entry(grandchild)
	greatGrandchildEntry := world.Entry(greatGrandchild)

	svc.SetPosition(parentEntry, 10, 0)
	svc.SetPosition(childEntry, 1, 0)
	svc.SetPosition(grandchildEntry, 2, 0)
	svc.SetPosition(greatGrandchildEntry, 3, 0)

	if err := svc.SetParent(childEntry, parentEntry); err != nil {
		t.Fatalf("SetParent child->parent failed: %v", err)
	}
	if err := svc.SetParent(grandchildEntry, childEntry); err != nil {
		t.Fatalf("SetParent grandchild->child failed: %v", err)
	}
	if err := svc.SetParent(greatGrandchildEntry, grandchildEntry); err != nil {
		t.Fatalf("SetParent greatGrandchild->grandchild failed: %v", err)
	}

	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())

	// Move the root only; all descendants should receive updated world transforms.
	svc.SetPosition(parentEntry, 20, 0)
	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())

	assertMatrixOrigin(t, svc.GetWorldMatrix(parentEntry), 20, 0)
	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 21, 0)
	assertMatrixOrigin(t, svc.GetWorldMatrix(grandchildEntry), 23, 0)
	assertMatrixOrigin(t, svc.GetWorldMatrix(greatGrandchildEntry), 26, 0)
}

func TestTransformPipelineSupportsMultipleDirtyRoots(t *testing.T) {
	world := donburi.NewWorld()
	svc := services.NewTransformService(world, timing.NewTime(60))

	rootA := world.Create(transform.TransformArchetype[:]...)
	childA := world.Create(transform.TransformArchetype[:]...)
	rootB := world.Create(transform.TransformArchetype[:]...)
	childB := world.Create(transform.TransformArchetype[:]...)

	rootAEntry := world.Entry(rootA)
	childAEntry := world.Entry(childA)
	rootBEntry := world.Entry(rootB)
	childBEntry := world.Entry(childB)

	svc.SetPosition(rootAEntry, 10, 0)
	svc.SetPosition(childAEntry, 1, 0)
	svc.SetPosition(rootBEntry, 100, 0)
	svc.SetPosition(childBEntry, 5, 0)

	if err := svc.SetParent(childAEntry, rootAEntry); err != nil {
		t.Fatalf("SetParent childA->rootA failed: %v", err)
	}
	if err := svc.SetParent(childBEntry, rootBEntry); err != nil {
		t.Fatalf("SetParent childB->rootB failed: %v", err)
	}

	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())

	svc.SetPosition(rootAEntry, 20, 0)
	svc.SetPosition(rootBEntry, 200, 0)
	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())

	assertMatrixOrigin(t, svc.GetWorldMatrix(rootAEntry), 20, 0)
	assertMatrixOrigin(t, svc.GetWorldMatrix(childAEntry), 21, 0)
	assertMatrixOrigin(t, svc.GetWorldMatrix(rootBEntry), 200, 0)
	assertMatrixOrigin(t, svc.GetWorldMatrix(childBEntry), 205, 0)
}

func TestTransformPipelineSiblingReparentingRoundTrip(t *testing.T) {
	world := donburi.NewWorld()
	svc := services.NewTransformService(world, timing.NewTime(60))

	parentA := world.Create(transform.TransformArchetype[:]...)
	parentB := world.Create(transform.TransformArchetype[:]...)
	child := world.Create(transform.TransformArchetype[:]...)

	parentAEntry := world.Entry(parentA)
	parentBEntry := world.Entry(parentB)
	childEntry := world.Entry(child)

	svc.SetPosition(parentAEntry, 10, 0)
	svc.SetPosition(parentBEntry, 100, 0)
	svc.SetPosition(childEntry, 1, 0)

	if err := svc.SetParent(childEntry, parentAEntry); err != nil {
		t.Fatalf("SetParent child->A failed: %v", err)
	}
	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())
	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 11, 0)

	if err := svc.SetParent(childEntry, parentBEntry); err != nil {
		t.Fatalf("SetParent child->B failed: %v", err)
	}
	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())
	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 101, 0)

	if err := svc.SetParent(childEntry, parentAEntry); err != nil {
		t.Fatalf("SetParent child->A (again) failed: %v", err)
	}
	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())
	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 11, 0)
}

func TestTransformPipelineCombinesDirtyLocalAndDirtyWorldInSameFrame(t *testing.T) {
	world := donburi.NewWorld()
	svc := services.NewTransformService(world, timing.NewTime(60))

	parentA := world.Create(transform.TransformArchetype[:]...)
	parentB := world.Create(transform.TransformArchetype[:]...)
	child := world.Create(transform.TransformArchetype[:]...)

	parentAEntry := world.Entry(parentA)
	parentBEntry := world.Entry(parentB)
	childEntry := world.Entry(child)

	svc.SetPosition(parentAEntry, 10, 0)
	svc.SetPosition(parentBEntry, 100, 0)
	svc.SetPosition(childEntry, 1, 0)

	if err := svc.SetParent(childEntry, parentAEntry); err != nil {
		t.Fatalf("initial SetParent failed: %v", err)
	}
	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())
	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 11, 0)

	// Same frame: mutate local transform and reparent before consuming dirty entities.
	svc.SetPosition(childEntry, 7, 0)
	if err := svc.SetParent(childEntry, parentBEntry); err != nil {
		t.Fatalf("reparent SetParent failed: %v", err)
	}
	transform.ResolveHierarchy(world, svc.ConsumeDirtyEntities())

	assertMatrixOrigin(t, svc.GetWorldMatrix(childEntry), 107, 0)
}

func reverseEntities(entities []donburi.Entity) {
	for i, j := 0, len(entities)-1; i < j; i, j = i+1, j-1 {
		entities[i], entities[j] = entities[j], entities[i]
	}
}

func assertMatrixOrigin(t *testing.T, m ebiten.GeoM, wantX, wantY float64) {
	t.Helper()

	x, y := (&m).Apply(0, 0)
	if !almostEqual(x, wantX) || !almostEqual(y, wantY) {
		t.Fatalf("matrix origin = (%.6f, %.6f), want (%.6f, %.6f)", x, y, wantX, wantY)
	}
}

func almostEqual(a, b float64) bool {
	const eps = 1e-6
	return math.Abs(a-b) <= eps
}
