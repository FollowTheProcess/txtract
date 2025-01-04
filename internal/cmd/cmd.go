// Package cmd implements the CLI for txtract.
package cmd

import (
	"fmt"

	"github.com/FollowTheProcess/cli"
)

var (
	version string // txtar version, set at compile time
	commit  string // commit hash of the source tree when this version of txtar was build
	date    string // build date

)

// Build constructs and returns the root txtract CLI command.
func Build() (*cli.Command, error) {
	txtract, err := cli.New(
		"txtract",
		cli.Allow(cli.NoArgs()),
		cli.Short("A CLI to interact with the txtar archive format"),
		cli.Version(version),
		cli.Commit(commit),
		cli.BuildDate(date),
		cli.SubCommands(buildZipCommand, buildUnzipCommand),
	)
	if err != nil {
		return nil, fmt.Errorf("could not build txtract cli: %w", err)
	}

	return txtract, nil
}
