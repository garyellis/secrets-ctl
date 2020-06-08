package main

import (
	"os"

	"github.com/garyellis/secrets-ctl/pkg/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
