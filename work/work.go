package kraken

import (
	"bufio"
	"encoding/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	"io"
	"log"
)

func Work(reader io.Reader, writer io.Writer, logger *log.Logger) error {

	td := NewTempDir("")
	defer td.Cleanup()

	dir, err := td.NewDir()
	if err != nil {
		return err
	}

	git := NewGit(dir)
	if _, err := git.Run(nil, "init"); err != nil {
		return err
	}

	bufreader := bufio.NewReader(reader)
	decoder := gob.NewDecoder(bufreader)

	bufwriter := bufio.NewWriter(writer)
	encoder := gob.NewEncoder(bufwriter)

	for {
		var job GoGetJob
		if err := decoder.Decode(&job); err != nil {
			return err
		}
		logger.Printf("New Job: %#v", job)

		worker := &GoGetWorker{job}
		result := worker.Work(git)
		logger.Printf("Output:\n%sSuccess: %v\n\n", result.Output, result.Success)

		if err := encoder.Encode(result); err != nil {
			return err
		}
		if err := bufwriter.Flush(); err != nil {
			return err
		}
	}
}
