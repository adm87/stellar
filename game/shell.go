package game

import (
	"github.com/adm87/stellar/assets"
	"github.com/adm87/stellar/errs"
	"github.com/adm87/stellar/game/scenes/gameplay"
	"github.com/adm87/stellar/game/scenes/splashscreen"
	"github.com/adm87/stellar/logging"
	"github.com/adm87/stellar/rendering"
	"github.com/adm87/stellar/scene"
	"github.com/adm87/stellar/timing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// --------------------------------------------------------------------------------
// Game Shell
// --------------------------------------------------------------------------------

// Shell is the main entry point for the game, responsible for initializing and running the game loop.
type Shell struct {
	config   *Config
	assets   *assets.Assets
	buffer   *rendering.ScreenBuffer
	director *scene.Director
	logger   *logging.Logger
	time     *timing.Time
}

func NewShell(config *Config) *Shell {
	return &Shell{
		config: config,
		assets: assets.NewAssets(),
		buffer: rendering.NewScreenBuffer(
			config.WindowWidth,
			config.WindowHeight,
			config.BackgroundColor,
		),
		director: scene.NewDirector(),
		logger:   logging.NewLogger(),
		time: timing.NewTime(
			config.FPS,
		),
	}
}

func (s *Shell) Assets() *assets.Assets {
	return s.assets
}

func (s *Shell) Buffer() *rendering.ScreenBuffer {
	return s.buffer
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
		"version", s.config.Version,
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

	if err := registerScenes(s.director); err != nil {
		return errs.BootFailure{
			Message: "Failed to register scenes: " + err.Error(),
		}
	}

	if err := addSceneTransitions(s.director); err != nil {
		return errs.BootFailure{
			Message: "Failed to add scene transitions: " + err.Error(),
		}
	}

	ebiten.SetWindowTitle(s.config.Name + " - " + s.config.Version + " (" + s.config.BuildMode + ")")
	ebiten.SetWindowSize(s.config.WindowWidth, s.config.WindowHeight)
	ebiten.SetFullscreen(s.config.Fullscreen)

	s.director.TransitionTo(splashscreen.SplashScreenScene)
	s.time.Start()

	return ebiten.RunGame(s)
}

func (s *Shell) Update() error {
	s.time.Tick()
	return s.director.Update(s.assets, s.buffer, s.logger, s.time)
}

func (s *Shell) Draw(screen *ebiten.Image) {
	s.buffer.Clear()

	if err := s.director.Draw(s.assets, s.buffer, s.logger, s.time); err != nil {
		s.logger.Error("Error drawing scene: " + err.Error())
		ebitenutil.DebugPrint(screen, "Error drawing scene: "+err.Error())
		return
	}
}

func (s *Shell) Layout(outsideWidth, outsideHeight int) (int, int) {
	width := int(float64(s.config.WindowWidth) * s.config.RenderScale)
	height := int(float64(s.config.WindowHeight) * s.config.RenderScale)
	return width, height
}

func registerScenes(d *scene.Director) error {
	if err := d.RegisterScene(splashscreen.SplashScreenScene, splashscreen.NewScene); err != nil {
		return err
	}
	if err := d.RegisterScene(gameplay.GameplayScene, gameplay.NewScene); err != nil {
		return err
	}
	return nil
}

func addSceneTransitions(d *scene.Director) error {
	if err := d.AddTransition(splashscreen.SplashScreenScene, splashscreen.SplashScreenComplete, gameplay.GameplayScene); err != nil {
		return err
	}
	return nil
}
