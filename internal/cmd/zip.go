package cmd

import (
	"fmt"

	"github.com/FollowTheProcess/cli"
)

const zipLong = `
The zip subcommand zips a filesystem directory into a txtar archive file
named <dir>.txtar saving the file to the current working directory.

The location of the saved archive can be controlled with the --output flag.

The directory is traversed recursively by default so nested sub directories and files
will be included in the final archive, with the names of the files being set to their
path relative to the top level directory.

If the directory does not exist or is empty, zip will return an error and exit.
`

// buildZipCommand constructs and returns the zip subcommand.
func buildZipCommand() (*cli.Command, error) {
	var output string
	return cli.New(
		"zip",
		cli.Short("Zip an on-disk directory into a txtar archive"),
		cli.Long(zipLong),
		cli.Example("Zip up the testdata directory", "txtract zip ./testdata"),
		cli.Example(
			"Save to another location",
			"txtract zip ./mydir --output ../somewhere/else.txtar",
		),
		cli.Allow(cli.MaxArgs(1)), // Only 1 directory allowed
		cli.Arg("dir", "The directory to zip into a txtar archive", ""),
		cli.Flag(&output, "output", 'o', "", "Path to save the zipped txtar file"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintf(
				cmd.Stdout(),
				"zip called with dir: %s, output: %s\n",
				cmd.Arg("dir"),
				output,
			)
			return nil
		}),
	)
}
