package kraken

import (
	"os/exec"
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

func (g *Git) Cmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.dir
	return cmd
}
