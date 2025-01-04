package cmd

import (
	"fmt"

	"github.com/FollowTheProcess/cli"
)

const unzipLong = `
The unzip subcommand unpacks a txtar archive file into a real filesystem directory, files
inside the archive get unpacked into their on-disk equivalents with the contents as per
the archive file.

The archive is unpacked into a directory (under $CWD) named the same as the archive file, under which
all the contained files are written. This can be controlled with the --output flag.

The txtar format does not enforce that the file names are valid filesystem names, if any
archived files cannot be written to disk, unzip will return an error and exit.

If the output directory does not exist, it will be created automatically.
`

// buildUnzipCommand constructs and returns the unzip subcommand.
func buildUnzipCommand() (*cli.Command, error) {
	var output string
	return cli.New(
		"unzip",
		cli.Short("Unzip a txtar archive into the filesystem"),
		cli.Long(unzipLong),
		cli.Example("Unzip a txtar test case to testdata", "txtract unzip ./TestMyThing.txtar"),
		cli.Example(
			"Save to another location",
			"txtract unzip ./TestMyThing.txtar --output ../somewhere/else/",
		),
		cli.Allow(cli.MaxArgs(1)), // Only 1 txtar file allowed
		cli.Arg("archive", "The path to the txtar archive to unzip", ""),
		cli.Flag(
			&output,
			"output",
			'o',
			"",
			"Path to a directory under which to unzip the archive",
		),
		cli.Run(func(cmd *cli.Command, args []string) error {
			fmt.Fprintf(
				cmd.Stdout(),
				"unzip called with archive: %s, output: %s\n",
				cmd.Arg("archive"),
				output,
			)
			return nil
		}),
	)
}
