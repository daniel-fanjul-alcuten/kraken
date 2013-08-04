package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	. "github.com/daniel-fanjul-alcuten/kraken/config"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	r "github.com/daniel-fanjul-alcuten/kraken/request"
	"io"
	"strconv"
	"strings"
)

func submit(git *Git, repoquest string, conn io.Writer, requests ...string) error {

	writer := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(writer)

	for _, request := range requests {

		repository := ""
		reference := ""
		time := int64(0)
		if output, err := git.Cmd("cat-file", "tag", request).Output(); err != nil {
			return fmt.Errorf("git cat-file: %s", err)
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
			req, err := r.ParseRequest(buffer)
			if err != nil {
				return fmt.Errorf("json decoding: request: %s", err)
			}
			repository = req.Repository
			reference = req.Reference
		}

		output, err := git.Cmd("show", request+":kraken.json").Output()
		if err != nil {
			return fmt.Errorf("git show: kraken.json not found: %s", err)
		}
		configuration, err := ParseConfiguration(bytes.NewBuffer(output))
		if err != nil {
			return fmt.Errorf("json deconding: kraken.json: %s", err)
		}

		jobs := make([]GoJob, len(configuration.Jobs))
		for i, job := range configuration.Jobs {
			jobs[i] = GoJob{job.ImportPath}
		}

		req := Request{repoquest, request, repository, reference, time, jobs}
		if err := encoder.Encode(req); err != nil {
			return fmt.Errorf("gob encoding: %s", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}
