package hierarchy

import (
	"github.com/adm87/stellar/errs"
	"github.com/yohamta/donburi"
)

// -----------------------------------------------
// Hierarchy Component
// -----------------------------------------------

// Model represents the hierarchical relationship of an entity in the scene graph.
// It contains references to the parent entity, the first child entity, and the next and previous sibling entities.
type Model struct {
	Parent      donburi.Entity
	FirstChild  donburi.Entity
	NextSibling donburi.Entity
	PrevSibling donburi.Entity
}

// Component defines the hierarchy component type.
var Component = donburi.NewComponentType[Model]()

// -----------------------------------------------
// Hierarchy Component Helper Functions
// -----------------------------------------------

// SetParent sets the parent of the child entity to the specified parent entity, updating the hierarchy accordingly.
func SetParent(world donburi.World, child, parent donburi.Entity) error {
	childEntry := world.Entry(child)
	childModel := Component.Get(childEntry)

	if child == parent {
		return errs.InvalidOperation{
			Message: "cannot set parent: child and parent are the same entity",
		}
	}

	// Detach child from current parent if it has one.
	if childModel.Parent != donburi.Null {
		currentParentModel := Component.Get(world.Entry(childModel.Parent))

		if currentParentModel.FirstChild == child {
			currentParentModel.FirstChild = childModel.NextSibling
		}

		if childModel.PrevSibling != donburi.Null {
			Component.Get(world.Entry(childModel.PrevSibling)).NextSibling = childModel.NextSibling
		}

		if childModel.NextSibling != donburi.Null {
			Component.Get(world.Entry(childModel.NextSibling)).PrevSibling = childModel.PrevSibling
		}

		// Clear child's parent and sibling references.
		childModel.Parent = donburi.Null
		childModel.NextSibling = donburi.Null
		childModel.PrevSibling = donburi.Null
	}

	// If the new parent is null, we're done after detaching from the current parent.
	if parent == donburi.Null {
		return nil
	}

	parentEntry := world.Entry(parent)
	parentModel := Component.Get(parentEntry)

	// Prevent creating cycles in the hierarchy by ensuring the new parent is not a descendant of the child.
	if isDescendant(world, child, parentModel) {
		return errs.InvalidOperation{
			Message: "cannot set parent: would create a cycle in the hierarchy",
		}
	}

	// Attach child to new parent.
	childModel.Parent = parent
	childModel.NextSibling = parentModel.FirstChild
	childModel.PrevSibling = donburi.Null

	if parentModel.FirstChild != donburi.Null {
		Component.Get(world.Entry(parentModel.FirstChild)).PrevSibling = child
	}

	parentModel.FirstChild = child

	return nil
}

// isDescendant checks if the given descendantModel is a descendant of the ancestor entity in the hierarchy.
func isDescendant(world donburi.World, ancestor donburi.Entity, descendantModel *Model) bool {
	current := descendantModel.Parent

	// Traverse up the hierarchy from the descendant to check if we encounter the ancestor.
	// Hotpath: Deliberately avoiding recursion and additional function calls to minimize overhead.
	for current != donburi.Null {
		if current == ancestor {
			return true
		}

		current = Component.Get(world.Entry(current)).Parent
	}

	return false
}
