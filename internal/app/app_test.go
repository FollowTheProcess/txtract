package app_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"go.followtheprocess.codes/test"
	"go.followtheprocess.codes/txtar"
	"go.followtheprocess.codes/txtract/internal/app"
)

func TestZip(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := io.Discard // Only logs go here anyway
	testApp := app.New(stdout, stderr, false)

	tmp := t.TempDir()
	target := filepath.Join("testdata", "TestZip")

	options := app.ZipOptions{
		Dir:    target,
		Output: tmp,
		Force:  false,
		Debug:  false,
	}

	err := testApp.Zip(t.Context(), options)
	test.Ok(t, err)

	// Now should be a txtar file called TestZip.txtar inside tmp
	archive := filepath.Join(tmp, "TestZip.txtar")

	contents, err := os.ReadFile(archive)
	test.Ok(t, err)

	contents = bytes.ReplaceAll(contents, []byte("\r\n"), []byte("\n"))
	golden := filepath.Join("testdata", "TestZip.txtar")

	goldenContents, err := os.ReadFile(golden)
	goldenContents = bytes.ReplaceAll(goldenContents, []byte("\r\n"), []byte("\n"))
	test.Ok(t, err)

	test.DiffBytes(t, contents, goldenContents)
}

func TestUnzip(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := io.Discard // Only logs go here anyway
	testApp := app.New(stdout, stderr, false)

	tmp := t.TempDir()

	archiveFile := filepath.Join("testdata", "TestUnzip.txtar")
	file, err := os.Open(archiveFile)
	test.Ok(t, err)
	defer file.Close()

	archive, err := txtar.Parse(file)
	test.Ok(t, err)

	options := app.UnzipOptions{
		Output:  tmp,
		Archive: archiveFile,
		Force:   false,
		Debug:   false,
	}

	err = testApp.Unzip(t.Context(), options)
	test.Ok(t, err)

	// We should now have real files and directories, all under location
	for path, wantContents := range archive.Files() {
		path = filepath.Join(tmp, "TestUnzip", path)
		gotContents, err := os.ReadFile(path)
		test.Ok(t, err)

		test.Diff(t, string(gotContents), wantContents)
	}
}
