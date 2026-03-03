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

// registerScenes registers all the scenes in the game with the provided Director, allowing them to be transitioned to during gameplay.
// It returns an error if there is an issue registering any of the scenes.
func registerScenes(d *scene.Director) error {
	if err := d.RegisterScene(splashscreen.SplashScreenScene, splashscreen.NewScene); err != nil {
		return err
	}
	if err := d.RegisterScene(gameplay.GameplayScene, gameplay.NewScene); err != nil {
		return err
	}
	return nil
}

// addSceneTransitions defines the valid transitions between scenes in the game, allowing the Director to manage scene changes based on specific conditions.
// It returns an error if there is an issue adding any of the scene transitions.
func addSceneTransitions(d *scene.Director) error {
	if err := d.AddTransition(splashscreen.SplashScreenScene, splashscreen.SplashScreenComplete, gameplay.GameplayScene); err != nil {
		return err
	}
	return nil
}

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

// NewShell creates a new Shell instance with the provided configuration, initializing all necessary components for the game.
func NewShell(config *Config) *Shell {
	return &Shell{
		config: config,
		assets: assets.NewAssets(),
		buffer: rendering.NewScreenBuffer(
			config.WindowWidth,
			config.WindowHeight,
			config.BackgroundColor,
			config.ResizeMode,
		),
		director: scene.NewDirector(),
		logger:   logging.NewLogger(),
		time: timing.NewTime(
			config.FPS,
		),
	}
}

// Assets returns the game's asset manager, which is responsible for loading and managing game assets such as images, sounds, and fonts.
func (s *Shell) Assets() *assets.Assets {
	return s.assets
}

// Buffer returns the game's screen buffer, which is used as a rendering target each frame.
func (s *Shell) Buffer() *rendering.ScreenBuffer {
	return s.buffer
}

// Logger returns the game's logger, which is used for logging messages and errors.
func (s *Shell) Logger() *logging.Logger {
	return s.logger
}

// Time returns the game's timing manager, which is responsible for tracking time and managing the game loop's timing.
func (s *Shell) Time() *timing.Time {
	return s.time
}

// Run starts the game loop, initializing the game and handling any errors that occur during execution.
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

// Update updates the game state, including the current scene and any necessary transitions.
// It returns an error if there is an issue updating the game state.
func (s *Shell) Update() error {
	s.time.Tick()
	return s.director.Update(s.assets, s.buffer, s.logger, s.time)
}

// Draw renders the current game state to the screen, including the current scene.
// It returns an error if there is an issue drawing the game state.
func (s *Shell) Draw(screen *ebiten.Image) {
	s.buffer.Clear()

	if err := s.director.Draw(s.assets, s.buffer, s.logger, s.time); err != nil {
		s.logger.Error("Error drawing scene: " + err.Error())
		ebitenutil.DebugPrint(screen, "Error drawing scene: "+err.Error())
		return
	}

	s.buffer.ApplyTo(screen)
}

// Layout calculates the layout of the game window based on the configured window size and render scale,
// returning the width and height for the game screen.
func (s *Shell) Layout(outsideWidth, outsideHeight int) (int, int) {
	width := int(float64(s.config.WindowWidth) * s.config.RenderScale)
	height := int(float64(s.config.WindowHeight) * s.config.RenderScale)
	return width, height
}
