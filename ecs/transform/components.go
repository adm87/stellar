package transform

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// -----------------------------------------------
// Transform Component
// -----------------------------------------------

type Transform struct {
	x, y           float64
	scaleX, scaleY float64
	rotation       float64
	isDirty        bool
}

var TransformComponent = donburi.NewComponentType[Transform](
	Transform{
		x:        0,
		y:        0,
		scaleX:   1,
		scaleY:   1,
		rotation: 0,
		isDirty:  true,
	},
)

type ComputedTransform struct {
	local ebiten.GeoM
	world ebiten.GeoM
}

var ComputedTransformComponent = donburi.NewComponentType[ComputedTransform]()

// -----------------------------------------------
// Transform Component Helper Functions
// -----------------------------------------------

// GetPosition retrieves the local position (x, y) of the specified entry's transform component.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func GetPosition(world donburi.World, entry *donburi.Entry) (float64, float64) {
	t := TransformComponent.Get(entry)
	return t.x, t.y
}

// SetPosition sets the local position (x, y) of the specified entry's transform component and marks it as dirty.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func SetPosition(world donburi.World, entry *donburi.Entry, x, y float64) {
	t := TransformComponent.Get(entry)

	if t.x == x && t.y == y {
		return
	}

	t.x = x
	t.y = y
	t.isDirty = true
}

// GetScale retrieves the local scale (scaleX, scaleY) of the specified entry's transform component.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func GetScale(world donburi.World, entry *donburi.Entry) (float64, float64) {
	t := TransformComponent.Get(entry)
	return t.scaleX, t.scaleY
}

// SetScale sets the local scale (scaleX, scaleY) of the specified entry's transform component and marks it as dirty.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func SetScale(world donburi.World, entry *donburi.Entry, scaleX, scaleY float64) {
	t := TransformComponent.Get(entry)

	if t.scaleX == scaleX && t.scaleY == scaleY {
		return
	}

	t.scaleX = scaleX
	t.scaleY = scaleY
	t.isDirty = true
}

// GetRotation retrieves the local rotation (in radians) of the specified entry's transform component.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func GetRotation(world donburi.World, entry *donburi.Entry) float64 {
	return TransformComponent.Get(entry).rotation
}

// SetRotation sets the local rotation (in radians) of the specified entry's transform component and marks it as dirty.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func SetRotation(world donburi.World, entry *donburi.Entry, rotation float64) {
	t := TransformComponent.Get(entry)

	if t.rotation == rotation {
		return
	}

	t.rotation = rotation
	t.isDirty = true
}

// IsDirty checks if the specified entry's transform component is marked as dirty, indicating that it has been modified since the last update.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func IsDirty(world donburi.World, entry *donburi.Entry) bool {
	t := TransformComponent.Get(entry)
	return t.isDirty
}

// ClearDirty clears the dirty flag of the specified entry's transform component, indicating that it has been updated and is no longer dirty.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func ClearDirty(world donburi.World, entry *donburi.Entry) {
	TransformComponent.Get(entry).isDirty = false
}

// GetLocalTransform retrieves the local transformation matrix of the specified entry's computed transform component.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func GetLocalTransform(world donburi.World, entry *donburi.Entry) ebiten.GeoM {
	return ComputedTransformComponent.Get(entry).local
}

// GetWorldTransform retrieves the world transformation matrix of the specified entry's computed transform component.
//
// Assumptions:
//   - entry is a valid entry in the world.
//   - The entity associated with the entry has a Transform component.
func GetWorldTransform(world donburi.World, entry *donburi.Entry) ebiten.GeoM {
	return ComputedTransformComponent.Get(entry).world
}
