// Package app implements the functionality of the txtract command, the CLI
// does the argument parsing and dispatches to this package to perform the
// requested action.
package app

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"go.followtheprocess.codes/txtar"
)

const (
	defaultFilePermissions = 0o644 // Default permissions for writing files, same as unix touch
	defaultDirPermissions  = 0o755 // Default permissions for creating directories, same as unix mkdir
)

// App represents the txtract application.
type App struct {
	stdout io.Writer    // Where normal application output is written
	stderr io.Writer    // Errors or logs are written here
	logger *slog.Logger // Logger logging to stderr
	debug  bool         // Output debug logs to stderr
}

// New returns a new [App].
func New(stdout, stderr io.Writer, debug bool) App {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	logger := slog.New(
		slog.NewTextHandler(stderr, &slog.HandlerOptions{
			Level: level,
		}),
	)
	return App{
		stdout: stdout,
		stderr: stderr,
		logger: logger,
		debug:  false,
	}
}

// ZipOptions are the options passed to the zip subcommand.
type ZipOptions struct {
	// Output is the path to save the zipped txtar file.
	Output string

	// Dir is the directory to zip into a txtar archive.
	Dir string

	// Name is the name of the txtar file, defaults to directory name.
	Name string

	// Force, if true, overwrites an existing archive.
	Force bool

	// Debug, if true, outputs debug logs.
	Debug bool
}

// Zip zips up a filesystem directory into a txtar archive named name under
// location.
//
// Force controls whether or not to overwrite the archive if it already exists.
func (a App) Zip(options ZipOptions) error {
	if options.Name == "" {
		// No override for name so use the target dir's name
		options.Name = filepath.Base(options.Dir)
	}

	outPath := filepath.Join(options.Output, options.Name)
	outPath += ".txtar"

	if exists(outPath) && !options.Force {
		return fmt.Errorf("path %s exists and will not be overwritten without --force", outPath)
	}

	a.log(
		"zipping dir into txtar archive",
		slog.String("dir", options.Dir),
		slog.String("archive", outPath),
		slog.Bool("force", options.Force),
	)

	archive, err := txtar.New()
	if err != nil {
		return err
	}

	err = filepath.WalkDir(options.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("WalkDirFunc err was not nil: %w", err)
		}

		if !d.IsDir() {
			contents, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("could not read %s: %w", path, err)
			}
			err = archive.Write(filepath.ToSlash(path), string(contents))
			if err != nil {
				return fmt.Errorf("could not add %s to archive: %w", path, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not walk directory tree from %s: %w", options.Dir, err)
	}

	// Now write the archive to the given location
	dir := filepath.Dir(outPath)
	if err = os.MkdirAll(dir, defaultDirPermissions); err != nil {
		return fmt.Errorf("could not create target directory %s: %w", dir, err)
	}

	file, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR, defaultFilePermissions)
	if err != nil {
		return fmt.Errorf("could not open target file %s: %w", outPath, err)
	}
	defer file.Close()

	if err := txtar.Dump(file, archive); err != nil {
		return fmt.Errorf("could not write archive to %s: %w", outPath, err)
	}

	return nil
}

// UnzipOptions are the options passed to the unzip subcommand.
type UnzipOptions struct {
	// Output is the based directory to unzip under, defaults to ".".
	Output string

	// Archive is the filepath to the txtar archive to unzip.
	Archive string

	// Force, if true, overwrites existing files and directories.
	Force bool

	// Debug, if true, outputs debug logs.
	Debug bool
}

// Unzip unzips a txtar archive back into real filesystem directories and files.
//
// Force controls whether or not to overwrite existing files and directories with
// the archive contents.
func (a App) Unzip(options UnzipOptions) error {
	a.log(
		"unzipping archive onto filesystem",
		slog.String("archive", options.Archive),
		slog.String("output", options.Output),
		slog.Bool("force", options.Force),
	)

	file, err := os.Open(options.Archive)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", options.Archive, err)
	}
	defer file.Close()

	archive, err := txtar.Parse(file)
	if err != nil {
		return fmt.Errorf("could not parse %s as txtar archive: %w", options.Archive, err)
	}

	// Name of the txtar archive file, to be used as the directory under which to unzip
	// the archive
	name := strings.TrimSuffix(filepath.Base(options.Archive), filepath.Ext(options.Archive))

	for path, contents := range archive.Files() {
		path = filepath.Join(options.Output, name, path)
		if exists(path) && !options.Force {
			return fmt.Errorf("path %s exists and will not be overwritten without --force", path)
		}
		// Ensure that if the path is nested, all the directories get created
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, defaultDirPermissions); err != nil {
			return fmt.Errorf("could not create directory %s: %w", dir, err)
		}

		if err := os.WriteFile(path, []byte(contents), defaultFilePermissions); err != nil {
			return fmt.Errorf("could not write to %s: %w", path, err)
		}
	}

	return nil
}

// log writes a debug log to stderr.
func (a App) log(msg string, attrs ...slog.Attr) {
	a.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

// exists reports whether a path exists in the filesystem.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
