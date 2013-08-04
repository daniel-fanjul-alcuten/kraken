package kraken

import (
	"fmt"
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

// Run 'git init'
func (g *Git) Init() error {
	if err := exec.Command(g.git, "init", g.dir).Run(); err != nil {
		return fmt.Errorf("git init %s: %s", g.dir, err)
	}
	return nil
}

// Run 'git init --bare'
func (g *Git) InitBare() error {
	if err := exec.Command(g.git, "init", "--bare", g.dir).Run(); err != nil {
		return fmt.Errorf("git init --bare %s: %s", g.dir, err)
	}
	return nil
}

// Run 'git config <name> <value>'
func (g *Git) SetConfig(name, value string) error {
	if err := g.Cmd("config", name, value).Run(); err != nil {
		return fmt.Errorf("git config %s %s: %s", name, value, err)
	}
	return nil
}

// Run 'git config <name>'
func (g *Git) Config(name string) (string, error) {
	output, err := g.Cmd("config", name).Output()
	if err != nil {
		return "", fmt.Errorf("git config %s: %s", name, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// Run 'git <args>...'
func (g *Git) Cmd(args ...string) *exec.Cmd {
	cmd := exec.Command(g.git, args...)
	cmd.Dir = g.dir
	return cmd
}
