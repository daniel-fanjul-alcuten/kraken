package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestGitDir(t *testing.T) {

	git := NewGit("foo")
	if git.Dir() != "foo" {
		t.Error(git.Dir())
	}
}

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

func TestGitInitError(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Fatal(err)
	}

	git := NewGitFullPath("git-unknown", dir)
	if err := git.Init(); err == nil {
		t.Error()
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

func TestGitInitBareError(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Fatal(err)
	}

	git := NewGitFullPath("git-unknown", dir)
	if err := git.InitBare(); err == nil {
		t.Error(err)
	}
}

func TestGitConfig(t *testing.T) {

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

	if err := git.SetConfig("foo.bar", "baz"); err != nil {
		t.Error(err)
	}
	if value, err := git.Config("foo.bar"); err != nil {
		t.Error(err)
	} else if value != "baz" {
		t.Error(value)
	}
}

func TestGitConfigError(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Fatal(err)
	}

	git := NewGitFullPath("git-unknown", dir)
	if err := git.InitBare(); err == nil {
		t.Error(err)
	}

	if err := git.SetConfig("foo.bar", "baz"); err == nil {
		t.Error(err)
	}
	if _, err := git.Config("foo.bar"); err == nil {
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
