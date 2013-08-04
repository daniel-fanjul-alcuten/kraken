package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	"os"
	"path"
	"strings"
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

func TestGitInitBare(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Fatal(err)
	}

	git := NewGit(dir)
	if err := git.InitBare(); err != nil {
		t.Error(err)
	}

	name := path.Join(dir, "config")
	if _, err := os.Stat(name); err != nil {
		t.Error(err)
	}
}

func TestGitCmd(t *testing.T) {

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

	output, err := git.Cmd("config", "core.bare").Output()
	value := strings.TrimSpace(string(output))
	if value != "false" {
		t.Error(value)
	}
	if err != nil {
		t.Error(err)
	}
}
