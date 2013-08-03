package main

import (
	"bytes"
	. "github.com/daniel-fanjul-alcuten/kraken"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestPush(t *testing.T) {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		t.Error(err)
	}

	if err := ioutil.WriteFile(path.Join(dir, "foo"), []byte("bar"), 0666); err != nil {
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

	refs, err := push(git, dir2, "host", "repo", "master")
	if err != nil {
		t.Error(err)
	}

	if len(refs) != 1 {
		t.Error(len(refs))
	}
	if err := git2.Cmd("show", refs[0]+":foo").Run(); err != nil {
		t.Error(err)
	}
	if output, err := git2.Cmd("cat-file", "tag", refs[0]).Output(); err != nil {
		t.Error(err)
	} else {
		buffer := bytes.NewBuffer(output)
		buffer.ReadString('\n') // remove object
		buffer.ReadString('\n') // remove type
		buffer.ReadString('\n') // remove tag
		buffer.ReadString('\n') // remove tagger
		buffer.ReadString('\n') // remove empty line
		json := strings.TrimSpace(buffer.String())
		if json != `{"Repository":"repo","Reference":"refs/heads/master"}` {
			t.Error(json)
		}
	}
}
