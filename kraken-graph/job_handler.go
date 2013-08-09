package main

import (
	"bufio"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/queue"
	"io"
)

type result struct {
	job    GoGetJob
	ack    bool
	result JobResult
}

type jobHandler struct {
	queue   *Queue
	results chan<- result
}

func (h *jobHandler) handle(reader io.Reader, writer io.Writer) error {
	input := h.queue.GetOutput()

	bufreader := bufio.NewReader(reader)
	decoder := gob.NewDecoder(bufreader)

	bufwriter := bufio.NewWriter(writer)
	encoder := gob.NewEncoder(bufwriter)

	for job := range input {
		var jobResult JobResult

		if err := encoder.Encode(job); err != nil {
			h.results <- result{job, false, jobResult}
			return err
		}
		if err := bufwriter.Flush(); err != nil {
			h.results <- result{job, false, jobResult}
			return err
		}

		if err := decoder.Decode(&jobResult); err != nil {
			h.results <- result{job, false, jobResult}
			return err
		}
		h.results <- result{job, true, jobResult}
	}
	return nil
}
