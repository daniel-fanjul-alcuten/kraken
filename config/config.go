package kraken

import (
	"encoding/json"
	"io"
)

// A go getable repository.
type GoJobConfig struct {
	// The folder where the sources must be placed: $GOPATH/src/<ImportPath>/<working copy>
	ImportPath string
}

// The configuration of the project is a list of jobs, that can run in different workers, be cached, and depend on others.
type Config struct {
	Jobs []GoJobConfig
}

func (config *Config) Decode(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(config)
}

func (config *Config) Encode(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(config)
}
