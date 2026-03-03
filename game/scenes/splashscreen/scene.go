package splashscreen

import (
	"fmt"

	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/content"
	"github.com/adm87/stellar/logging"
	"github.com/adm87/stellar/rendering"
	"github.com/adm87/stellar/scene"
	"github.com/adm87/stellar/timing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

const (
	SplashScreenScene    scene.SceneID         = "splashscreen"
	SplashScreenComplete scene.SceneTransition = iota + 1
)

type Scene struct {
	img      *ebiten.Image
	sequence *gween.Sequence

	opacity float32
}

func NewScene() scene.Scene {
	return &Scene{}
}

func (s *Scene) EnterScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	if err := assets.Load(content.EmbeddedFS, content.EmbeddedSplash1920x1080); err != nil {
		logger.Error(fmt.Sprintf("Failed to load splash screen asset: %s", err.Error()))
		return nil
	}

	img, err := assets.Images().GetByPath(content.EmbeddedSplash1920x1080)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get splash screen image: %s", err.Error()))
		return nil
	}

	buffer.Resize(img.Bounds().Dx(), img.Bounds().Dy())
	buffer.SetFilter(ebiten.FilterLinear)

	s.sequence = gween.NewSequence(
		gween.New(0, 1, 1.0, ease.Linear),
		gween.New(1, 1, 2.0, ease.Linear),
		gween.New(1, 0, 1.0, ease.Linear),
	)
	s.img = img

	return nil
}

func (s *Scene) ExitScene(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	s.sequence = nil
	s.img = nil

	if err := assets.Unload(content.EmbeddedSplash1920x1080); err != nil {
		logger.Error(fmt.Sprintf("Failed to unload splash screen asset: %s", err.Error()))
	}

	return nil
}

func (s *Scene) Update(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) (scene.SceneTransition, error) {
	if s.img == nil {
		logger.Error("Splash screen image is nil, cannot update splash screen")
		return SplashScreenComplete, nil
	}

	value, _, complete := s.sequence.Update(time.Delta32())

	if complete {
		return SplashScreenComplete, nil
	}

	s.opacity = value

	return scene.ContinueScene, nil
}

func (s *Scene) Draw(assets *assets.Assets, buffer *rendering.ScreenBuffer, logger *logging.Logger, time *timing.Time) error {
	if s.img == nil {
		logger.Error("Splash screen image is nil, cannot draw splash screen")
		return nil
	}

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(s.opacity)

	buffer.Image().DrawImage(s.img, op)

	return nil
}
