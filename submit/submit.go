package kraken

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/config"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	. "github.com/daniel-fanjul-alcuten/kraken/json"
	"io"
	"strconv"
	"strings"
)

func Submit(git *Git, repoquest string, writer io.Writer, requests ...string) error {

	bufwriter := bufio.NewWriter(writer)
	encoder := gob.NewEncoder(bufwriter)

	for _, request := range requests {

		repository := ""
		reference := ""
		time := int64(0)
		if output, err := git.Run(nil, "cat-file", "tag", request); err != nil {
			return err
		} else {
			buffer := bytes.NewBuffer(output)
			buffer.ReadString('\n')              // remove object
			buffer.ReadString('\n')              // remove type
			buffer.ReadString('\n')              // remove tag
			tagger, _ := buffer.ReadString('\n') // tagger
			split := strings.Split(tagger, " ")
			if len(split) > 2 {
				time, _ = strconv.ParseInt(split[len(split)-2], 10, 64)
			}
			buffer.ReadString('\n') // remove empty line
			var ref RequestRef
			if err := ref.Decode(buffer); err != nil {
				return fmt.Errorf("json decoding: request: %s", err)
			}
			repository = ref.Repository
			reference = ref.Reference
		}

		output, err := git.Run(nil, "show", request+":kraken.json")
		if err != nil {
			return err
		}
		var config Config
		if config.Decode(bytes.NewBuffer(output)); err != nil {
			return fmt.Errorf("json deconding: kraken.json: %s", err)
		}

		jobs := make([]GoGetRequest, len(config.Jobs))
		for i, job := range config.Jobs {
			jobs[i] = GoGetRequest{job.ImportPath}
		}

		req := Request{repoquest, request, repository, reference, time, jobs}
		if err := encoder.Encode(req); err != nil {
			return fmt.Errorf("gob encoding: %s", err)
		}
	}

	if err := bufwriter.Flush(); err != nil {
		return err
	}
	return nil
}
