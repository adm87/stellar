package hierarchy

import (
	"github.com/adm87/stellar/errs"
	"github.com/yohamta/donburi"
)

// -----------------------------------------------
// Hierarchy Component
// -----------------------------------------------

// Hierarchy represents the hierarchical relationship of an entity in the scene graph.
// It contains references to the parent entity, the first child entity, and the next and previous sibling entities.
type Hierarchy struct {
	parent      donburi.Entity
	firstChild  donburi.Entity
	nextSibling donburi.Entity
	prevSibling donburi.Entity
}

// HierarchyComponent defines the hierarchy component type.
var HierarchyComponent = donburi.NewComponentType[Hierarchy]()

// -----------------------------------------------
// Hierarchy Component Helper Functions
// -----------------------------------------------

// GetParent retrieves the parent entity of the specified entry's hierarchy component.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Hierarchy component.
func GetParent(world donburi.World, entry *donburi.Entry) donburi.Entity {
	return HierarchyComponent.Get(entry).parent
}

// SetParent sets the parent of the child entity to the specified parent entity, updating the hierarchy accordingly.
//
// Assumptions:
//   - parentEntry can be nil, which indicates detaching the child from its current parent.
//   - Both childEntry and parentEntry (if not nil) are valid entries in the world.
//   - Both entities have a Hierarchy component.
//
// Errors:
//   - Returns an error if attempting to set an entity as its own parent.
//   - Returns an error if attempting to set a parent that is a descendant of the child, which would create a cycle.
func SetParent(world donburi.World, childEntry, parentEntry *donburi.Entry) error {
	if childEntry == nil {
		return errs.InvalidOperation{
			Message: "child entry cannot be nil",
		}
	}

	childModel := HierarchyComponent.Get(childEntry)
	childEntity := childEntry.Entity()

	// If the new parent is nil/null, detach and return after unlinking from current parent.
	isDetach := parentEntry == nil || parentEntry.Entity() == donburi.Null

	if !isDetach && childEntity == parentEntry.Entity() {
		return errs.InvalidOperation{
			Message: "cannot set parent: child and parent are the same entity",
		}
	}

	// Detach child from current parent if it has one.
	if childModel.parent != donburi.Null {
		currentParentEntry := world.Entry(childModel.parent)
		currentParentModel := HierarchyComponent.Get(currentParentEntry)

		if currentParentModel.firstChild == childEntity {
			currentParentModel.firstChild = childModel.nextSibling
		}

		if childModel.prevSibling != donburi.Null {
			prevSiblingModel := HierarchyComponent.Get(world.Entry(childModel.prevSibling))
			prevSiblingModel.nextSibling = childModel.nextSibling
		}

		if childModel.nextSibling != donburi.Null {
			nextSiblingModel := HierarchyComponent.Get(world.Entry(childModel.nextSibling))
			nextSiblingModel.prevSibling = childModel.prevSibling
		}

		// Clear child's parent and sibling references.
		childModel.parent = donburi.Null
		childModel.nextSibling = donburi.Null
		childModel.prevSibling = donburi.Null
	}

	// If the new parent is null, we're done after detaching from the current parent.
	if isDetach {
		return nil
	}

	parentEntity := parentEntry.Entity()
	parentModel := HierarchyComponent.Get(parentEntry)

	// Prevent creating cycles in the hierarchy by ensuring the new parent is not a descendant of the child.
	for current := parentModel.parent; current != donburi.Null; {
		if current == childEntity {
			return errs.InvalidOperation{
				Message: "cannot set parent: parent is a descendant of the child",
			}
		}

		currentEntry := world.Entry(current)
		current = HierarchyComponent.Get(currentEntry).parent
	}

	// Attach child to new parent.
	childModel.parent = parentEntity
	childModel.nextSibling = parentModel.firstChild
	childModel.prevSibling = donburi.Null

	if parentModel.firstChild != donburi.Null {
		HierarchyComponent.Get(world.Entry(parentModel.firstChild)).prevSibling = childEntity
	}

	parentModel.firstChild = childEntity

	return nil
}

// EachChild iterates over the child entities of the specified entry's hierarchy component, invoking the provided function for each child.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Hierarchy component.
//   - The provided function does not modify the hierarchy structure (e.g., by changing parent-child relationships) during iteration,
//     as this could lead to undefined behavior.
func EachChild(world donburi.World, entry *donburi.Entry, fn func(*donburi.Entry) bool) {
	model := HierarchyComponent.Get(entry)

	for current := model.firstChild; current != donburi.Null; {
		childEntry := world.Entry(current)
		childModel := HierarchyComponent.Get(childEntry)

		next := childModel.nextSibling
		if !fn(childEntry) {
			break
		}
		current = next
	}
}
