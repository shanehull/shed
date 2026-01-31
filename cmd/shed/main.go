package main

import (
	"context"
	"fmt"
	"os"
	"time"

	cli "github.com/urfave/cli/v2"
)

var bashCompletionsMode bool

func main() {
	if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
		bashCompletionsMode = true
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var app = &cli.App{
	Usage:                "The Shed toolbox.",
	Description:          "A toolbox of tools for The Shed.",
	EnableBashCompletion: true,
	Compiled:             time.Now(),
	Before: func(_ *cli.Context) error {
		if bashCompletionsMode {
			return nil
		}

		return nil
	},
	Commands: []*cli.Command{
		zetCommand,
		checkcrtCommand,
	},
}
