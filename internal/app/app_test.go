package app_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/FollowTheProcess/test"
	"github.com/FollowTheProcess/txtract/internal/app"
)

func TestZip(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := io.Discard // Only logs go here anyway
	testApp := app.New(stdout, stderr, false)

	tmp, err := os.MkdirTemp("", "TestZip*")
	test.Ok(t, err)
	t.Cleanup(func() {
		os.RemoveAll(tmp)
	})

	target := filepath.Join("testdata", "TestZip")

	err = testApp.Zip(target, "", tmp, false)
	test.Ok(t, err)

	// Now should be a txtar file called TestZip.txtar inside tmp
	archive := filepath.Join(tmp, "TestZip.txtar")

	contents, err := os.ReadFile(archive)
	test.Ok(t, err)

	contents = bytes.ReplaceAll(contents, []byte("\r\n"), []byte("\n"))
	golden := filepath.Join("testdata", "TestZip.txtar")

	goldenContents, err := os.ReadFile(golden)
	test.Ok(t, err)

	test.DiffBytes(t, contents, goldenContents)
}
