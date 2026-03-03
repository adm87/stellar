package splashscreen

import (
	"github.com/adm87/stellar/scene"
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

func (s *Scene) EnterScene(ctx scene.Context) error {
	ctx.Logger().Info("Entering splash screen scene")
	return nil
}

func (s *Scene) ExitScene(ctx scene.Context) error {
	ctx.Logger().Info("Exiting splash screen scene")
	return nil
}

func (s *Scene) Update(ctx scene.Context) (scene.SceneTransition, error) {
	return SplashScreenComplete, nil
}

func (s *Scene) Draw(ctx scene.Context) error {
	return nil
}
