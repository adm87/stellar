package scene

import (
	"github.com/adm87/stellar/assets"
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

// Scene represents a single screen or state in the game.
// It is responsible for updating and drawing itself based on the provided context.
type Scene interface {
	EnterScene(ctx Context) error
	ExitScene(ctx Context) error
	Update(ctx Context) error
	Draw(ctx Context) error
}

// --------------------------------------------------------------------------------
// Director
// --------------------------------------------------------------------------------

// Director manages the current scene and handles transitions between scenes.
// It is responsible for calling the update and draw methods of the active scene.
type Director struct {
	current Scene
}

// NewDirector creates a new Director with no active scene.
func NewDirector() *Director {
	return &Director{}
}

func (d *Director) Update(ctx Context) error {
	if d.current == nil {
		return nil
	}
	return d.current.Update(ctx)
}

func (d *Director) Draw(ctx Context) error {
	if d.current == nil {
		return nil
	}
	return d.current.Draw(ctx)
}
