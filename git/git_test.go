package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	"strings"
	"testing"
)

func TestGitDir(t *testing.T) {

	git := NewGit("foo")
	if git.Dir() != "foo" {
		t.Error(git.Dir())
	}
}

func TestGitString(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Fatal(err)
	}

	git := NewGit(dir)
	if _, err := git.Run(nil, "init", "--bare"); err != nil {
		t.Error(err)
	}

	if output, err := git.String(nil, "config", "core.bare"); err != nil {
		t.Error(err)
	} else if output != "true" {
		t.Error(output)
	}
}

func TestGitStringError(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Fatal(err)
	}

	git := NewGit(dir)
	if _, err := git.Run(nil, "init", "--bare"); err != nil {
		t.Error(err)
	}

	if output, err := git.String(nil, "config", "-x", "foo.bar"); err == nil {
		t.Error(err)
	} else if !strings.Contains(output, "usage: git config [options]") {
		t.Error(output)
	}
}
