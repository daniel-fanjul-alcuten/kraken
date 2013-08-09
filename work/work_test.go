package kraken

import (
	"bytes"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	"io"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"testing"
)

func TestWork(t *testing.T) {

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

	reader := &bytes.Buffer{}
	encoder := gob.NewEncoder(reader)
	if err := encoder.Encode(GoGetJob{"foo", dir1, "master"}); err != nil {
		t.Error(err)
	}

	writer := &bytes.Buffer{}
	logger := log.New(ioutil.Discard, "", 0)
	if err := Work(reader, writer, logger); err != io.EOF {
		t.Error(err)
	}

	decoder := gob.NewDecoder(writer)
	var result JobResult
	if err := decoder.Decode(&result); err != nil {
		t.Error(err)
	}
	if !result.Success {
		t.Error(result.Success)
	}
	if !strings.Contains(result.Output, "ok  \tfoo\t") {
		t.Error(result.Output)
	}
}
