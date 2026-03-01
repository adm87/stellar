package game

import (
	"github.com/adm87/stellar/content"
	"github.com/adm87/stellar/data"
	"github.com/adm87/stellar/engine/assets"
	"github.com/adm87/stellar/engine/structures/store"
	"github.com/adm87/stellar/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Shell struct {
	config *data.GameConfig
	time   *Time

	testId store.StoreID
}

func NewShell(config *data.GameConfig) *Shell {
	loader := assets.NewLoader(content.EmbeddedFS, content.EmbeddedImage10x10)

	if err := loader.Load(); err != nil {
		panic("failed to load assets: " + err.Error())
	}

	testId, exists := images.GetStoreID(content.EmbeddedImage10x10)

	if !exists {
		panic("failed to retrieve test image from cache")
	}

	return &Shell{
		config: config,
		time:   NewTime(config.FPS),
		testId: testId,
	}
}

func (s *Shell) Update() error {
	s.time.Tick()
	return nil
}

func (s *Shell) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello")
	images.RenderImage(screen, s.testId, nil)
}

func (s *Shell) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	screenWidth = int(float64(s.config.WindowWidth) * s.config.RenderScale)
	screenHeight = int(float64(s.config.WindowHeight) * s.config.RenderScale)
	return screenWidth, screenHeight
}
