package gameplay

import "github.com/adm87/stellar/scene"

const GameplayScene scene.SceneID = "gameplay"

type Scene struct {
}

func NewScene() scene.Scene {
	return &Scene{}
}

func (s *Scene) EnterScene(ctx scene.Context) error {
	return nil
}

func (s *Scene) ExitScene(ctx scene.Context) error {
	return nil
}

func (s *Scene) Update(ctx scene.Context) error {
	return nil
}

func (s *Scene) Draw(ctx scene.Context) error {
	return nil
}
