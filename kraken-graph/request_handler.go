package main

import (
	"bufio"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	"io"
)

type requestHandler struct {
	requests chan<- Request
}

func (h *requestHandler) handle(reader io.Reader, writer io.Writer) error {

	bufreader := bufio.NewReader(reader)
	decoder := gob.NewDecoder(bufreader)
	var request Request
	for {
		if err := decoder.Decode(&request); err != nil {
			return err
		}
		h.requests <- request
	}
}
