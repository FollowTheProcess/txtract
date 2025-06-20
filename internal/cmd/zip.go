package cmd

import (
	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/txtract/internal/app"
)

const zipLong = `
The zip subcommand zips a filesystem directory into a txtar archive file
named <dir>.txtar saving the file to the current working directory.

The location of the saved archive can be controlled with the --output flag and the
name may be overridden using the --name flag.

The directory is traversed recursively by default so nested sub directories and files
will be included in the final archive, with the names of the files being set to their
path relative to the top level directory. Note that file paths are stored in the archive
with the unix separator ('/') regardless of the host platform.

If the directory does not exist or is empty, zip will return an error and exit.

The default behaviour is not to overwrite an existing archive, this can be controlled
by using the --force flag.
`

// buildZipCommand constructs and returns the zip subcommand.
func buildZipCommand() (*cli.Command, error) {
	var (
		output string
		name   string
		debug  bool
		force  bool
	)
	return cli.New(
		"zip",
		cli.Short("Zip an on-disk directory into a txtar archive"),
		cli.Long(zipLong),
		cli.Example("Zip up the testdata directory", "txtract zip ./testdata"),
		cli.Example("Save to another location", "txtract zip ./mydir --output ../somewhere/else"),
		cli.Example("Use a different name", "txtract zip ./mydir --name myarchive"),
		cli.Allow(cli.MaxArgs(1)), // Only 1 directory allowed
		cli.RequiredArg("dir", "The directory to zip into a txtar archive"),
		cli.Flag(&output, "output", 'o', ".", "Path to save the zipped txtar file"),
		cli.Flag(&name, "name", 'n', "", "Name of the txtar file, defaults to directory name"),
		cli.Flag(&force, "force", 'f', false, "Overwrite an existing archive"),
		cli.Flag(&debug, "debug", cli.NoShortHand, false, "Output debug info to stderr"),
		cli.Run(func(cmd *cli.Command, _ []string) error {
			txtract := app.New(cmd.Stdout(), cmd.Stderr(), debug)
			return txtract.Zip(cmd.Arg("dir"), name, output, force)
		}),
	)
}
