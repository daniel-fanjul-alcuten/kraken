package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestGoGetWorker(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir1, err := td.NewDir()
	if err != nil {
		t.Error(err)
	}

	if err := ioutil.WriteFile(path.Join(dir1, "foo.go"), []byte("package foo\ntype foo int\n"), 0666); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(path.Join(dir1, "foo_test.go"), []byte("package foo\nimport \"testing\"\nfunc TestX(t *testing.T){}\n"), 0666); err != nil {
		t.Error(err)
	}

	git1 := NewGit(dir1)
	if _, err := git1.Run(nil, "init"); err != nil {
		t.Error(err)
	}
	if _, err := git1.Run(nil, "add", "-A"); err != nil {
		t.Error(err)
	}
	if _, err := git1.Run(nil, "commit", "-m", "foo"); err != nil {
		t.Error(err)
	}

	dir2, err := td.NewDir()
	if err != nil {
		t.Error(err)
	}
	git2 := NewGit(dir2)
	if _, err := git2.Run(nil, "init"); err != nil {
		t.Error(err)
	}

	var w Worker = &GoGetWorker{GoGetJob{"foo", dir1, "master"}}
	result := w.Work(git2)
	if !result.Success {
		t.Error(result.Success)
	}
	if !strings.Contains(result.Output, "ok  \tfoo\t") {
		t.Error(result.Output)
	}
}
