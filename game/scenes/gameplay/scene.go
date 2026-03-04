package gameplay

import (
	"image/color"

	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/logging"
	"github.com/adm87/stellar/rendering"
	"github.com/adm87/stellar/scene"
	"github.com/adm87/stellar/timing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
)

const (
	GameplayScene scene.SceneID         = "gameplay"
	GameplayQuit  scene.SceneTransition = iota + 1
)

type Scene struct {
	world donburi.World
}

func NewScene() scene.Scene {
	return &Scene{
		world: donburi.NewWorld(),
	}
}

func (s *Scene) EnterScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	logger.Info("Entering gameplay scene")

	buffer.Resize(800, 450)
	buffer.SetBackgroundColor(color.RGBA{R: 100, G: 149, B: 237, A: 255})
	buffer.SetFilter(ebiten.FilterPixelated)

	return nil
}

func (s *Scene) ExitScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	logger.Info("Exiting gameplay scene")
	return nil
}

func (s *Scene) Update(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) (scene.SceneTransition, error) {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return GameplayQuit, nil
	}
	return scene.ContinueScene, nil
}

func (s *Scene) Draw(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	ebitenutil.DebugPrint(buffer.Image(), "Gameplay Scene - Press ESC to return to Splash Screen")
	return nil
}
