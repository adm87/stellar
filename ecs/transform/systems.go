package transform

import (
	"github.com/adm87/stellar/ecs/hierarchy"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	// TransformQuery retrieves all entities that have both Transform and ComputedTransform components but do not have a Hierarchy component.
	// This query is used for processing entities that are not part of a hierarchy (i.e., root-level entities).
	TransformQuery = donburi.NewQuery(
		filter.And(
			filter.Contains(
				TransformComponent,
				ComputedTransformComponent,
			),
			filter.Not(
				filter.Contains(
					hierarchy.HierarchyComponent,
				),
			),
		),
	)

	// TransformHierarchyQuery retrieves all entities that have Transform, ComputedTransform, and Hierarchy components.
	// This query is used for processing entities that are part of a hierarchy (i.e., child entities).
	TransformHierarchyQuery = donburi.NewQuery(
		filter.Contains(
			TransformComponent,
			ComputedTransformComponent,
			hierarchy.HierarchyComponent,
		),
	)
)

// Update processes all entities with Transform components, calculating their local and world transformation matrices.
//
// Assumptions:
//   - Entities with Transform components must also have ComputedTransform components to store the calculated matrices.
//   - Entities that are part of a hierarchy (i.e., have a Hierarchy component) will have their world transformations
//     calculated based on their local transformations and their parent's world transformations.
//   - Root-level entities (i.e., those without a parent) will have their world transformations set directly from their local transformations.
func Update(world donburi.World) {
	TransformQuery.Each(world, func(e *donburi.Entry) {
		t := TransformComponent.Get(e)

		if t.isDirty {
			ct := ComputedTransformComponent.Get(e)
			updateLocal(t, ct)
			ct.world = ct.local
		}
	})

	TransformHierarchyQuery.Each(world, func(e *donburi.Entry) {
		if hierarchy.GetParent(world, e) == donburi.Null {
			updateHierarchy(world, e, ebiten.GeoM{}, false)
		}
	})
}

// updateLocal calculates the local transformation matrix for a given Transform component and updates the corresponding ComputedTransform component.
func updateLocal(t *Transform, ct *ComputedTransform) {
	ct.local.Reset()
	ct.local.Scale(t.scaleX, t.scaleY)
	ct.local.Rotate(t.rotation)
	ct.local.Translate(t.x, t.y)
	t.isDirty = false
}

// updateHierarchy recursively updates the world transformation matrix for an entity and its descendants in the hierarchy, taking into account any changes in the parent entities.
func updateHierarchy(world donburi.World, e *donburi.Entry, parentMatrix ebiten.GeoM, parentDirty bool) {
	t := TransformComponent.Get(e)
	ct := ComputedTransformComponent.Get(e)

	shouldUpdate := t.isDirty || parentDirty

	if shouldUpdate {
		if t.isDirty {
			updateLocal(t, ct)
		}

		ct.world.Reset()
		ct.world.Concat(ct.local)
		ct.world.Concat(parentMatrix)
	}

	updateHierarchyDescendants(world, e, ct.world, shouldUpdate)
}

// updateHierarchyDescendants recursively updates the world transformation matrices for all descendant entities in the hierarchy,
// ensuring that any changes in the parent entities are propagated down the hierarchy.
func updateHierarchyDescendants(world donburi.World, e *donburi.Entry, matrix ebiten.GeoM, parentDirty bool) {
	hierarchy.EachChild(world, e, func(ce *donburi.Entry) bool {
		if ce.HasComponent(TransformComponent) && ce.HasComponent(ComputedTransformComponent) {
			updateHierarchy(world, ce, matrix, parentDirty)
		} else {
			updateHierarchyDescendants(world, ce, matrix, parentDirty)
		}
		return true
	})
}
