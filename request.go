package kraken

import (
	"encoding/json"
	"io"
)

type Request struct {
	Repository string
	Reference  string
}

func ParseRequest(reader io.Reader) (*Request, error) {
	decoder := json.NewDecoder(reader)
	var request Request
	err := decoder.Decode(&request)
	return &request, err
}
