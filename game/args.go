package game

// Args holds the command-line arguments for the game.
type Args struct {
	RootDir    string // Root directory of the game
	LogLevel   string // Logging level (e.g., "debug", "info", "warn", "error")
	Fullscreen bool   // Whether to start the game in fullscreen mode
}
