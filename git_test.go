package kraken

import (
	"os"
	"path"
	"testing"
)

func TestGitInit(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Fatal(err)
	}

	git := NewGit(dir)
	if err := git.Init(); err != nil {
		t.Error(err)
	}

	name := path.Join(dir, ".git", "config")
	if _, err := os.Stat(name); err != nil {
		t.Error(err)
	}
}
