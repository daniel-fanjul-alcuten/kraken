package main

import (
	"bytes"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	"testing"
	"time"
)

func TestRequestHandler(t *testing.T) {

	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(Request{"rq", "req", "repo", "ref", 5, []GoGetRequest{GoGetRequest{"ip"}}}); err != nil {
		t.Error(err)
	}

	requests := make(chan Request, 1)
	handler := &requestHandler{requests}
	handler.handle(buffer, nil)

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
}
