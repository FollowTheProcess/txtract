package main

import (
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
	cli, err := cmd.Build()
	if err != nil {
		return err
	}

	return cli.Execute()
}
