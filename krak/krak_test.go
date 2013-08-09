package main

import (
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	"io/ioutil"
	"path"
	"testing"
)

func TestKrak(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Error(err)
	}

	if err := ioutil.WriteFile(path.Join(dir, "kraken.json"), []byte(`{"Jobs":[{"ImportPath":"foo"}]}`), 0666); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(path.Join(dir, "foo.go"), []byte("package foo\ntype foo int\n"), 0666); err != nil {
		t.Error(err)
	}

	git := NewGit(dir)
	if _, err := git.Run(nil, "init"); err != nil {
		t.Error(err)
	}
	if _, err := git.Run(nil, "add", "-A"); err != nil {
		t.Error(err)
	}
	if _, err := git.Run(nil, "commit", "-m", "foo"); err != nil {
		t.Error(err)
	}

	if err := krak(dir, "master"); err != nil {
		t.Error(err)
	}

	if _, err := git.Run(nil, "show", "refs/build/master:kraken.log"); err != nil {
		t.Error(err)
	}

	if _, err := git.Run(nil, "show", "refs/build/master:pkg"); err != nil {
		t.Error(err)
	}
}
