package kraken

import (
	"bytes"
	"testing"
)

func TestConfigurationParse(t *testing.T) {

	buffer := &bytes.Buffer{}
	if _, err := buffer.WriteString(`{"Jobs":[{"ImportPath":"foo/bar/baz"}]}`); err != nil {
		t.Error(err)
	}

	config, err := Parse(buffer)
	if config == nil {
		t.Fatal()
	}
	if err != nil {
		t.Error(err)
	}

	if len(config.Jobs) != 1 {
		t.Fatal(len(config.Jobs))
	}
	if config.Jobs[0].ImportPath != "foo/bar/baz" {
		t.Error(config.Jobs[0].ImportPath)
	}
}
