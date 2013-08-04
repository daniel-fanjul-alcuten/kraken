package kraken

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func (config *Config) Equal(config2 *Config) bool {
	if len(config.Jobs) != len(config2.Jobs) {
		return false
	}
	for i, job := range config.Jobs {
		if job.ImportPath != config2.Jobs[i].ImportPath {
			return false
		}
	}
	return true
}

func TestEncode(t *testing.T) {

	buffer := &bytes.Buffer{}

	config := &Config{[]GoJobConfig{GoJobConfig{"foo/bar/baz"}}}
	if err := config.Encode(buffer); err != nil {
		t.Error(err)
	}

	json := strings.TrimSpace(buffer.String())
	if json != `{"Jobs":[{"ImportPath":"foo/bar/baz"}]}` {
		t.Error(json)
	}
}

func TestDecode(t *testing.T) {

	buffer := &bytes.Buffer{}

	config1 := &Config{[]GoJobConfig{GoJobConfig{"foo/bar/baz"}}}
	if err := config1.Encode(buffer); err != nil {
		t.Error(err)
	}

	var config2 Config
	if err := config2.Decode(buffer); err != nil {
		t.Error(err)
	}
	if !config2.Equal(config1) {
		t.Error(config2)
	}
}

func TestKrakenConfiguration(t *testing.T) {

	file, err := os.Open("../kraken.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	config1 := &Config{[]GoJobConfig{GoJobConfig{"github.com/daniel-fanjul-alcuten/kraken"}}}

	var config2 Config
	if err := config2.Decode(file); err != nil {
		t.Error(err)
	}
	if !config2.Equal(config1) {
		t.Error(config2)
	}
}
