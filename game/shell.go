package game

import (
	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/errs"
	"github.com/adm87/stellar/logging"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Shell struct {
	config *Config
	assets *assets.Assets
	logger *logging.Logger
	time   *Time
}

func NewShell(config *Config) *Shell {
	return &Shell{
		config: config,
		assets: assets.NewAssets(),
		logger: logging.NewLogger().With("version", config.Version),
		time:   NewTime(config.FPS),
	}
}

func (s *Shell) Run() error {
	s.logger.SetLevel(logging.LogLevel(s.config.LogLevel))
	s.logger.Debug("Starting game shell...",
		"buildMode", s.config.BuildMode,
		"name", s.config.Name,
		"rootDir", s.config.RootDir,
		"fullscreen", s.config.Fullscreen,
		"windowWidth", s.config.WindowWidth,
		"windowHeight", s.config.WindowHeight,
		"renderScale", s.config.RenderScale,
		"fullscreen", s.config.Fullscreen,
		"logLevel", s.config.LogLevel,
		"fps", s.config.FPS,
	)

	if err := s.assets.Initialize(); err != nil {
		return errs.BootFailure{
			Message: "Failed to register asset loaders: " + err.Error(),
		}
	}

	ebiten.SetWindowTitle(s.config.Name + " - " + s.config.Version + " (" + s.config.BuildMode + ")")
	ebiten.SetWindowSize(s.config.WindowWidth, s.config.WindowHeight)
	ebiten.SetFullscreen(s.config.Fullscreen)

	s.time.start()
	return ebiten.RunGame(s)
}

func (s *Shell) Update() error {
	s.time.tick()
	return nil
}

func (s *Shell) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, Stellar!")
}

func (s *Shell) Layout(outsideWidth, outsideHeight int) (int, int) {
	width := int(float64(s.config.WindowWidth) * s.config.RenderScale)
	height := int(float64(s.config.WindowHeight) * s.config.RenderScale)
	return width, height
}
