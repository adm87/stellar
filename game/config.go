package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Config holds the configuration settings for the game.
type Config struct {
	Name            string                  // Name of the game
	Version         string                  // Version of the game
	RootDir         string                  // Root directory of the game
	LogLevel        string                  // Logging level (e.g., "debug", "info", "warn", "error")
	BuildMode       string                  // Build mode (e.g., "dev", "release")
	FPS             int                     // Target frames per second
	WindowWidth     int                     // Width of the game window
	WindowHeight    int                     // Height of the game window
	RenderScale     float64                 // Scale factor for rendering
	Fullscreen      bool                    // Whether to start the game in fullscreen mode
	BackgroundColor color.RGBA              // Background color of the game window
	DrawOptions     ebiten.DrawImageOptions // Default draw options for rendering
}

// NewDefaultConfig creates a new Config with default values, allowing overrides from command-line arguments.
func NewDefaultConfig(version string, args Args) *Config {
	return &Config{
		Name:            "Stellar",
		Version:         version,
		FPS:             60,
		WindowWidth:     1280,
		WindowHeight:    720,
		RenderScale:     1.0,
		Fullscreen:      args.Fullscreen,
		RootDir:         args.RootDir,
		LogLevel:        args.LogLevel,
		BackgroundColor: color.RGBA{R: 100, G: 149, B: 237, A: 255},
		DrawOptions: ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
		},
	}
}
