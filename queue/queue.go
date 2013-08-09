package kraken

import (
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
)

type Queue struct {
	jobs   []GoGetJob
	input  chan GoGetJob
	output chan GoGetJob
}

func NewQueue(size int) *Queue {
	jobs := make([]GoGetJob, 0, size)
	input := make(chan GoGetJob, size)
	output := make(chan GoGetJob)
	q := &Queue{jobs, input, output}
	go q.run()
	return q
}

func (q *Queue) run() {
	input := q.input
	for {
		output := q.output
		var first GoGetJob
		if len(q.jobs) > 0 {
			first = q.jobs[0]
		} else if input == nil {
			break
		} else {
			output = nil
		}
		select {
		case job, ok := <-input:
			if ok {
				q.jobs = append(q.jobs, job)
			} else {
				input = nil
			}
		case output <- first:
			q.jobs = q.jobs[1:]
		}
	}
	close(q.output)
}

func (q *Queue) GetInput() chan<- GoGetJob {
	return q.input
}

func (q *Queue) GetOutput() <-chan GoGetJob {
	return q.output
}
