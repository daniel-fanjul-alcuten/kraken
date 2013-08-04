package main

import (
	"bufio"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	. "github.com/daniel-fanjul-alcuten/kraken/push"
	"io/ioutil"
	"net"
	"path"
	"testing"
	"time"
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

	listener, err := net.ListenTCP("tcp", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)

	channel := make(chan interface{})
	go func() {
		defer close(channel)
		conn, err := listener.Accept()
		if err != nil {
			channel <- err
			return
		}
		defer conn.Close()
		reader := bufio.NewReader(conn)
		decoder := gob.NewDecoder(reader)
		var request Request
		if err := decoder.Decode(&request); err != nil {
			channel <- err
			return
		}
		if err := conn.Close(); err != nil {
			channel <- err
			return
		}
		channel <- request
	}()

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = submit(git2, dir2, conn, requests...)
	if err != nil {
		t.Error(err)
	}

	select {
	case v := <-channel:
		req, ok := v.(Request)
		if !ok {
			t.Errorf("%#v", req)
		}
		if req.Repoquest != dir2 {
			t.Error(req.Repoquest)
		}
		if req.Request != requests[0] {
			t.Error(req.Request)
		}
		if req.Repository != dir {
			t.Error(req.Repository)
		}
		if req.Reference != "refs/heads/master" {
			t.Error(req.Reference)
		}
		if req.Time < 1375609864 {
			t.Error(req.Time)
		}
		if len(req.Jobs) != 1 {
			t.Error(len(req.Jobs))
		} else if req.Jobs[0].ImportPath != "foo/bar/baz" {
			t.Error(req.Jobs[0].ImportPath)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}

	if err := conn.Close(); err != nil {
		t.Error(err)
	}

	if err := listener.Close(); err != nil {
		t.Error(err)
	}
}
