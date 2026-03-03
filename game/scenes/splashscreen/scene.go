package splashscreen

import (
	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/logging"
	"github.com/adm87/stellar/rendering"
	"github.com/adm87/stellar/scene"
	"github.com/adm87/stellar/timing"
)

const (
	SplashScreenScene    scene.SceneID         = "splashscreen"
	SplashScreenComplete scene.SceneTransition = iota + 1
)

type Scene struct {
}

func NewScene() scene.Scene {
	return &Scene{}
}

func (s *Scene) EnterScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	logger.Info("Entering splash screen scene")
	return nil
}

func (s *Scene) ExitScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	logger.Info("Exiting splash screen scene")
	return nil
}

func (s *Scene) Update(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) (scene.SceneTransition, error) {
	return scene.ContinueScene, nil
}

func (s *Scene) Draw(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	return nil
}
