package game

// Config holds the configuration settings for the game.
type Config struct {
	Name         string  // Name of the game
	Version      string  // Version of the game
	RootDir      string  // Root directory of the game
	LogLevel     string  // Logging level (e.g., "debug", "info", "warn", "error")
	BuildMode    string  // Build mode (e.g., "dev", "release")
	FPS          int     // Target frames per second
	WindowWidth  int     // Width of the game window
	WindowHeight int     // Height of the game window
	RenderScale  float64 // Scale factor for rendering
	Fullscreen   bool    // Whether to start the game in fullscreen mode
}

func NewDefaultConfig(version string, args Args) *Config {
	return &Config{
		Name:         "Stellar",
		Version:      version,
		FPS:          60,
		WindowWidth:  1280,
		WindowHeight: 720,
		RenderScale:  1.0,
		Fullscreen:   args.Fullscreen,
		RootDir:      args.RootDir,
		LogLevel:     args.LogLevel,
	}
}
