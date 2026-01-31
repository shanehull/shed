package main

import (
	"context"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v3"
)

var bashCompletionsMode bool

func main() {
	if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
		bashCompletionsMode = true
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

var app = &cli.Command{
	Usage:                 "The Shed toolbox.",
	Description:           "A toolbox of tools for The Shed.",
	EnableShellCompletion: true,
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		if bashCompletionsMode {
			return ctx, nil
		}

		return nil, nil
	},
	Commands: []*cli.Command{
		zetCommand,
		checkcrtCommand,
	},
}
