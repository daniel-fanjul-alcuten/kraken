package kraken

import (
	"fmt"
	"os/exec"
	"strings"
)

type Git struct {
	dir string
}

func NewGit(dir string) *Git {
	return &Git{dir}
}

func (g *Git) Init() error {
	return exec.Command("git", "init", g.dir).Run()
}

func (g *Git) InitBare() error {
	return exec.Command("git", "init", "--bare", g.dir).Run()
}

func (g *Git) SetConfig(name, value string) error {
	if err := g.Cmd("config", name, value).Run(); err != nil {
		return fmt.Errorf("git config %s %s", name, value)
	}
	return nil
}

func (g *Git) Config(name string) (string, error) {
	output, err := g.Cmd("config", name).Output()
	if err != nil {
		return "", fmt.Errorf("git config %s: not found", name)
	}
	return strings.TrimSpace(string(output)), nil
}

func (g *Git) Cmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.dir
	return cmd
}
