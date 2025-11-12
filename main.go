package main

import (
	"context"
	"os"

	"go.followtheprocess.codes/msg"
	"go.followtheprocess.codes/txtract/internal/cmd"
)

func main() {
	if err := run(); err != nil {
		msg.Err(err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	cli, err := cmd.Build()
	if err != nil {
		return err
	}

	return cli.Execute(ctx)
}
