package main

import (
	"bufio"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	"net"
	"testing"
	"time"
)

func TestListen(t *testing.T) {

	listener, err := net.ListenTCP("tcp", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)

	requests := make(chan Request, 1)
	errs := make(chan error, 1)
	go func() {
		errs <- listen(listener, requests)
	}()

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(Request{"rq", "req", "repo", "ref", 5, []GoJob{GoJob{"ip"}}}); err != nil {
		t.Error(err)
	}
	if err := writer.Flush(); err != nil {
		t.Error(err)
	}

	select {
	case req, ok := <-requests:
		if !ok {
			t.Error(ok)
		}
		if req.Repoquest != "rq" {
			t.Error(req.Repoquest)
		}
		if req.Request != "req" {
			t.Error(req.Request)
		}
		if req.Repository != "repo" {
			t.Error(req.Repository)
		}
		if req.Reference != "ref" {
			t.Error(req.Reference)
		}
		if req.Time != 5 {
			t.Error(req.Time)
		}
		if len(req.Jobs) != 1 {
			t.Error(len(req.Jobs))
		} else if req.Jobs[0].ImportPath != "ip" {
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

	select {
	// "use of closed network connection" expected
	case _, ok := <-errs:
		if !ok {
			t.Error(ok)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}
