package gameplay

import (
	"image/color"

	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/ecs/transform"
	"github.com/adm87/stellar/logging"
	"github.com/adm87/stellar/rendering"
	"github.com/adm87/stellar/scene"
	"github.com/adm87/stellar/services"
	"github.com/adm87/stellar/timing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

const (
	GameplayScene scene.SceneID         = "gameplay"
	GameplayQuit  scene.SceneTransition = iota + 1
)

type Scene struct {
	world donburi.World
	root  donburi.Entity

	transformService services.ITransformService
}

func NewScene(time *timing.Time) scene.Scene {
	w := donburi.NewWorld()
	return &Scene{
		world:            w,
		root:             w.Create(transform.TransformArchetype[:]...),
		transformService: services.NewTransformService(w, time),
	}
}

func (s *Scene) EnterScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	logger.Info("Entering gameplay scene")

	buffer.Resize(800, 450)
	buffer.SetBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255})
	buffer.SetFilter(ebiten.FilterPixelated)

	s.world.OnCreate(s.onEntityCreated)
	s.world.OnRemove(s.onEntityRemoved)

	return nil
}

func (s *Scene) ExitScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	logger.Info("Exiting gameplay scene")
	return nil
}

func (s *Scene) Update(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) (scene.SceneTransition, error) {
	transform.ResolveHierarchy(s.world, s.transformService.ConsumeDirtyEntities())
	return scene.ContinueScene, nil
}

func (s *Scene) Draw(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	return nil
}

func (s *Scene) onEntityCreated(world donburi.World, entity donburi.Entity) {
	entry := world.Entry(entity)

	if !entry.HasComponent(transform.TransformHierarchyComponent) {
		return
	}

	// Attach the entity to the scene root, however it can be immediately re-parented after creation if needed.
	s.transformService.SetParent(entry, world.Entry(s.root))
}

func (s *Scene) onEntityRemoved(world donburi.World, entity donburi.Entity) {
	entry := world.Entry(entity)

	if !entry.HasComponent(transform.TransformHierarchyComponent) {
		return
	}

	// Detach the entity from its current parent to preserve the integrity of the hierarchy and prevent dangling references in the transform system.
	s.transformService.SetParent(entry, nil)
}
