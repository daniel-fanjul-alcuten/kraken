package kraken

import (
	"bytes"
	"testing"
)

func TestParseRequest(t *testing.T) {

	buffer := &bytes.Buffer{}
	if _, err := buffer.WriteString(`{"Repository":"foo","Reference":"bar"}`); err != nil {
		t.Error(err)
	}

	request, err := ParseRequest(buffer)
	if request == nil {
		t.Fatal()
	}
	if err != nil {
		t.Error(err)
	}

	if request.Repository != "foo" {
		t.Fatal(request.Repository)
	}
	if request.Reference != "bar" {
		t.Fatal(request.Reference)
	}
}
