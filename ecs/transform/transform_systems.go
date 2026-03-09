package transform

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type stackEntry struct {
	entity       donburi.Entity
	state        *TransformState
	transform    *Transform
	matrix       *TransformMatrix
	parentMatrix *TransformMatrix
	hierarchy    *TransformHierarchy
}

// TransformQuery defines a query for entities that have the complete set of transform components,
// allowing systems to efficiently retrieve and operate on these entities.
var TransformQuery = donburi.NewQuery(
	filter.Contains(TransformArchetype[:]...),
)

var (
	roots = make([]stackEntry, 0, 100)
	stack = make([]stackEntry, 0, 100)
)

// ------------------------------------------------
// Transform Update System
// ------------------------------------------------

// ResolveHierarchy resolves the hierarchy of entities that need their transforms updated based on their dirty flags.
//
// Assumption 1: All provided entities fulfill the TransformArchetype and are dirty.
// The caller is responsible for ensuring this. This function does not validate archetype membership,
// and only applies a defensive skip for entities whose dirty flag is not set.
//
// Assumption 2: All provided entities are valid within the world, and are not dangling references.
// The caller is responsible for ensuring this, as this function does not perform any checks or validation.
//
// Note: This is still work in progress and not thread safe.
func ResolveHierarchy(world donburi.World, entities []donburi.Entity) {
	roots = buildUpdatePaths(world, entities)

	for len(roots) > 0 {
		n := len(roots) - 1

		current := roots[n]
		roots = roots[:n]

		traverseHierarchy(world, current)
	}
}

func buildUpdatePaths(world donburi.World, entities []donburi.Entity) []stackEntry {
	roots = roots[:0]

	for len(entities) > 0 {
		n := len(entities) - 1

		entry := world.Entry(entities[n])
		entities = entities[:n]

		state := TransformStateComponent.Get(entry)
		if state.DirtyFlag == 0 {
			continue // Skip entities that are not dirty
		}

		hierarchy := TransformHierarchyComponent.Get(entry)
		ancestor := hierarchy.Parent

		// If our parent is null, we are a root entity.
		if ancestor == donburi.Null {
			roots = append(roots, stackEntry{
				entity:       entry.Entity(),
				state:        state,
				transform:    TransformComponent.Get(entry),
				matrix:       TransformMatrixComponent.Get(entry),
				parentMatrix: nil,
				hierarchy:    hierarchy,
			})

			state.Queued = true
			continue // Root entities will be processed in the next loop
		}

		ancestorEntry := world.Entry(ancestor)
		ancestorState := TransformStateComponent.Get(ancestorEntry)

		// Traverse up the hierarchy to find the nearest ancestor that is already queued for update.
		// If we find one, we can skip this entity for now, as it will be processed when its ancestor is processed.
		for ancestor != donburi.Null && !ancestorState.Queued {
			if ancestorState.DirtyFlag != 0 {
				break // This ancestor is dirty and should already be or will be queued for update, so we can stop traversing up.
			}

			ancestorHierarchy := TransformHierarchyComponent.Get(ancestorEntry)
			ancestor = ancestorHierarchy.Parent

			if ancestor != donburi.Null {
				ancestorEntry = world.Entry(ancestor)
				ancestorState = TransformStateComponent.Get(ancestorEntry)
			}
		}

		// If we reached the top of the hierarchy without finding a queued parent, we are considered a root for this update cycle.
		if ancestor == donburi.Null {
			parentEntry := world.Entry(hierarchy.Parent)
			roots = append(roots, stackEntry{
				entity:       entry.Entity(),
				state:        state,
				transform:    TransformComponent.Get(entry),
				matrix:       TransformMatrixComponent.Get(entry),
				parentMatrix: TransformMatrixComponent.Get(parentEntry),
				hierarchy:    hierarchy,
			})
		}

		// We're done, this entity is now queued.
		state.Queued = true
	}

	return roots
}

func traverseHierarchy(world donburi.World, current stackEntry) {
	stack = stack[:0]
	stack = append(stack, current)

	for len(stack) > 0 {
		n := len(stack) - 1

		current := stack[n]
		stack = stack[:n]

		if current.state.DirtyFlag&DirtyLocal != 0 {
			updateLocalTransform(current)
		}

		updateWorldTransform(current)

		current.state.DirtyFlag = 0
		current.state.Queued = false

		child := current.hierarchy.FirstChild
		for child != donburi.Null {
			childEntry := world.Entry(child)
			childHierarchy := TransformHierarchyComponent.Get(childEntry)

			stack = append(stack, stackEntry{
				entity:       child,
				state:        TransformStateComponent.Get(childEntry),
				transform:    TransformComponent.Get(childEntry),
				matrix:       TransformMatrixComponent.Get(childEntry),
				parentMatrix: current.matrix,
				hierarchy:    childHierarchy,
			})

			child = childHierarchy.NextSibling
		}
	}
}

func updateLocalTransform(current stackEntry) {
	current.matrix.Local.Reset()
	current.matrix.Local.Rotate(current.transform.Rotation)
	current.matrix.Local.Scale(current.transform.ScaleX, current.transform.ScaleY)
	current.matrix.Local.Translate(current.transform.X, current.transform.Y)
}

func updateWorldTransform(current stackEntry) {
	if current.hierarchy.Parent == donburi.Null {
		current.matrix.World = current.matrix.Local
		return
	}
	current.matrix.World = current.parentMatrix.World
	current.matrix.World.Concat(current.matrix.Local)
}
