package main

import (
	"context"
	"os"
	"os/signal"

	"go.xargs.dev/bindl/command/cli"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cli.Root.AddCommand(versionCmd)
	err := cli.Root.ExecuteContext(ctx)
	return err
}
