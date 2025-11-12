package cmd

import (
	"context"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/txtract/internal/app"
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
	var options app.UnzipOptions
	return cli.New(
		"unzip",
		cli.Short("Unzip a txtar archive into the filesystem"),
		cli.Long(unzipLong),
		cli.Example("Unzip a txtar test case to testdata", "txtract unzip ./TestMyThing.txtar"),
		cli.Example("Save to another location", "txtract unzip ./TestMyThing.txtar --output ../somewhere/else/"),
		cli.Arg(&options.Archive, "archive", "The path to the txtar archive to unzip"),
		cli.Flag(&options.Output, "output", 'o', "Base directory to unzip under", cli.FlagDefault(".")),
		cli.Flag(&options.Force, "force", 'f', "Overwrite existing files and directories"),
		cli.Flag(&options.Debug, "debug", flag.NoShortHand, "Output debug info to stderr"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			txtract := app.New(cmd.Stdout(), cmd.Stderr(), options.Debug)
			return txtract.Unzip(ctx, options)
		}),
	)
}
