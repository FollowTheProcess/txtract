// Package app implements the functionality of the txtract command, the CLI
// does the argument parsing and dispatches to this package to perform the
// requested action.
package app

import (
	"fmt"
	"io"
	"path/filepath"
)

// App represents the txtract application.
type App struct {
	stdout io.Writer // Where normal application output is written
	stderr io.Writer // Errors or logs are written here
	debug  bool      // Output debug logs to stderr
}

// New returns a new [App].
func New(stdout, stderr io.Writer, debug bool) App {
	return App{
		stdout: stdout,
		stderr: stderr,
		debug:  false,
	}
}

// Zip zips up a filesystem directory into a txtar archive named name under
// location.
//
// Force controls whether or not to overwrite the archive if it already exists.
func (a App) Zip(target, name, location string, force bool) error {
	if name == "" {
		// No override for name so use the target dir
		name = target
	}

	fmt.Fprintf(a.stdout, "zipping %s into %s.txtar, force: %v\n", target, filepath.Join(location, name), force)
	return nil
}

// Unzip unzips a txtar archive back into real filesystem directories and files.
//
// Force controls whether or not to overwrite existing files and directories with
// the archive contents.
func (a App) Unzip(target, location string, force bool) error {
	fmt.Fprintf(a.stdout, "unzipping %s.txtar into %s, force: %v\n", target, location, force)
	return nil
}
