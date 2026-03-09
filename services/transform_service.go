package services

import (
	"github.com/adm87/stellar/ecs/transform"
	"github.com/adm87/stellar/errs"
	"github.com/adm87/stellar/timing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// ------------------------------------------------
// Transform Services
// ------------------------------------------------

// TransformService provides an interface for managing entity transformations and their hierarchical relationships.
//
// This includes setting and getting local transformation properties, as well as retrieving world transformation
// matrices that take into account the entity's position in the hierarchy.
type ITransformService interface {
	GetPosition(entry *donburi.Entry) (float64, float64)
	SetPosition(entry *donburi.Entry, x, y float64)

	GetScale(entry *donburi.Entry) (float64, float64)
	SetScale(entry *donburi.Entry, scaleX, scaleY float64)

	GetRotation(entry *donburi.Entry) float64
	SetRotation(entry *donburi.Entry, rotation float64)

	GetParent(entry *donburi.Entry) donburi.Entity
	SetParent(childEntry, parentEntry *donburi.Entry) error

	GetLocalMatrix(entry *donburi.Entry) ebiten.GeoM
	GetWorldMatrix(entry *donburi.Entry) ebiten.GeoM

	ConsumeDirtyEntities() []donburi.Entity
}

type TransformService struct {
	world donburi.World
	time  *timing.Time

	dirty    []donburi.Entity
	dirtySet map[donburi.Entity]struct{}
}

func NewTransformService(world donburi.World, time *timing.Time) ITransformService {
	return &TransformService{
		world:    world,
		time:     time,
		dirtySet: make(map[donburi.Entity]struct{}, 128),
	}
}

func (s *TransformService) markDirty(entry *donburi.Entry, reason transform.DirtyFlag) {
	state := transform.TransformStateComponent.Get(entry)
	state.DirtyFlag |= reason

	// If the entity is already marked dirty, we don't need to track it again.
	if _, dirty := s.dirtySet[entry.Entity()]; dirty {
		return
	}

	s.dirtySet[entry.Entity()] = struct{}{}
	s.dirty = append(s.dirty, entry.Entity())
}

func (s *TransformService) ConsumeDirtyEntities() []donburi.Entity {
	out := s.dirty

	s.dirty = s.dirty[:0]
	clear(s.dirtySet)

	return out
}

func (s *TransformService) GetWorldMatrix(entry *donburi.Entry) ebiten.GeoM {
	return transform.TransformMatrixComponent.Get(entry).World
}

func (s *TransformService) GetLocalMatrix(entry *donburi.Entry) ebiten.GeoM {
	return transform.TransformMatrixComponent.Get(entry).Local
}

func (s *TransformService) GetPosition(entry *donburi.Entry) (float64, float64) {
	t := transform.TransformComponent.Get(entry)
	return t.X, t.Y
}

func (s *TransformService) SetPosition(entry *donburi.Entry, x, y float64) {
	t := transform.TransformComponent.Get(entry)

	if t.X == x && t.Y == y {
		return
	}

	t.X = x
	t.Y = y

	s.markDirty(entry, transform.DirtyLocal)
}

func (s *TransformService) GetScale(entry *donburi.Entry) (float64, float64) {
	t := transform.TransformComponent.Get(entry)
	return t.ScaleX, t.ScaleY
}

func (s *TransformService) SetScale(entry *donburi.Entry, scaleX, scaleY float64) {
	t := transform.TransformComponent.Get(entry)

	if t.ScaleX == scaleX && t.ScaleY == scaleY {
		return
	}

	t.ScaleX = scaleX
	t.ScaleY = scaleY

	s.markDirty(entry, transform.DirtyLocal)
}

func (s *TransformService) GetRotation(entry *donburi.Entry) float64 {
	t := transform.TransformComponent.Get(entry)
	return t.Rotation
}

func (s *TransformService) SetRotation(entry *donburi.Entry, rotation float64) {
	t := transform.TransformComponent.Get(entry)

	if t.Rotation == rotation {
		return
	}

	t.Rotation = rotation

	s.markDirty(entry, transform.DirtyLocal)
}

func (s *TransformService) GetParent(entry *donburi.Entry) donburi.Entity {
	return transform.TransformHierarchyComponent.Get(entry).Parent
}

func (s *TransformService) SetParent(childEntry, parentEntry *donburi.Entry) error {
	if childEntry == nil {
		return errs.InvalidArg{
			Message: "child entry cannot be nil",
		}
	}

	if err := s.internalSetParent(childEntry, parentEntry); err != nil {
		return err
	}

	s.markDirty(childEntry, transform.DirtyWorld)
	return nil
}

func (s *TransformService) internalSetParent(childEntry, parentEntry *donburi.Entry) error {
	if childEntry == nil {
		return errs.InvalidOperation{Message: "child entry cannot be nil"}
	}

	childEntity := childEntry.Entity()
	childModel := transform.TransformHierarchyComponent.Get(childEntry)

	// Determine if this operation is detaching or reparenting the child entity.
	reparenting := parentEntry != nil && parentEntry.Entity() != donburi.Null

	if reparenting {
		parentEntity := parentEntry.Entity()
		if childEntity == parentEntity {
			return errs.InvalidOperation{Message: "cannot set parent: child and parent are the same entity"}
		}

		parentModel := transform.TransformHierarchyComponent.Get(parentEntry)
		for current := parentModel.Parent; current != donburi.Null; {
			if current == childEntity {
				return errs.InvalidOperation{Message: "cannot set parent: parent is a descendant of the child"}
			}
			current = transform.TransformHierarchyComponent.Get(s.world.Entry(current)).Parent
		}
	}

	// Detach from current parent if attached
	if childModel.Parent != donburi.Null {
		currentParentEntry := s.world.Entry(childModel.Parent)
		currentParentModel := transform.TransformHierarchyComponent.Get(currentParentEntry)

		// Unlink from the sibling chain
		if currentParentModel.FirstChild == childEntity {
			currentParentModel.FirstChild = childModel.NextSibling
		}
		if childModel.PrevSibling != donburi.Null {
			transform.TransformHierarchyComponent.Get(s.world.Entry(childModel.PrevSibling)).NextSibling = childModel.NextSibling
		}
		if childModel.NextSibling != donburi.Null {
			transform.TransformHierarchyComponent.Get(s.world.Entry(childModel.NextSibling)).PrevSibling = childModel.PrevSibling
		}

		// Reset local links
		childModel.Parent = donburi.Null
		childModel.NextSibling = donburi.Null
		childModel.PrevSibling = donburi.Null
	}

	// If not reparenting, we're done after detaching from the current parent.
	if !reparenting {
		return nil
	}

	// Attach to new parent (Head Insertion)
	parentEntity := parentEntry.Entity()
	parentModel := transform.TransformHierarchyComponent.Get(parentEntry)

	childModel.Parent = parentEntity
	childModel.NextSibling = parentModel.FirstChild
	childModel.PrevSibling = donburi.Null

	if parentModel.FirstChild != donburi.Null {
		transform.TransformHierarchyComponent.Get(s.world.Entry(parentModel.FirstChild)).PrevSibling = childEntity
	}
	parentModel.FirstChild = childEntity

	return nil
}
