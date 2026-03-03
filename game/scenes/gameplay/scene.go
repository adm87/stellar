package gameplay

import (
	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/logging"
	"github.com/adm87/stellar/rendering"
	"github.com/adm87/stellar/scene"
	"github.com/adm87/stellar/timing"
)

const GameplayScene scene.SceneID = "gameplay"

type Scene struct {
}

func NewScene() scene.Scene {
	return &Scene{}
}

func (s *Scene) EnterScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	logger.Info("Entering gameplay scene")
	return nil
}

func (s *Scene) ExitScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	logger.Info("Exiting gameplay scene")
	return nil
}

func (s *Scene) Update(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) (scene.SceneTransition, error) {
	return scene.ContinueScene, nil
}

func (s *Scene) Draw(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	return nil
}
