package main

import (
	"os"

	"github.com/adm87/stellar/cmd"
	"github.com/adm87/stellar/images"
)

var version = "0.0.0-unreleased"

func init() {
	images.Register()
}

func main() {
	if err := cmd.Stellar(version); err != nil {
		os.Exit(1)
	}
}
