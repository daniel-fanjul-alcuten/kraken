package kraken

import (
	"encoding/json"
	"io"
)

// A ref to be built. The repository is not accessed directly, so it is an id instead of an url.
type RequestRef struct {
	Repository string
	Reference  string
}

func (ref *RequestRef) Decode(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(ref)
}

func (ref *RequestRef) Encode(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(ref)
}
