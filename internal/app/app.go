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

	"github.com/FollowTheProcess/txtar"
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

// log writes a debug log to stderr.
func (a App) log(msg string, attrs ...slog.Attr) {
	a.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
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

	path := filepath.Join(location, name)
	path += ".txtar"

	a.log(
		"zipping dir into txtar archive",
		slog.String("dir", target),
		slog.String("archive", path),
		slog.Bool("force", force),
	)

	archive, err := txtar.New()
	if err != nil {
		return err
	}

	err = filepath.WalkDir(target, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("WalkDirFunc err was not nil: %w", err)
		}

		if !d.IsDir() {
			contents, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("could not read %s: %w", path, err)
			}
			err = archive.Add(filepath.ToSlash(path), contents)
			if err != nil {
				return fmt.Errorf("could not add %s to archive: %w", path, err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("could not walk directory tree from %s: %w", target, err)
	}

	// Now write the archive to the given location
	dir := filepath.Dir(path)
	if err = os.MkdirAll(dir, defaultDirPermissions); err != nil {
		return fmt.Errorf("could not create target directory %s: %w", dir, err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, defaultFilePermissions)
	if err != nil {
		return fmt.Errorf("could not open target file %s: %w", path, err)
	}
	defer file.Close()

	if err := txtar.Dump(file, archive); err != nil {
		return fmt.Errorf("could not write archive to %s: %w", path, err)
	}

	return nil
}

// Unzip unzips a txtar archive back into real filesystem directories and files.
//
// Force controls whether or not to overwrite existing files and directories with
// the archive contents.
func (a App) Unzip(target, location string, force bool) error {
	a.log(
		"unzipping archive onto filesystem",
		slog.String("archive", target),
		slog.String("location", location),
		slog.Bool("force", force),
	)
	return nil
}
