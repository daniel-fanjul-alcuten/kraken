package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {

	queue := NewQueue(3)
	input := queue.GetInput()
	output := queue.GetOutput()

	select {
	case input <- GoGetJob{ImportPath: "foo"}:
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}

	select {
	case input <- GoGetJob{ImportPath: "bar"}:
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}

	close(input)

	select {
	case job := <-output:
		if job.ImportPath != "foo" {
			t.Error(job.ImportPath)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}

	select {
	case job := <-output:
		if job.ImportPath != "bar" {
			t.Error(job.ImportPath)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout")
	}
}
