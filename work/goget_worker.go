package kraken

import (
	"bytes"
	. "github.com/daniel-fanjul-alcuten/kraken/git"
	. "github.com/daniel-fanjul-alcuten/kraken/gob"
	"os"
	"os/exec"
	"strings"
)

type GoGetWorker struct {
	Job GoGetJob
}

const gopath = "GOPATH"
const gopath_prefix = gopath + "="

func (w *GoGetWorker) Work(git *Git) JobResult {
	dir := git.Dir()
	importPath := w.Job.ImportPath
	repository := w.Job.Repoquest
	reference := w.Job.Request

	log_buffer := &bytes.Buffer{}

	// clean index and working copy
	mktree, err := git_string(log_buffer, git, "mktree")
	if err != nil {
		return JobResult{false, log_buffer.String()}
	}
	committree, err := git_string(log_buffer, git, "commit-tree", mktree)
	if err != nil {
		return JobResult{false, log_buffer.String()}
	}
	if _, err := git_run(log_buffer, git, "checkout", "-f", committree); err != nil {
		return JobResult{false, log_buffer.String()}
	}
	if _, err := git_run(log_buffer, git, "clean", "-fdx"); err != nil {
		return JobResult{false, log_buffer.String()}
	}

	// fetch and checkout working copy
	if _, err := git_run(log_buffer, git, "fetch", repository, "+"+reference+":"+reference); err != nil {
		return JobResult{false, log_buffer.String()}
	}
	if _, err := git_run(log_buffer, git, "read-tree", "-u", "--prefix=src/"+importPath+"/", reference+":"); err != nil {
		return JobResult{false, log_buffer.String()}
	}

	// set GOPATH
	environ := os.Environ()
	ok := false
	for index, value := range environ {
		if strings.HasPrefix(value, gopath_prefix) {
			environ[index] = gopath_prefix + dir
			ok = true
			break
		}
	}
	if !ok {
		environ = append(environ, gopath_prefix+dir)
	}

	// compile and test
	for _, op := range []string{"install", "test"} {
		log_buffer.WriteString("$ go ")
		log_buffer.WriteString(op)
		log_buffer.WriteString(" ./...\n")
		install := exec.Command("go", op, "./...")
		install.Dir = dir
		install.Env = environ
		install.Stdout = log_buffer
		install.Stderr = log_buffer
		if err := install.Run(); err != nil {
			switch e := err.(type) {
			case *exec.ExitError:
				if !e.Success() {
					return JobResult{false, log_buffer.String()}
				}
			default:
				return JobResult{false, log_buffer.String()}
			}
		}
	}

	if _, err := git_run(log_buffer, git, "gc", "--auto"); err != nil {
		return JobResult{false, log_buffer.String()}
	}

	return JobResult{true, log_buffer.String()}
}

func git_run(log_buffer *bytes.Buffer, git *Git, args ...string) ([]byte, error) {
	log_buffer.WriteString("$ git ")
	log_buffer.WriteString(strings.Join(args, " "))
	log_buffer.WriteString("\n")
	output, err := git.Run(nil, args...)
	log_buffer.Write(output)
	return output, err
}

func git_string(log_buffer *bytes.Buffer, git *Git, args ...string) (string, error) {
	output, err := git_run(log_buffer, git, args...)
	return strings.TrimSpace(string(output)), err
}
