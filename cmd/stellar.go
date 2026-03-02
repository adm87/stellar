package cmd

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/adm87/stellar/errs"
	"github.com/adm87/stellar/game"
	"github.com/adm87/stellar/logging"
)

func Stellar(version string) error {
	args, err := parseArgs()

	if err != nil {
		return &errs.Fatal{
			Message: fmt.Sprintf("Failed to parse arguments: %v", err),
		}
	}

	cfg := game.NewConfig(version, args)

	if err := game.NewShell(cfg).Run(); err != nil {
		return &errs.Fatal{
			Message: fmt.Sprintf("Failed to run game shell: %v", err),
		}
	}

	return nil
}

func parseArgs() (game.Args, error) {
	args := game.Args{}

	flag.StringVar(&args.RootDir, "root", ".", "Root directory for game assets")
	flag.StringVar(&args.LogLevel, "log-level", "error", "Logging level (debug, info, warn, error)")
	flag.BoolVar(&args.Fullscreen, "fullscreen", false, "Start the game in fullscreen mode")
	flag.Parse()

	return args, validateArgs(&args)
}

func validateArgs(args *game.Args) error {
	absRootPath, err := filepath.Abs(args.RootDir)
	if err != nil {
		return &errs.InvalidArg{
			Message: fmt.Sprintf("Invalid root directory: %v", err),
		}
	}
	args.RootDir = absRootPath

	lvl := logging.LogLevel(args.LogLevel)
	if !lvl.IsValid() {
		return &errs.InvalidArg{
			Message: fmt.Sprintf("Invalid log level: %s", args.LogLevel),
		}
	}

	return nil
}
