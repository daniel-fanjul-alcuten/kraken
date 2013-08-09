package kraken

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Wrapper for the git command line
type Git struct {
	git string
	dir string
}

// If dir = "", the current directory is used
func NewGit(dir string) *Git {
	return &Git{"git", dir}
}

// If dir = "", the current directory is used
func NewGitFullPath(git, dir string) *Git {
	return &Git{git, dir}
}

func (g *Git) Dir() string {
	return g.dir
}

// Run 'git <args>'
func (g *Git) Run(stdin io.Reader, args ...string) ([]byte, error) {
	cmd := exec.Command(g.git, args...)
	cmd.Dir = g.dir
	if stdin != nil {
		cmd.Stdin = stdin
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("git %s: %s:\n%s\n", strings.Join(args, " "), err, output)
	}
	return output, nil
}

// Run 'git <args>'
func (g *Git) String(stdin io.Reader, args ...string) (string, error) {
	output, err := g.Run(stdin, args...)
	if err == nil {
		return strings.TrimSpace(string(output)), err
	}
	return string(output), err
}
