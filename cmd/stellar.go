package cmd

import (
	"flag"
	"path/filepath"

	"github.com/adm87/stellar/data"
	"github.com/adm87/stellar/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func Stellar(version string) error {
	a := data.GameArgs{}

	if err := parseArgs(&a); err != nil {
		return err
	}

	return ebiten.RunGame(game.NewShell(&data.GameConfig{
		Name:         "Stellar",
		Version:      version,
		RootDir:      a.RootDir,
		FPS:          60,
		WindowWidth:  800,
		WindowHeight: 600,
		RenderScale:  1.0,
		Fullscreen:   a.Fullscreen,
	}))
}

func parseArgs(args *data.GameArgs) error {
	flag.StringVar(&args.RootDir, "root", ".", "Root directory for game assets")
	flag.BoolVar(&args.Fullscreen, "fullscreen", false, "Start the game in fullscreen mode")
	flag.Parse()
	return validateArgs(args)
}

func validateArgs(args *data.GameArgs) error {
	absRootDir, err := filepath.Abs(args.RootDir)
	if err != nil {
		return err
	}
	args.RootDir = absRootDir

	return nil
}
