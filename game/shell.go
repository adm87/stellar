package game

import (
	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/errs"
	"github.com/adm87/stellar/logging"
	"github.com/adm87/stellar/scene"
	"github.com/adm87/stellar/timing"
	"github.com/hajimehoshi/ebiten/v2"
)

// --------------------------------------------------------------------------------
// Game Shell
// --------------------------------------------------------------------------------

// Shell is the main entry point for running the game. It manages the game loop and provides access to the game context.
type Shell struct {
	config   *Config
	director *scene.Director
	assets   *assets.Assets
	logger   *logging.Logger
	time     *timing.Time
}

func NewShell(config *Config) *Shell {
	return &Shell{
		config:   config,
		director: scene.NewDirector(),
		assets:   assets.NewAssets(),
		logger:   logging.NewLogger().With("version", config.Version),
		time:     timing.NewTime(config.FPS),
	}
}

func (s *Shell) Assets() *assets.Assets {
	return s.assets
}

func (s *Shell) Logger() *logging.Logger {
	return s.logger
}

func (s *Shell) Time() *timing.Time {
	return s.time
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

	s.time.Start()
	return ebiten.RunGame(s)
}

func (s *Shell) Update() error {
	s.time.Tick()
	return s.director.Update(s)
}

func (s *Shell) Draw(screen *ebiten.Image) {
	if err := s.director.Draw(s); err != nil {
		s.logger.Error("Error drawing scene: " + err.Error())
	}
}

func (s *Shell) Layout(outsideWidth, outsideHeight int) (int, int) {
	width := int(float64(s.config.WindowWidth) * s.config.RenderScale)
	height := int(float64(s.config.WindowHeight) * s.config.RenderScale)
	return width, height
}
