package kraken

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {

	buffer := &bytes.Buffer{}

	ref := &RequestRef{"foo", "bar"}
	if err := ref.Encode(buffer); err != nil {
		t.Error(err)
	}

	json := strings.TrimSpace(buffer.String())
	if json != `{"Repository":"foo","Reference":"bar"}` {
		t.Error(json)
	}
}

func TestDecode(t *testing.T) {

	buffer := &bytes.Buffer{}

	ref1 := &RequestRef{"foo", "bar"}
	if err := ref1.Encode(buffer); err != nil {
		t.Error(err)
	}

	var ref2 RequestRef
	if err := ref2.Decode(buffer); err != nil {
		t.Error(err)
	}
	if ref2 != *ref1 {
		t.Error(ref2)
	}
}
