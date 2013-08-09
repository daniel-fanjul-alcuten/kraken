package kraken

import (
	"bytes"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
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
	if _, err := git.Run(nil, "init"); err != nil {
		t.Error(err)
	}
	if _, err := git.Run(nil, "add", "-A"); err != nil {
		t.Error(err)
	}
	if _, err := git.Run(nil, "commit", "-m", "foo"); err != nil {
		t.Error(err)
	}

	dir2, err := td.NewDir()
	if err != nil {
		t.Error(err)
	}
	git2 := NewGit(dir2)
	if _, err := git2.Run(nil, "init", "--bare"); err != nil {
		t.Error(err)
	}

	refs, err := Push(git, dir2, "host", "repo", "master")
	if err != nil {
		t.Error(err)
	}

	if len(refs) != 1 {
		t.Error(len(refs))
	}
	if _, err := git2.Run(nil, "show", refs[0]+":foo"); err != nil {
		t.Error(err)
	}
	if output, err := git2.Run(nil, "cat-file", "tag", refs[0]); err != nil {
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
