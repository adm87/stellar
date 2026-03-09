package transform

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type DirtyFlag uint8

const (
	DirtyLocal DirtyFlag = 1 << iota // Indicates that the local transform of the entity has changed and needs to be updated.
	DirtyWorld                       // Indicates that the world transform of the entity has changed and needs to be updated.
)

// ------------------------------------------------
// Transform Archetype
// ------------------------------------------------

// TransformArchetype defines the set of components that make up a complete transform for an entity.
var TransformArchetype = [4]donburi.IComponentType{
	TransformComponent,
	TransformMatrixComponent,
	TransformHierarchyComponent,
	TransformStateComponent,
}

// -----------------------------------------------
// Transform Component Models
// -----------------------------------------------

// Transform represents a local transformation of an entity.
type Transform struct {
	X        float64
	Y        float64
	ScaleX   float64
	ScaleY   float64
	Rotation float64
}

// TransformMatrix represents the local and world transformation matrices of an entity.
type TransformMatrix struct {
	Local ebiten.GeoM
	World ebiten.GeoM
}

// TransformHierarchy represents the hierarchical relationship of an entity in the scene graph for transformation purposes.
// It contains references to the parent entity, the first child entity, and the next and previous sibling entities.
type TransformHierarchy struct {
	Parent      donburi.Entity
	FirstChild  donburi.Entity
	NextSibling donburi.Entity
	PrevSibling donburi.Entity
}

// TransformState is a component that tracks the state of an entity's transformation, including
// a dirty flag to indicate when the transform has changed and needs to be updated.
type TransformState struct {
	DirtyFlag DirtyFlag
	Queued    bool
}

// -----------------------------------------------
// Transform Component Types
// -----------------------------------------------

// TransformComponent is the component type for storing local transformation properties of an entity.
var TransformComponent = donburi.NewComponentType[Transform](
	Transform{
		ScaleX: 1,
		ScaleY: 1,
	},
)

// TransformMatrixComponent is the component type for storing the transformation matrices of an entity.
var TransformMatrixComponent = donburi.NewComponentType[TransformMatrix]()

// TransformHierarchyComponent is the component type for storing the hierarchical relationship of an entity in the scene graph.
var TransformHierarchyComponent = donburi.NewComponentType[TransformHierarchy]()

// TransformStateComponent is the component type for storing the state of an entity's transformation, including dirty flags.
var TransformStateComponent = donburi.NewComponentType[TransformState]()

// ------------------------------------------------
// Transform Component Accessors
// ------------------------------------------------

// TryGetTransform attempts to retrieve the Transform component from the given entity entry.
// It returns a pointer to the Transform component if it exists, or nil if the entity does not have a Transform component.
//
// Note: Entities should contain the full TransformArchetype to ensure that all related components are present and properly initialized.
func TryGetTransform(entry *donburi.Entry) (*Transform, bool) {
	if !entry.HasComponent(TransformComponent) {
		return nil, false
	}
	return TransformComponent.Get(entry), true
}

// TryGetTransformMatrix attempts to retrieve the TransformMatrix component from the given entity entry.
// It returns a pointer to the TransformMatrix component if it exists, or nil if the entity does not have a TransformMatrix component.
//
// Note: Entities should contain the full TransformArchetype to ensure that all related components are present and properly initialized.
func TryGetTransformMatrix(entry *donburi.Entry) (*TransformMatrix, bool) {
	if !entry.HasComponent(TransformMatrixComponent) {
		return nil, false
	}
	return TransformMatrixComponent.Get(entry), true
}

// TryGetTransformHierarchy attempts to retrieve the TransformHierarchy component from the given entity entry.
// It returns a pointer to the TransformHierarchy component if it exists, or nil if the entity does not have a TransformHierarchy component.
//
// Note: Entities should contain the full TransformArchetype to ensure that all related components are present and properly initialized.
func TryGetTransformHierarchy(entry *donburi.Entry) (*TransformHierarchy, bool) {
	if !entry.HasComponent(TransformHierarchyComponent) {
		return nil, false
	}
	return TransformHierarchyComponent.Get(entry), true
}

// TryGetTransformState attempts to retrieve the TransformState component from the given entity entry.
// It returns a pointer to the TransformState component if it exists, or nil if the entity does not have a TransformState component.
//
// Note: Entities should contain the full TransformArchetype to ensure that all related components are present and properly initialized.
func TryGetTransformState(entry *donburi.Entry) (*TransformState, bool) {
	if !entry.HasComponent(TransformStateComponent) {
		return nil, false
	}
	return TransformStateComponent.Get(entry), true
}
