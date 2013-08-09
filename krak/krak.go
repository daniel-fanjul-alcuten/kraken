package main

import (
	"bytes"
	. "github.com/daniel-fanjul-alcuten/kraken/config"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/ioutil"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const GOPATH = "GOPATH"
const GOPATH_PREFIX = GOPATH + "="

func krak(dir string, refs ...string) error {

	td := NewTempDir("")
	defer td.Cleanup()

	work, err := td.NewDir()
	if err != nil {
		return err
	}

	git := NewGit(work)
	if _, err := git.Run(nil, "init"); err != nil {
		return err
	}

	for _, ref := range refs {

		// fetch and parse configuration
		if _, err := git.Run(nil, "fetch", dir, ref); err != nil {
			return err
		}
		config_bytes, err := git.Run(nil, "show", "FETCH_HEAD:kraken.json")
		if err != nil {
			return err
		}
		var config Config
		if err := config.Decode(bytes.NewBuffer(config_bytes)); err != nil {
			return err
		}

		for _, job := range config.Jobs {
			log_buffer := &bytes.Buffer{}

			// clean working copy
			mktree, err := git.String(nil, "mktree")
			if err != nil {
				return err
			}
			if _, err := git.Run(nil, "read-tree", mktree); err != nil {
				return err
			}
			if _, err := git.Run(nil, "clean", "-fdx"); err != nil {
				return err
			}

			// setup working copy
			if _, err := git.Run(nil, "read-tree", "-u", "--prefix=src/"+job.ImportPath+"/", "FETCH_HEAD:"); err != nil {
				return err
			}

			// commit
			writetree, err := git.String(nil, "write-tree")
			if err != nil {
				return err
			}
			committree, err := git.String(nil, "commit-tree", writetree)
			if err != nil {
				return err
			}

			// set GOPATH
			environ := os.Environ()
			ok := false
			for index, value := range environ {
				if strings.HasPrefix(value, GOPATH_PREFIX) {
					environ[index] = GOPATH_PREFIX + work
					ok = true
					break
				}
			}
			if !ok {
				environ = append(environ, GOPATH_PREFIX+work)
			}

			// compile and test
			for _, op := range []string{"install", "test"} {
				log_buffer.WriteString("$ go ")
				log_buffer.WriteString(op)
				log_buffer.WriteString(" ./...\n")
				install := exec.Command("go", op, "./...")
				install.Dir = work
				install.Env = environ
				install.Stdout = log_buffer
				install.Stderr = log_buffer
				if err := install.Run(); err != nil {
					switch err.(type) {
					case *exec.ExitError:
					default:
						return err
					}
					break
				}
			}

			// write log
			if err := ioutil.WriteFile(path.Join(work, "kraken.log"), log_buffer.Bytes(), 0666); err != nil {
				return err
			}
			if _, err := git.Run(nil, "add", "-A"); err != nil {
				return err
			}

			// commit
			writetree2, err := git.String(nil, "write-tree")
			if err != nil {
				return err
			}
			committree2, err := git.String(nil, "commit-tree", "-p", committree, writetree2)
			if err != nil {
				return err
			}

			// push
			if _, err := git.Run(nil, "push", dir, "+"+committree2+":refs/build/"+ref); err != nil {
				return err
			}
		}
	}

	return nil
}
