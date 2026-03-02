package main

import (
	"os"

	"github.com/adm87/stellar/cmd"
)

var version = "0.0.0-unreleased"

func main() {
	if err := cmd.Stellar(version); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
