package kraken

import (
	"encoding/json"
	"io"
)

type GoRepositoryJob struct {
	ImportPath string
}

type Configuration struct {
	Jobs []GoRepositoryJob
}

func ParseConfiguration(reader io.Reader) (*Configuration, error) {
	decoder := json.NewDecoder(reader)
	var config Configuration
	err := decoder.Decode(&config)
	return &config, err
}
