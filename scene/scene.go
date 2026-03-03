package scene

import (
	"fmt"

	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/errs"
	"github.com/adm87/stellar/logging"
	"github.com/adm87/stellar/timing"
)

// --------------------------------------------------------------------------------
// Scene Context
// --------------------------------------------------------------------------------

// Context provides access to shared resources and services that scenes can use during their lifecycle.
type Context interface {
	Assets() *assets.Assets
	Logger() *logging.Logger
	Time() *timing.Time
}

// --------------------------------------------------------------------------------
// Scene Interface
// --------------------------------------------------------------------------------

// SceneID is a unique identifier for a scene, used by the Director to manage scene transitions.
type SceneID string

const InvalidSceneID SceneID = ""

// Scene represents a single screen or state in the game.
// It is responsible for updating and drawing itself based on the provided context.
type Scene interface {
	EnterScene(ctx Context) error
	ExitScene(ctx Context) error
	Update(ctx Context) (SceneTransition, error)
	Draw(ctx Context) error
}

// SceneFactory is a function type that creates a new instance of a Scene.
type SceneFactory func() Scene

// SceneTransition represents the condition under which a scene should exit, used to determine when to transition to another scene.
type SceneTransition uint8

// NoSceneTransition indicates that a scene should not transition to another scene and should continue running until explicitly transitioned by the Director.
const NoSceneTransition SceneTransition = 0

// --------------------------------------------------------------------------------
// Director
// --------------------------------------------------------------------------------

// Director manages the current scene and handles transitions between scenes.
// It is responsible for calling the update and draw methods of the active scene.
type Director struct {
	ctors       map[SceneID]SceneFactory
	transitions map[SceneID]map[SceneTransition]SceneID
	nextID      SceneID
	currentID   SceneID
	current     Scene
}

// NewDirector creates a new Director with no active scene.
func NewDirector() *Director {
	return &Director{
		ctors:       make(map[SceneID]SceneFactory),
		transitions: make(map[SceneID]map[SceneTransition]SceneID),
		nextID:      InvalidSceneID,
		currentID:   InvalidSceneID,
		current:     nil,
	}
}

// RegisterScene registers a new scene with the Director using a unique SceneID and a SceneFactory function.
// It returns an error if a scene with the given ID already exists.
func (d *Director) RegisterScene(id SceneID, ctor SceneFactory) error {
	if _, exists := d.ctors[id]; exists {
		return errs.DuplicateEntry{
			Message: "scene with the given ID already exists",
		}
	}

	d.ctors[id] = ctor
	return nil
}

// AddTransition adds a transition from one scene to another based on a specified condition.
// It returns an error if the scene IDs are invalid, if the transition condition is invalid, or if either scene ID is not registered.
//
// If a transition exists for the given 'from' scene and condition, it will be overwritten with the new 'to' scene ID.
func (d *Director) AddTransition(from SceneID, condition SceneTransition, to SceneID) error {
	if from == InvalidSceneID || to == InvalidSceneID {
		return errs.InvalidArg{
			Message: "scene IDs cannot be invalid",
		}
	}

	if from == to {
		return errs.InvalidArg{
			Message: "cannot add transition from a scene to itself",
		}
	}

	if _, exists := d.ctors[from]; !exists {
		return errs.InvalidArg{
			Message: fmt.Sprintf("no scene registered with ID '%s'", from),
		}
	}

	if _, exists := d.ctors[to]; !exists {
		return errs.InvalidArg{
			Message: fmt.Sprintf("no scene registered with ID '%s'", to),
		}
	}

	if _, exists := d.transitions[from]; !exists {
		d.transitions[from] = make(map[SceneTransition]SceneID)
	}

	d.transitions[from][condition] = to
	return nil
}

// TransitionTo initiates a transition to the scene with the specified SceneID.
// It returns an error if the target SceneID is invalid or if no scene is registered with the given ID.
func (d *Director) TransitionTo(id SceneID) error {
	if id == InvalidSceneID {
		return errs.InvalidArg{
			Message: "cannot transition to invalid scene ID",
		}
	}

	if id == d.currentID {
		return nil // No transition needed if already in the target scene
	}

	if _, exists := d.ctors[id]; !exists {
		return errs.InvalidOperation{
			Message: fmt.Sprintf("no scene registered with ID '%s'", id),
		}
	}

	d.nextID = id
	return nil
}

// Update updates the current scene and handles any necessary transitions based on the scene's update logic.
// It returns an error if there is an issue updating the current scene or transitioning to a new scene.
func (d *Director) Update(ctx Context) error {
	if d.current == nil && d.nextID == InvalidSceneID {
		return nil
	}

	if d.nextID != InvalidSceneID {
		if d.current != nil {
			if err := d.current.ExitScene(ctx); err != nil {
				return errs.InvalidOperation{
					Message: fmt.Sprintf("failed to exit current scene '%s': %s", d.currentID, err.Error()),
				}
			}

			d.current = nil
			d.currentID = InvalidSceneID

			return nil
		}

		nextCtor, exists := d.ctors[d.nextID]

		if !exists {
			return errs.InvalidOperation{
				Message: fmt.Sprintf("no scene registered with ID '%s'", d.nextID),
			}
		}

		d.current = nextCtor()
		d.currentID = d.nextID
		d.nextID = InvalidSceneID

		if err := d.current.EnterScene(ctx); err != nil {
			return errs.InvalidOperation{
				Message: fmt.Sprintf("failed to enter next scene '%s': %s", d.currentID, err.Error()),
			}
		}

		return nil
	}

	transition, err := d.current.Update(ctx)

	if err != nil {
		return errs.InvalidOperation{
			Message: fmt.Sprintf("failed to update current scene '%s': %s", d.currentID, err.Error()),
		}
	}

	if transition != NoSceneTransition {
		nextID, exists := d.transitions[d.currentID][transition]

		if !exists {
			return errs.InvalidOperation{
				Message: fmt.Sprintf("no transition defined from scene '%s' for condition '%d'", d.currentID, transition),
			}
		}

		d.nextID = nextID
	}

	return nil
}

// Draw calls the Draw method of the current scene to render it on the screen.
// It returns an error if there is an issue drawing the current scene.
func (d *Director) Draw(ctx Context) error {
	if d.current == nil {
		return nil
	}

	if err := d.current.Draw(ctx); err != nil {
		return errs.InvalidOperation{
			Message: fmt.Sprintf("failed to draw current scene '%s': %s", d.currentID, err.Error()),
		}
	}

	return nil
}
