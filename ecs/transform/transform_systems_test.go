package transform

import (
	"testing"

	"github.com/yohamta/donburi"
)

func TestBuildUpdatePathsSkipsChildRootWhenDirtyAncestorAppearsLater(t *testing.T) {
	world := donburi.NewWorld()

	parent := world.Create(TransformArchetype[:]...)
	child := world.Create(TransformArchetype[:]...)

	parentEntry := world.Entry(parent)
	childEntry := world.Entry(child)

	parentHierarchy := TransformHierarchyComponent.Get(parentEntry)
	childHierarchy := TransformHierarchyComponent.Get(childEntry)

	parentHierarchy.FirstChild = child
	childHierarchy.Parent = parent

	TransformStateComponent.Get(parentEntry).DirtyFlag = DirtyLocal
	TransformStateComponent.Get(childEntry).DirtyFlag = DirtyLocal

	roots := buildUpdatePaths(world, []donburi.Entity{child, parent})
	if len(roots) != 1 {
		t.Fatalf("expected exactly one root, got %d", len(roots))
	}
	if roots[0].entity != parent {
		t.Fatalf("expected parent to be selected as root, got entity %v", roots[0].entity)
	}
}
