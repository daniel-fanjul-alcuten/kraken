package main

import (
	"bytes"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/queue"
	"testing"
	"time"
)

func TestJobHandler1(t *testing.T) {

	queue := NewQueue(1)
	queue.GetInput() <- GoGetJob{ImportPath: "foo"}
	close(queue.GetInput())

	reader := &bytes.Buffer{}
	writer := &bytes.Buffer{}
	results := make(chan result, 1)
	handler := &jobHandler{queue, results}
	handler.handle(reader, writer)

	decoder := gob.NewDecoder(writer)
	var job GoGetJob
	if err := decoder.Decode(&job); err != nil {
		t.Error(err)
	}
	if job.ImportPath != "foo" {
		t.Error(job.ImportPath)
	}

	select {
	case res, ok := <-results:
		if !ok {
			t.Error(ok)
		}
		if res.job.ImportPath != "foo" {
			t.Error(res.job.ImportPath)
		}
		if res.ack {
			t.Error(res.ack)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}

func TestJobHandler2(t *testing.T) {

	queue := NewQueue(1)
	queue.GetInput() <- GoGetJob{ImportPath: "foo"}
	close(queue.GetInput())

	reader := &bytes.Buffer{}
	encoder := gob.NewEncoder(reader)
	if err := encoder.Encode(JobResult{true, "bar"}); err != nil {
		t.Error(err)
	}

	writer := &bytes.Buffer{}
	results := make(chan result, 1)
	handler := &jobHandler{queue, results}
	handler.handle(reader, writer)

	decoder := gob.NewDecoder(writer)
	var job GoGetJob
	if err := decoder.Decode(&job); err != nil {
		t.Error(err)
	}
	if job.ImportPath != "foo" {
		t.Error(job.ImportPath)
	}

	select {
	case res, ok := <-results:
		if !ok {
			t.Error(ok)
		}
		if res.job.ImportPath != "foo" {
			t.Error(res.job.ImportPath)
		}
		if !res.ack {
			t.Error(res.ack)
		}
		if !res.result.Success {
			t.Error(res.result.Success)
		}
		if res.result.Output != "bar" {
			t.Error(res.result.Output)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}
