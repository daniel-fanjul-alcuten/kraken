package kraken

import (
	"bytes"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	. "github.com/daniel-fanjul-alcuten/kraken/push"
	"io/ioutil"
	"path"
	"testing"
)

func TestSubmit(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Error(err)
	}

	json := `{"Jobs":[{"ImportPath":"foo/bar/baz"}]}`
	if err := ioutil.WriteFile(path.Join(dir, "kraken.json"), []byte(json), 0666); err != nil {
		t.Error(err)
	}

	git := NewGit(dir)
	if err := git.Init(); err != nil {
		t.Error(err)
	}
	if err := git.Cmd("add", "-A").Run(); err != nil {
		t.Error(err)
	}
	if err := git.Cmd("commit", "-m", "foo").Run(); err != nil {
		t.Error(err)
	}

	dir2, err := td.NewDir()
	if err != nil {
		t.Error(err)
	}
	git2 := NewGit(dir2)
	if err := git2.InitBare(); err != nil {
		t.Error(err)
	}

	requests, err := Push(git, dir2, "host", dir, "master")
	if err != nil {
		t.Error(err)
	}

	buffer := &bytes.Buffer{}
	err = Submit(git2, dir2, buffer, requests...)
	if err != nil {
		t.Error(err)
	}

	decoder := gob.NewDecoder(buffer)
	var request Request
	if err := decoder.Decode(&request); err != nil {
		t.Error(err)
	}

	if request.Repoquest != dir2 {
		t.Error(request.Repoquest)
	}
	if request.Request != requests[0] {
		t.Error(request.Request)
	}
	if request.Repository != dir {
		t.Error(request.Repository)
	}
	if request.Reference != "refs/heads/master" {
		t.Error(request.Reference)
	}
	if request.Time < 1375609864 {
		t.Error(request.Time)
	}
	if len(request.Jobs) != 1 {
		t.Error(len(request.Jobs))
	} else if request.Jobs[0].ImportPath != "foo/bar/baz" {
		t.Error(request.Jobs[0].ImportPath)
	}
}
