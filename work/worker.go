package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
)

type Worker interface {
	Work(git *Git) JobResult
}
